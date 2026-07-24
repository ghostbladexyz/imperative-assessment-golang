package runner

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

const (
	MaxSourceBytes = 192 * 1024
	MaxOutputBytes = 256 * 1024
)

var packagePattern = regexp.MustCompile(`(?m)^\s*package\s+main\s*(?:\r?\n)?`)

type TestResult struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Purpose    string  `json:"purpose"`
	Input      string  `json:"input"`
	Expected   string  `json:"expected"`
	Actual     string  `json:"actual"`
	Passed     bool    `json:"passed"`
	Status     string  `json:"status"`
	Failure    string  `json:"failure,omitempty"`
	DurationMS float64 `json:"durationMs"`
}

type RunResult struct {
	LevelID       int          `json:"levelId"`
	Passed        bool         `json:"passed"`
	PassedCount   int          `json:"passedCount"`
	TotalCount    int          `json:"totalCount"`
	CompileError  string       `json:"compileError,omitempty"`
	RuntimeError  string       `json:"runtimeError,omitempty"`
	TimedOut      bool         `json:"timedOut"`
	Stopped       bool         `json:"stopped"`
	Stdout        string       `json:"stdout"`
	Stderr        string       `json:"stderr"`
	DurationMS    float64      `json:"durationMs"`
	Results       []TestResult `json:"results"`
	FormattedCode string       `json:"formattedCode,omitempty"`
	SourceHash    string       `json:"sourceHash"`
	Receipt       string       `json:"receipt,omitempty"`
}

type ReceiptIssuer interface {
	Issue(levelID int, sourceHash string) (string, error)
}

type Runner struct {
	goBinary string
	slots    chan struct{}
	receipts ReceiptIssuer
}

func New(goBinary string, maxConcurrent int, receipts ReceiptIssuer) *Runner {
	if goBinary == "" {
		goBinary = "go"
	}
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	return &Runner{goBinary: goBinary, slots: make(chan struct{}, maxConcurrent), receipts: receipts}
}

func PrepareSource(source string) string {
	source = strings.TrimPrefix(source, "\uFEFF")
	source = packagePattern.ReplaceAllString(source, "")
	return "package main\n\n" + strings.TrimSpace(source) + "\n"
}

func (runner *Runner) Format(ctx context.Context, source string) (string, error) {
	if len(source) > MaxSourceBytes {
		return "", fmt.Errorf("source exceeds %d bytes", MaxSourceBytes)
	}
	command := exec.CommandContext(ctx, "gofmt")
	command.Stdin = strings.NewReader(PrepareSource(source))
	var output bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &output
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
	}
	formatted := packagePattern.ReplaceAllString(output.String(), "")
	return strings.TrimSpace(formatted) + "\n", nil
}

