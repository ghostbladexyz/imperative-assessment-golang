package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"sync"
)

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
	return runCommandInput(parent, limit, nil, name, args...)
}

func runCommandInput(parent context.Context, limit int, stdin []byte, name string, args ...string) (string, string, error, bool) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	stdout := &limitedBuffer{limit: limit, cancel: cancel}
	stderr := &limitedBuffer{limit: limit, cancel: cancel}
	command := exec.CommandContext(ctx, name, args...)
	if stdin != nil {
		command.Stdin = bytes.NewReader(stdin)
	}
	command.Stdout = stdout
	command.Stderr = stderr
	err := command.Run()
	return stdout.String(), stderr.String(), err, stdout.limited || stderr.limited
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
