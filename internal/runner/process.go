package runner

import (
	"bytes"
	"context"
	"os/exec"
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
