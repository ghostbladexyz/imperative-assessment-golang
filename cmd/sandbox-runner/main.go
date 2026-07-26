package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go/format"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/sandboxprotocol"
)

const (
	maxRequestBytes = 1024 * 1024
	maxSourceBytes  = 192 * 1024
	resultMarker    = "__IMPERATIVE_ASSESSMENT_RESULT__"
)

func main() {
	response := execute(readRequest(os.Stdin))
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(response); err != nil {
		os.Exit(1)
	}
}

func readRequest(input io.Reader) (sandboxprotocol.Request, error) {
	var request sandboxprotocol.Request
	decoder := json.NewDecoder(io.LimitReader(input, maxRequestBytes+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return request, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return request, errors.New("request must contain exactly one JSON value")
	}
	if err := request.Validate(maxSourceBytes); err != nil {
		return request, err
	}
	return request, nil
}

func execute(request sandboxprotocol.Request, requestErr error) sandboxprotocol.Response {
	if requestErr != nil {
		return internalFailure("The sandbox request was invalid.")
	}
	if err := primeBuildCache("/opt/go-build", "/tmp/go-build"); err != nil {
		return internalFailure("The sandbox could not prepare its compiler cache.")
	}
	workspace, err := os.MkdirTemp("/workspace", "submission-")
	if err != nil {
		return internalFailure("The sandbox could not prepare its temporary workspace.")
	}
	defer os.RemoveAll(workspace)

	formattedSource := request.Source
	if formatted, formatErr := format.Source([]byte(request.Source)); formatErr == nil {
		formattedSource = string(formatted)
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(formattedSource), 0o600); err != nil {
		return internalFailure("The sandbox could not prepare the submitted source.")
	}
	if err := os.WriteFile(filepath.Join(workspace, "assessment_harness.go"), []byte(request.Harness), 0o600); err != nil {
		return internalFailure("The sandbox could not prepare the controlled tests.")
	}

	executable := filepath.Join(workspace, "assessment-program")
	compileCtx, cancelCompile := context.WithTimeout(
		context.Background(),
		time.Duration(request.CompileTimeoutMS)*time.Millisecond,
	)
	compileStdout, compileStderr, compileErr, compileLimited := runCommand(
		compileCtx,
		workspace,
		request.OutputLimitBytes,
		"go",
		"build",
		"-p=1",
		"-trimpath",
		"-o",
		executable,
		"solution.go",
		"assessment_harness.go",
	)
	cancelCompile()
	if compileErr != nil {
		response := sandboxprotocol.Response{
			Status:          sandboxprotocol.StatusCompile,
			CompileError:    cleanCompilerError(compileStderr),
			Stdout:          compileStdout,
			Stderr:          compileStderr,
			FormattedSource: formattedSource,
			Results:         []sandboxprotocol.TestResult{},
		}
		switch {
		case compileLimited:
			response.Status = sandboxprotocol.StatusOutput
			response.CompileError = "Compiler output exceeded the configured limit."
		case errors.Is(compileCtx.Err(), context.DeadlineExceeded):
			response.Status = sandboxprotocol.StatusCompileTimeout
			response.CompileError = "Compilation exceeded the configured time limit."
		}
		return response
	}

	runtimeCtx, cancelRuntime := context.WithTimeout(
		context.Background(),
		time.Duration(request.RuntimeTimeoutMS)*time.Millisecond,
	)
	runStdout, runStderr, runErr, runLimited := runCommand(
		runtimeCtx,
		workspace,
		request.OutputLimitBytes,
		executable,
	)
	cancelRuntime()
	results := parseResults(runStdout)
	response := sandboxprotocol.Response{
		Status:          sandboxprotocol.StatusSuccess,
		Stdout:          stripMarkers(runStdout),
		Stderr:          runStderr,
		FormattedSource: formattedSource,
		Results:         results,
	}
	if runErr != nil {
		switch {
		case runLimited:
			response.Status = sandboxprotocol.StatusOutput
			response.RuntimeError = "Program output exceeded the configured limit and execution was terminated."
		case errors.Is(runtimeCtx.Err(), context.DeadlineExceeded):
			response.Status = sandboxprotocol.StatusRuntimeTimeout
			response.RuntimeError = "Execution exceeded the configured time limit and was terminated."
		default:
			response.Status = sandboxprotocol.StatusRuntime
			response.RuntimeError = strings.TrimSpace(runStderr)
			if response.RuntimeError == "" {
				response.RuntimeError = "The program exited before completing every test."
			}
		}
		return response
	}

	expected := make(map[string]string, len(request.Tests))
	for _, test := range request.Tests {
		expected[test.ID] = test.Expected
	}
	seen := make(map[string]bool, len(results))
	for _, result := range results {
		seen[result.ID] = true
		if result.Failure != "" {
			response.Status = sandboxprotocol.StatusRuntime
		} else if expected[result.ID] != result.Actual && response.Status == sandboxprotocol.StatusSuccess {
			response.Status = sandboxprotocol.StatusAssertion
		}
	}
	if len(seen) != len(expected) {
		response.Status = sandboxprotocol.StatusRuntime
		response.RuntimeError = "The program stopped before every test produced a result."
	}
	return response
}

func internalFailure(message string) sandboxprotocol.Response {
	return sandboxprotocol.Response{
		Status:       sandboxprotocol.StatusInternal,
		RuntimeError: message,
		Results:      []sandboxprotocol.TestResult{},
	}
}

type limitedBuffer struct {
	mu      sync.Mutex
	buffer  bytes.Buffer
	limit   int
	limited bool
	cancel  context.CancelFunc
}

func (buffer *limitedBuffer) Write(data []byte) (int, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	remaining := buffer.limit - buffer.buffer.Len()
	if remaining > 0 {
		count := len(data)
		if count > remaining {
			count = remaining
		}
		_, _ = buffer.buffer.Write(data[:count])
		remaining = count
	} else {
		remaining = 0
	}
	if len(data) > remaining && !buffer.limited {
		buffer.limited = true
		buffer.cancel()
	}
	return len(data), nil
}

func (buffer *limitedBuffer) String() string {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.String()
}

func runCommand(parent context.Context, directory string, limit int, name string, args ...string) (string, string, error, bool) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	stdout := &limitedBuffer{limit: limit, cancel: cancel}
	stderr := &limitedBuffer{limit: limit, cancel: cancel}
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = directory
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err, stdout.limited || stderr.limited
}

