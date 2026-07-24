package runner

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

func TestDockerRunnerIntegration(t *testing.T) {
	if os.Getenv("IMPERATIVE_DOCKER_INTEGRATION") != "1" {
		t.Skip("set IMPERATIVE_DOCKER_INTEGRATION=1 to run tests against a real Docker daemon")
	}
	commands := processExecutor{}
	checkCtx, cancelCheck := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCheck()
	if err := CheckDocker(checkCtx, commands, "docker"); err != nil {
		t.Skipf("real Docker integration skipped: %v", err)
	}
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Minute)
	sandbox, err := NewDocker(buildCtx, DockerOptions{
		MaxConcurrent: 1,
		Commands:      commands,
	})
	cancelBuild()
	if err != nil {
		t.Fatal(err)
	}
	level, _ := assessment.FindLevel(1)
	isolationLevel := level
	isolationLevel.Instructions.AllowedPackages = append(
		append([]string(nil), level.Instructions.AllowedPackages...),
		"net",
		"os",
		"time",
	)

	t.Run("correct solution passes", func(t *testing.T) {
		result := sandbox.Run(context.Background(), level, correctLevelOneSolution, nil)
		if !result.Passed || result.PassedCount != len(level.Tests) {
			t.Fatalf("correct solution failed: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("incorrect solution fails assertions", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			level,
			`func NormalizeTokens(input string) []string { return []string{input} }`,
			[]string{"l1-case"},
		)
		if result.Passed || len(result.Results) != 1 || result.Results[0].Status != "assertion" {
			t.Fatalf("incorrect solution had unexpected result: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("invalid Go returns compiler error", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			level,
			`func NormalizeTokens(input string) []string { definitely not go }`,
			[]string{"l1-word"},
		)
		if result.CompileError == "" || !strings.Contains(result.CompileError, "solution.go:") {
			t.Fatalf("missing readable compiler error: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("infinite loop times out", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			level,
			`func NormalizeTokens(input string) []string { for {} }`,
			[]string{"l1-word"},
		)
		if !result.TimedOut {
			t.Fatalf("expected timeout: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("cancellation removes container", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan RunResult, 1)
		go func() {
			done <- sandbox.Run(
				ctx,
				level,
				`func NormalizeTokens(input string) []string { for {} }`,
				[]string{"l1-word"},
			)
		}()
		time.Sleep(750 * time.Millisecond)
		cancel()
		result := <-done
		if !result.Stopped {
			t.Fatalf("expected stopped result: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("external network is unavailable", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			isolationLevel,
			`import (
	"net"
	"time"
)
func NormalizeTokens(input string) []string {
	connection, err := net.DialTimeout("tcp", "1.1.1.1:80", 500*time.Millisecond)
	if err == nil {
		connection.Close()
		return []string{"network-was-available"}
	}
	return []string{input}
}`,
			[]string{"l1-word"},
		)
		if result.PassedCount != 1 {
			t.Fatalf("network isolation probe failed: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("container root is read only", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			isolationLevel,
			`import "os"
func NormalizeTokens(input string) []string {
	if os.WriteFile("/root-write-probe", []byte("unsafe"), 0600) == nil {
		return []string{"root-was-writable"}
	}
	return []string{input}
}`,
			[]string{"l1-word"},
		)
		if result.PassedCount != 1 {
			t.Fatalf("read-only root probe failed: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("host files are not mounted", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			isolationLevel,
			`import "os"
func NormalizeTokens(input string) []string {
	if _, err := os.ReadFile("/imperative-assessment-host-probe"); err == nil {
		return []string{"host-file-was-visible"}
	}
	return []string{input}
}`,
			[]string{"l1-word"},
		)
		if result.PassedCount != 1 {
			t.Fatalf("host file isolation probe failed: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("excessive output is capped", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			level,
			`func NormalizeTokens(input string) []string {
	for i := 0; i < 300000; i++ { print("output") }
	return []string{input}
}`,
			[]string{"l1-word"},
		)
		if result.FailureKind != FailureOutput || len(result.Stderr) > MaxOutputBytes {
			t.Fatalf("output was not capped: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("every starter and harness compiles", func(t *testing.T) {
		for _, current := range assessment.Levels() {
			result := sandbox.Run(
				context.Background(),
				current,
				current.StarterCode,
				[]string{current.Tests[0].ID},
			)
			if result.CompileError != "" {
				t.Fatalf("level %d starter did not compile: %s", current.ID, result.CompileError)
			}
			assertNoAssessmentContainers(t)
		}
	})
}

func assertNoAssessmentContainers(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stdout, stderr, err, _ := runCommand(
		ctx,
		64*1024,
		"docker",
		"ps",
		"-a",
		"--filter",
		"name=imperative-go-assessment-",
		"--format",
		"{{.Names}}",
	)
	if err != nil {
		t.Fatalf("inspect Docker containers: %v: %s", err, stderr)
	}
	if strings.TrimSpace(stdout) != "" {
		t.Fatalf("assessment containers remain: %s", stdout)
	}
}

const correctLevelOneSolution = `func NormalizeTokens(input string) []string {
	result := make([]string, 0)
	current := make([]rune, 0)
	flush := func() {
		if len(current) > 0 {
			result = append(result, string(current))
			current = current[:0]
		}
	}
	for _, r := range input {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			current = append(current, r)
		} else {
			flush()
		}
	}
	flush()
	return result
}`
