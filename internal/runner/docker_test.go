package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/sandboxprotocol"
)

type recordedCommand struct {
	name  string
	args  []string
	stdin []byte
}

type fakeExecutor struct {
	mu       sync.Mutex
	commands []recordedCommand
	run      func(context.Context, int, []byte, string, ...string) CommandResult
}

func (fake *fakeExecutor) Run(ctx context.Context, limit int, stdin []byte, name string, args ...string) CommandResult {
	fake.mu.Lock()
	fake.commands = append(fake.commands, recordedCommand{
		name: name, args: append([]string(nil), args...), stdin: append([]byte(nil), stdin...),
	})
	fake.mu.Unlock()
	if fake.run != nil {
		return fake.run(ctx, limit, stdin, name, args...)
	}
	return CommandResult{}
}

func (fake *fakeExecutor) snapshot() []recordedCommand {
	fake.mu.Lock()
	defer fake.mu.Unlock()
	return append([]recordedCommand(nil), fake.commands...)
}

func TestCheckDockerAvailability(t *testing.T) {
	tests := []struct {
		name    string
		run     func(int, []string) CommandResult
		wantErr string
	}{
		{
			name: "available",
			run: func(call int, _ []string) CommandResult {
				if call == 1 {
					return CommandResult{Stdout: "Docker version 29"}
				}
				return CommandResult{Stdout: "29.2.1"}
			},
		},
		{
			name: "missing cli",
			run: func(call int, _ []string) CommandResult {
				return CommandResult{Err: errors.New("not found")}
			},
			wantErr: "Docker CLI is unavailable",
		},
		{
			name: "daemon stopped",
			run: func(call int, _ []string) CommandResult {
				if call == 1 {
					return CommandResult{Stdout: "Docker version 29"}
				}
				return CommandResult{Err: errors.New("daemon unavailable")}
			},
			wantErr: "Docker Desktop is not running",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			call := 0
			fake := &fakeExecutor{run: func(_ context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
				call++
				return test.run(call, args)
			}}
			err := CheckDocker(context.Background(), fake, "docker")
			if test.wantErr == "" && err != nil {
				t.Fatal(err)
			}
			if test.wantErr != "" && (err == nil || !strings.Contains(err.Error(), test.wantErr)) {
				t.Fatalf("expected %q, got %v", test.wantErr, err)
			}
		})
	}
}