func parseResults(stdout string) []sandboxprotocol.TestResult {
	results := make([]sandboxprotocol.TestResult, 0)
	var pendingOutput []string
	for _, line := range strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n") {
		prefix, payload, found := splitResultLine(line)
		if !found {
			pendingOutput = append(pendingOutput, line)
			continue
		}
		if prefix != "" {
			pendingOutput = append(pendingOutput, prefix)
		}
		var result sandboxprotocol.TestResult
		if json.Unmarshal([]byte(payload), &result) == nil {
			result.Stdout = strings.TrimSpace(strings.Join(pendingOutput, "\n"))
			results = append(results, result)
			pendingOutput = nil
		}
	}
	return results
}

func stripMarkers(stdout string) string {
	lines := strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n")
	kept := lines[:0]
	for _, line := range lines {
		prefix, _, found := splitResultLine(line)
		if found {
			kept = append(kept, prefix)
		} else {
			kept = append(kept, line)
		}
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

func splitResultLine(line string) (string, string, bool) {
	index := strings.Index(line, resultMarker)
	if index < 0 {
		return "", "", false
	}
	payload := line[index+len(resultMarker):]
	if !json.Valid([]byte(payload)) {
		return "", "", false
	}
	return line[:index], payload, true
}

func cleanCompilerError(message string) string {
	message = strings.ReplaceAll(message, "\\", "/")
	message = strings.ReplaceAll(message, "./solution.go:", "solution.go:")
	message = strings.ReplaceAll(message, "./assessment_harness.go:", "assessment tests:")
	if strings.TrimSpace(message) == "" {
		return "Compilation failed."
	}
	return strings.TrimSpace(message)
}

func primeBuildCache(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0o600)
	})
}
