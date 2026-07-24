package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type localAdapter struct {
	goBinary  string
	goVersion string
}

func NewLocal(goBinary string, maxConcurrent int, receipts ReceiptIssuer) *Engine {
	if goBinary == "" {
		goBinary = "go"
	}
	version := "local Go toolchain"
	versionCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if output, _, err, _ := runCommand(versionCtx, 8*1024, goBinary, "version"); err == nil {
		version = strings.TrimSpace(output)
	}
	return newEngine(&localAdapter{
		goBinary:  goBinary,
		goVersion: version,
	}, maxConcurrent, receipts)
}

// New preserves the original constructor for packages that explicitly depend on
// the trusted local runner.
func New(goBinary string, maxConcurrent int, receipts ReceiptIssuer) *Engine {
	return NewLocal(goBinary, maxConcurrent, receipts)
}

func (local *localAdapter) Info() Info {
	return Info{
		Mode:         ModeLocal,
		SandboxReady: false,
		GoVersion:    local.goVersion,
		Message:      "Local runner selected. Submitted code executes with the current user's permissions.",
	}
}

func (local *localAdapter) Execute(ctx context.Context, plan executionPlan) executionOutcome {
	tempDir, err := os.MkdirTemp("", "imperative-assessment-*")
	if err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "Could not create an isolated temporary run directory.",
		}
	}
	defer os.RemoveAll(tempDir)

	solutionPath := filepath.Join(tempDir, "solution.go")
	harnessPath := filepath.Join(tempDir, "assessment_harness.go")
	executablePath := filepath.Join(tempDir, "assessment-runner")
	if runtime.GOOS == "windows" {
		executablePath += ".exe"
	}
	if err := os.WriteFile(solutionPath, []byte(plan.prepared), 0o600); err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "Could not prepare the submitted source.",
		}
	}
	if err := os.WriteFile(harnessPath, []byte(plan.harness), 0o600); err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "Could not prepare the controlled test harness.",
		}
	}

	buildCtx, cancelBuild := context.WithTimeout(ctx, CompileTimeout)
	buildStdout, buildStderr, buildErr, buildLimited := runCommand(
		buildCtx, MaxOutputBytes, local.goBinary,
		"build", "-trimpath", "-o", executablePath, solutionPath, harnessPath,
	)
	cancelBuild()
	if buildErr != nil {
		outcome := executionOutcome{status: executionCompile, stdout: buildStdout, stderr: buildStderr}
		switch {
		case errors.Is(buildCtx.Err(), context.DeadlineExceeded):
			outcome.status = executionCompileTimeout
		case errors.Is(ctx.Err(), context.Canceled):
			outcome.status = executionStopped
		case buildLimited:
			outcome.status = executionOutput
			outcome.compileError = "Compiler output exceeded the 256 KiB limit."
		default:
			outcome.compileError = buildStderr
		}
		return outcome
	}

	runCtx, cancelRun := context.WithTimeout(ctx, RuntimeTimeout)
	runStdout, runStderr, runErr, runLimited := runCommand(runCtx, MaxOutputBytes, executablePath)
	cancelRun()
	outcome := executionOutcome{
		status:  executionSuccess,
		stdout:  stripMarkers(runStdout),
		stderr:  runStderr,
		results: parseResults(runStdout),
	}
	if runErr != nil {
		switch {
		case errors.Is(runCtx.Err(), context.DeadlineExceeded):
			outcome.status = executionRuntimeTimeout
		case errors.Is(ctx.Err(), context.Canceled):
			outcome.status = executionStopped
		case runLimited:
			outcome.status = executionOutput
			outcome.runtimeError = "Program output exceeded the 256 KiB limit and execution was terminated."
		default:
			outcome.status = executionRuntime
			outcome.runtimeError = strings.TrimSpace(runStderr)
			if outcome.runtimeError == "" {
				outcome.runtimeError = "The program exited before completing every test."
			}
		}
	}
	return outcome
}