func TestImageInputHashChangesOnlyWithImageInputs(t *testing.T) {
	root := makeImageFixture(t)
	first, err := imageInputHash(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("unrelated"), 0o600); err != nil {
		t.Fatal(err)
	}
	second, err := imageInputHash(root)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("unrelated file changed the image hash")
	}
	if err := os.WriteFile(filepath.Join(root, "cmd", "sandbox-runner", "main.go"), []byte("package main\n// changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	third, err := imageInputHash(root)
	if err != nil {
		t.Fatal(err)
	}
	if first == third {
		t.Fatal("sandbox source did not change the image hash")
	}
}

func TestNewDockerReusesCachedImageAndBuildsMissingImage(t *testing.T) {
	for _, cached := range []bool{true, false} {
		t.Run(map[bool]string{true: "cached", false: "missing"}[cached], func(t *testing.T) {
			root := makeImageFixture(t)
			builds := 0
			fake := &fakeExecutor{run: func(_ context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
				switch args[0] {
				case "--version":
					return CommandResult{Stdout: "Docker version 29"}
				case "info":
					return CommandResult{Stdout: "29.2.1"}
				case "image":
					if cached {
						return CommandResult{Stdout: "sha256:image"}
					}
					return CommandResult{Err: errors.New("missing")}
				case "build":
					builds++
					return CommandResult{}
				default:
					t.Fatalf("unexpected command: %v", args)
					return CommandResult{}
				}
			}}
			instance, err := NewDocker(context.Background(), DockerOptions{
				ProjectRoot: root, Commands: fake, MaxConcurrent: 1,
			})
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(instance.Info().DockerImage, dockerImageRepository+":") {
				t.Fatalf("unexpected image %q", instance.Info().DockerImage)
			}
			wantBuilds := 0
			if !cached {
				wantBuilds = 1
			}
			if builds != wantBuilds {
				t.Fatalf("got %d builds, want %d", builds, wantBuilds)
			}
		})
	}
}

func TestDockerRunArgumentsAreRestrictedAndUseNoHostMounts(t *testing.T) {
	name := "imperative-go-assessment-00112233445566778899aabb"
	args := dockerRunArgs(name, "runner:image")
	for _, required := range []string{
		"--rm", "--interactive", "none", "--read-only", "ALL",
		"no-new-privileges", "256m", "nofile=128:128", "65532:65532",
		"CGO_ENABLED=0", "GOPROXY=off", "runner:image",
	} {
		if !slices.Contains(args, required) {
			t.Errorf("missing restricted argument %q in %v", required, args)
		}
	}
	joined := " " + strings.Join(args, " ") + " "
	for _, forbidden := range []string{" --volume ", " -v ", "/var/run/docker.sock", `C:\`, "/home/", "/Users/"} {
		if strings.Contains(joined, forbidden) {
			t.Errorf("unsafe host mount or path %q in %v", forbidden, args)
		}
	}
	if args[len(args)-1] != "runner:image" {
		t.Fatal("image must be the final Docker argument")
	}
}

func TestContainerNamesAreUniqueAndValidated(t *testing.T) {
	first, err := newContainerName(bytes.NewReader(bytes.Repeat([]byte{1}, 12)))
	if err != nil {
		t.Fatal(err)
	}
	second, err := newContainerName(bytes.NewReader(bytes.Repeat([]byte{2}, 12)))
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !containerNamePattern.MatchString(first) || !containerNamePattern.MatchString(second) {
		t.Fatalf("invalid names %q and %q", first, second)
	}
	if _, err := newContainerName(strings.NewReader("short")); err == nil {
		t.Fatal("expected entropy read failure")
	}
}

func TestDockerRunCleansUpAfterSuccessFailureAndOutputOverflow(t *testing.T) {
	tests := []struct {
		name        string
		runResult   CommandResult
		cleanup     CommandResult
		wantFailure string
	}{
		{
			name:      "success",
			runResult: sandboxSuccessResponse(t),
			cleanup:   CommandResult{Err: errors.New("gone"), Stderr: "No such container"},
		},
		{
			name:        "startup failure",
			runResult:   CommandResult{Err: errors.New("start failed")},
			cleanup:     CommandResult{},
			wantFailure: FailureStartup,
		},
		{
			name:        "output overflow",
			runResult:   CommandResult{Err: errors.New("limited"), OutputLimited: true},
			cleanup:     CommandResult{},
			wantFailure: FailureOutput,
		},
		{
			name:        "cleanup failure",
			runResult:   sandboxSuccessResponse(t),
			cleanup:     CommandResult{Err: errors.New("daemon failed")},
			wantFailure: FailureCleanup,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cleanupCalls := 0
			fake := &fakeExecutor{run: func(_ context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
				if args[0] == "run" {
					return test.runResult
				}
				if args[0] == "rm" {
					cleanupCalls++
					return test.cleanup
				}
				t.Fatalf("unexpected command %v", args)
				return CommandResult{}
			}}
			instance := testDocker(fake)
			level := mustExercise(t, "zone01/29")
			result := instance.Run(context.Background(), level, "func CountAlpha(input string) int { return 10 }", []string{"l29-02"})
			if cleanupCalls != 1 {
				t.Fatalf("got %d cleanup calls", cleanupCalls)
			}
			if result.FailureKind != test.wantFailure {
				t.Fatalf("got failure %q, want %q: %#v", result.FailureKind, test.wantFailure, result)
			}
		})
	}
}

func TestDockerRunCancellationForcesCleanup(t *testing.T) {
	started := make(chan struct{})
	cleanup := make(chan struct{})
	fake := &fakeExecutor{run: func(ctx context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
		if args[0] == "run" {
			close(started)
			<-ctx.Done()
			return CommandResult{Err: ctx.Err()}
		}
		close(cleanup)
		return CommandResult{}
	}}
	instance := testDocker(fake)
	level := mustExercise(t, "zone01/29")
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan RunResult, 1)
	go func() {
		done <- instance.Run(ctx, level, "func CountAlpha(input string) int { for {} }", []string{"l29-02"})
	}()
	<-started
	cancel()
	result := <-done
	if !result.Stopped {
		t.Fatalf("expected stopped result, got %#v", result)
	}
	select {
	case <-cleanup:
	default:
		t.Fatal("container cleanup did not run after cancellation")
	}
}

func TestDockerMapsSandboxTimeoutAndRejectsMalformedResponse(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	for _, test := range []struct {
		name      string
		response  string
		wantTimed bool
		wantKind  string
	}{
		{
			name: "runtime timeout",
			response: mustJSON(t, sandboxprotocol.Response{
				Status: sandboxprotocol.StatusRuntimeTimeout, RuntimeError: "timed out", Results: []sandboxprotocol.TestResult{},
			}),
			wantTimed: true,
		},
		{name: "malformed", response: `{"status":"future"}`, wantKind: FailureInternal},
		{name: "unknown JSON field", response: `{"status":"success","surprise":true}`, wantKind: FailureStartup},
	} {
		t.Run(test.name, func(t *testing.T) {
			fake := &fakeExecutor{run: func(_ context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
				if args[0] == "run" {
					return CommandResult{Stdout: test.response}
				}
				return CommandResult{Err: errors.New("gone"), Stderr: "No such container"}
			}}
			result := testDocker(fake).Run(
				context.Background(), level,
				"func CountAlpha(input string) int { return 0 }",
				[]string{"l29-02"},
			)
			if result.TimedOut != test.wantTimed || result.FailureKind != test.wantKind {
				t.Fatalf("unexpected result %#v", result)
			}
		})
	}
}

func testDocker(commands CommandExecutor) *Engine {
	return newEngine(&dockerAdapter{
		dockerBinary: "docker",
		image:        "runner:test",
		commands:     commands,
		random:       bytes.NewReader(bytes.Repeat([]byte{7}, 128)),
	}, 1, nil)
}

func sandboxSuccessResponse(t *testing.T) CommandResult {
	t.Helper()
	return CommandResult{Stdout: mustJSON(t, sandboxprotocol.Response{
		Status:          sandboxprotocol.StatusSuccess,
		FormattedSource: "package main\n\nfunc CountAlpha(input string) int { return 10 }\n",
		Results: []sandboxprotocol.TestResult{{
			ID: "l29-02", Actual: "10", DurationMS: 0.1,
		}},
	})}
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

func makeImageFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, directory := range []string{
		filepath.Join(root, "docker"),
		filepath.Join(root, "cmd", "sandbox-runner"),
		filepath.Join(root, "internal", "sandboxprotocol"),
	} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{
		".dockerignore": ".git\n",
		"go.mod":        "module example.test/sandbox\ngo 1.23\n",
		filepath.Join("docker", "runner.Dockerfile"):         "FROM scratch\n",
		filepath.Join("cmd", "sandbox-runner", "main.go"):    "package main\n",
		filepath.Join("internal", "sandboxprotocol", "p.go"): "package sandboxprotocol\n",
	}
	for relative, content := range files {
		if err := os.WriteFile(filepath.Join(root, relative), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestDockerRunnerAtCapacity(t *testing.T) {
	block := make(chan struct{})
	started := make(chan struct{})
	fake := &fakeExecutor{run: func(_ context.Context, _ int, _ []byte, _ string, args ...string) CommandResult {
		if args[0] == "run" {
			close(started)
			<-block
			return sandboxSuccessResponse(t)
		}
		return CommandResult{Err: errors.New("gone"), Stderr: "No such container"}
	}}
	instance := testDocker(fake)
	level := mustExercise(t, "zone01/29")
	done := make(chan struct{})
	go func() {
		instance.Run(context.Background(), level, "func CountAlpha(input string) int { return 10 }", []string{"l29-02"})
		close(done)
	}()
	<-started
	result := instance.Run(context.Background(), level, "func CountAlpha(input string) int { return 0 }", []string{"l29-02"})
	if result.FailureKind != FailureCapacity {
		t.Fatalf("expected capacity result, got %#v", result)
	}
	close(block)
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("first run did not finish")
	}
}