func (runner *Runner) Run(ctx context.Context, level assessment.Level, source string, testIDs []string) RunResult {
	started := time.Now()
	result := RunResult{LevelID: level.ID}
	if len(source) > MaxSourceBytes {
		result.CompileError = fmt.Sprintf("Source is too large. The limit is %d KiB.", MaxSourceBytes/1024)
		result.DurationMS = elapsedMS(started)
		return result
	}
	selected, valid := assessment.SelectTests(level, testIDs)
	if !valid {
		result.CompileError = "The request contains an unknown test identifier."
		result.DurationMS = elapsedMS(started)
		return result
	}
	result.TotalCount = len(selected)
	for _, current := range selected {
		result.Results = append(result.Results, TestResult{
			ID: current.ID, Name: current.Name, Purpose: current.Purpose,
			Input: current.Input, Expected: current.Expected, Status: "pending",
		})
	}
	prepared := PrepareSource(source)
	sum := sha256.Sum256([]byte(prepared))
	result.SourceHash = hex.EncodeToString(sum[:])

	select {
	case runner.slots <- struct{}{}:
		defer func() { <-runner.slots }()
	case <-ctx.Done():
		result.Stopped = true
		result.DurationMS = elapsedMS(started)
		return result
	}

	tempDir, err := os.MkdirTemp("", "imperative-assessment-*")
	if err != nil {
		result.RuntimeError = "Could not create an isolated temporary run directory."
		result.DurationMS = elapsedMS(started)
		return result
	}
	defer os.RemoveAll(tempDir)

	solutionPath := filepath.Join(tempDir, "solution.go")
	harnessPath := filepath.Join(tempDir, "assessment_harness.go")
	executablePath := filepath.Join(tempDir, "assessment-runner")
	if runtime.GOOS == "windows" {
		executablePath += ".exe"
	}
	if err := os.WriteFile(solutionPath, []byte(prepared), 0o600); err != nil {
		result.RuntimeError = "Could not prepare the submitted source."
		result.DurationMS = elapsedMS(started)
		return result
	}
	if err := os.WriteFile(harnessPath, []byte(level.BuildHarness(selected)), 0o600); err != nil {
		result.RuntimeError = "Could not prepare the controlled test harness."
		result.DurationMS = elapsedMS(started)
		return result
	}

	buildCtx, cancelBuild := context.WithTimeout(ctx, 20*time.Second)
	buildStdout, buildStderr, buildErr, buildLimited := runCommand(buildCtx, MaxOutputBytes, runner.goBinary, "build", "-trimpath", "-o", executablePath, solutionPath, harnessPath)
	cancelBuild()
	result.Stdout = buildStdout
	result.Stderr = buildStderr
	if buildErr != nil {
		if errors.Is(buildCtx.Err(), context.DeadlineExceeded) {
			result.TimedOut = true
			result.CompileError = "Compilation exceeded the 20 second limit."
		} else if errors.Is(ctx.Err(), context.Canceled) {
			result.Stopped = true
		} else if buildLimited {
			result.CompileError = "Compiler output exceeded the 256 KiB limit."
		} else {
			result.CompileError = cleanCompilerError(buildStderr)
		}
		markAll(result.Results, "compile", result.CompileError)
		result.DurationMS = elapsedMS(started)
		return result
	}

	runCtx, cancelRun := context.WithTimeout(ctx, 4*time.Second)
	runStdout, runStderr, runErr, runLimited := runCommand(runCtx, MaxOutputBytes, executablePath)
	cancelRun()
	result.Stdout = stripMarkers(runStdout)
	result.Stderr = runStderr
	wireResults := parseResults(runStdout)
	byID := make(map[string]wireResult, len(wireResults))
	for _, item := range wireResults {
		byID[item.ID] = item
	}
	for index := range result.Results {
		wire, found := byID[result.Results[index].ID]
		if !found {
			result.Results[index].Status = "runtime"
			result.Results[index].Failure = "The program stopped before this test produced a result."
			continue
		}
		result.Results[index].Actual = wire.Actual
		result.Results[index].DurationMS = wire.DurationMS
		result.Results[index].Failure = wire.Failure
		if wire.Failure != "" {
			result.Results[index].Status = "runtime"
			continue
		}
		if wire.Actual == result.Results[index].Expected {
			result.Results[index].Passed = true
			result.Results[index].Status = "pass"
			result.PassedCount++
		} else {
			result.Results[index].Status = "assertion"
		}
	}
	if runErr != nil {
		switch {
		case errors.Is(runCtx.Err(), context.DeadlineExceeded):
			result.TimedOut = true
			result.RuntimeError = "Execution exceeded the 4 second limit and was terminated."
		case errors.Is(ctx.Err(), context.Canceled):
			result.Stopped = true
			result.RuntimeError = "Execution was stopped."
		case runLimited:
			result.RuntimeError = "Program output exceeded the 256 KiB limit and execution was terminated."
		default:
			result.RuntimeError = strings.TrimSpace(runStderr)
			if result.RuntimeError == "" {
				result.RuntimeError = "The program exited before completing every test."
			}
		}
	}
	result.Passed = len(testIDs) == 0 && result.PassedCount == len(level.Tests)
	if result.Passed && runner.receipts != nil {
		result.Receipt, _ = runner.receipts.Issue(level.ID, result.SourceHash)
	}
	formatCtx, cancelFormat := context.WithTimeout(context.Background(), 3*time.Second)
	result.FormattedCode, _ = runner.Format(formatCtx, source)
	cancelFormat()
	result.DurationMS = elapsedMS(started)
	return result
}

type wireResult struct {
	ID         string  `json:"id"`
	Actual     string  `json:"actual"`
	Failure    string  `json:"failure"`
	DurationMS float64 `json:"durationMs"`
}

func parseResults(stdout string) []wireResult {
	var results []wireResult
	for _, line := range strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n") {
		if !strings.HasPrefix(line, "__IMPERATIVE_ASSESSMENT_RESULT__") {
			continue
		}
		var result wireResult
		if json.Unmarshal([]byte(strings.TrimPrefix(line, "__IMPERATIVE_ASSESSMENT_RESULT__")), &result) == nil {
			results = append(results, result)
		}
	}
	return results
}

func stripMarkers(stdout string) string {
	lines := strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n")
	kept := lines[:0]
	for _, line := range lines {
		if !strings.HasPrefix(line, "__IMPERATIVE_ASSESSMENT_RESULT__") {
			kept = append(kept, line)
		}
	}
	return strings.TrimSpace(strings.Join(kept, "\n"))
}

func markAll(results []TestResult, status, failure string) {
	for index := range results {
		results[index].Status = status
		results[index].Failure = failure
	}
}

func cleanCompilerError(message string) string {
	message = strings.ReplaceAll(message, "\\", "/")
	lines := strings.Split(strings.TrimSpace(message), "\n")
	for index, line := range lines {
		if marker := strings.Index(line, "/solution.go:"); marker >= 0 {
			lines[index] = "solution.go:" + line[marker+len("/solution.go:"):]
		} else if marker := strings.Index(line, "/assessment_harness.go:"); marker >= 0 {
			lines[index] = "assessment tests:" + line[marker+len("/assessment_harness.go:"):]
		}
	}
	return strings.Join(lines, "\n")
}

func elapsedMS(started time.Time) float64 {
	return float64(time.Since(started).Microseconds()) / 1000
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
		if len(data) < remaining {
			remaining = len(data)
		}
		_, _ = buffer.buffer.Write(data[:remaining])
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

func runCommand(parent context.Context, limit int, name string, args ...string) (string, string, error, bool) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	stdout := &limitedBuffer{limit: limit, cancel: cancel}
	stderr := &limitedBuffer{limit: limit, cancel: cancel}
	command := exec.CommandContext(ctx, name, args...)
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err, stdout.limited || stderr.limited
}

func SortedStatuses(results []TestResult) []string {
	statuses := make([]string, 0, len(results))
	for _, result := range results {
		statuses = append(statuses, result.Status)
	}
	sort.Strings(statuses)
	return statuses
}
