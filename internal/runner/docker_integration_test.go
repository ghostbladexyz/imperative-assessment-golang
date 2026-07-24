package runner

import (
	"context"
	"encoding/json"
	"os"
	"slices"
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
	level := mustExercise(t, "zone01/29")
	isolationLevel := level
	isolationLevel.Instructions.AllowedPackages = append(
		append([]string(nil), level.Instructions.AllowedPackages...),
		"net",
		"os",
		"time",
	)

	t.Run("correct solution passes", func(t *testing.T) {
		result := sandbox.Run(context.Background(), level, correctCountAlphaSolution, nil)
		if !result.Passed || result.PassedCount != len(level.Tests) {
			t.Fatalf("correct solution failed: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("incorrect solution fails assertions", func(t *testing.T) {
		result := sandbox.Run(
			context.Background(),
			level,
			`func CountAlpha(input string) int { return 0 }`,
			[]string{"l29-02"},
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
			`func CountAlpha(input string) int { definitely not go }`,
			[]string{"l29-02"},
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
			`func CountAlpha(input string) int { for {} }`,
			[]string{"l29-02"},
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
				`func CountAlpha(input string) int { for {} }`,
				[]string{"l29-02"},
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

	t.Run("effective container configuration is isolated", func(t *testing.T) {
		assertNoAssessmentContainers(t)
		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan RunResult, 1)
		go func() {
			done <- sandbox.Run(
				ctx,
				level,
				`func CountAlpha(input string) int { for {} }`,
				[]string{"l29-02"},
			)
		}()

		name := waitForAssessmentContainer(t)
		inspection := inspectAssessmentContainer(t, name)
		if !inspection.State.Running {
			t.Fatal("sandbox container was not running while submission executed")
		}
		if inspection.Config.Image != sandbox.Info().DockerImage ||
			inspection.Config.User != "65532:65532" ||
			inspection.Path != "/sandbox-runner" {
			t.Fatalf("unexpected sandbox identity: %#v", inspection)
		}
		if inspection.HostConfig.NetworkMode != "none" ||
			inspection.HostConfig.IpcMode != "none" ||
			!inspection.HostConfig.ReadonlyRootfs ||
			inspection.HostConfig.LogConfig.Type != "none" {
			t.Fatalf("sandbox namespaces or filesystem are not isolated: %#v", inspection.HostConfig)
		}
		if !slices.Contains(inspection.HostConfig.CapDrop, "ALL") ||
			!slices.Contains(inspection.HostConfig.SecurityOpt, "no-new-privileges") {
			t.Fatalf("sandbox privilege restrictions are missing: %#v", inspection.HostConfig)
		}
		if inspection.HostConfig.Memory != 256*1024*1024 ||
			inspection.HostConfig.MemorySwap != 256*1024*1024 ||
			inspection.HostConfig.NanoCpus != 1_000_000_000 ||
			inspection.HostConfig.PidsLimit != 64 {
			t.Fatalf("sandbox resource limits are not effective: %#v", inspection.HostConfig)
		}
		if len(inspection.HostConfig.Binds) != 0 || len(inspection.HostConfig.Mounts) != 0 {
			t.Fatalf("sandbox unexpectedly has host mounts: %#v", inspection.HostConfig)
		}
		for _, required := range []string{"/tmp", "/tmp/go-build", "/workspace"} {
			if _, ok := inspection.HostConfig.Tmpfs[required]; !ok {
				t.Fatalf("sandbox is missing tmpfs %s: %#v", required, inspection.HostConfig.Tmpfs)
			}
		}
		for _, environment := range inspection.Config.Env {
			name, _, _ := strings.Cut(environment, "=")
			if !slices.Contains([]string{
				"CGO_ENABLED", "GOCACHE", "GOGC", "GOENV", "GOMAXPROCS",
				"GOMEMLIMIT", "GOPATH", "GOPROXY", "GOTELEMETRY", "GOTOOLCHAIN",
				"GOTMPDIR", "HOME", "PATH",
			}, name) {
				t.Fatalf("sandbox inherited unexpected environment variable %q", name)
			}
		}

		cancel()
		result := <-done
		if !result.Stopped {
			t.Fatalf("expected stopped result: %#v", result)
		}
		assertNoAssessmentContainers(t)
	})

	t.Run("startup removes a stale created container", func(t *testing.T) {
		assertNoAssessmentContainers(t)
		name := "imperative-go-assessment-ffeeddccbbaa998877665544"
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, _, _, _ = runCommand(ctx, 64*1024, "docker", "rm", "--force", name)
		})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, stderr, err, _ := runCommand(
			ctx,
			64*1024,
			"docker",
			"create",
			"--name",
			name,
			"--label",
			dockerContainerLabel,
			sandbox.Info().DockerImage,
		)
		cancel()
		if err != nil {
			t.Fatalf("create stale sandbox fixture: %v: %s", err, stderr)
		}
		if err := cleanupStaleContainers(context.Background(), commands, "docker"); err != nil {
			t.Fatal(err)
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
func CountAlpha(input string) int {
	connection, err := net.DialTimeout("tcp", "1.1.1.1:80", 500*time.Millisecond)
	if err == nil {
		connection.Close()
		return 0
	}
	return 10
}`,
			[]string{"l29-02"},
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
func CountAlpha(input string) int {
	if os.WriteFile("/root-write-probe", []byte("unsafe"), 0600) == nil {
		return 0
	}
	return 10
}`,
			[]string{"l29-02"},
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
func CountAlpha(input string) int {
	if _, err := os.ReadFile("/imperative-assessment-host-probe"); err == nil {
		return 0
	}
	return 10
}`,
			[]string{"l29-02"},
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
			`func CountAlpha(input string) int {
	for i := 0; i < 300000; i++ { print("output") }
	return 10
}`,
			[]string{"l29-02"},
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

type dockerInspection struct {
	Path   string
	Config struct {
		Env   []string
		Image string
		User  string
	}
	State struct {
		Running bool
	}
	HostConfig struct {
		Binds          []string
		CapDrop        []string
		IpcMode        string
		Memory         int64
		MemorySwap     int64
		Mounts         []struct{}
		NanoCpus       int64
		NetworkMode    string
		PidsLimit      int64
		ReadonlyRootfs bool
		SecurityOpt    []string
		Tmpfs          map[string]string
		LogConfig      struct {
			Type string
		}
	}
}

func waitForAssessmentContainer(t *testing.T) string {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		stdout, _, err, _ := runCommand(
			ctx,
			64*1024,
			"docker",
			"ps",
			"--filter",
			"label="+dockerContainerLabel,
			"--format",
			"{{.Names}}",
		)
		cancel()
		if err == nil && strings.TrimSpace(stdout) != "" {
			return strings.TrimSpace(stdout)
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("sandbox container did not start")
	return ""
}

func inspectAssessmentContainer(t *testing.T, name string) dockerInspection {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stdout, stderr, err, _ := runCommand(ctx, 1024*1024, "docker", "inspect", name)
	if err != nil {
		t.Fatalf("inspect sandbox container: %v: %s", err, stderr)
	}
	var inspections []dockerInspection
	if err := json.Unmarshal([]byte(stdout), &inspections); err != nil || len(inspections) != 1 {
		t.Fatalf("decode sandbox inspection: %v", err)
	}
	return inspections[0]
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

const correctCountAlphaSolution = `func CountAlpha(input string) int {
	count := 0
	for _, value := range input {
		if (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z') {
			count++
		}
	}
	return count
}`
