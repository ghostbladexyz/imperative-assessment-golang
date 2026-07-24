package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

type Local struct {
	goBinary  string
	goVersion string
	common    common
}

func NewLocal(goBinary string, maxConcurrent int, receipts ReceiptIssuer) *Local {
	if goBinary == "" {
		goBinary = "go"
	}
	version := "local Go toolchain"
	versionCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if output, _, err, _ := runCommand(versionCtx, 8*1024, goBinary, "version"); err == nil {
		version = strings.TrimSpace(output)
	}
	return &Local{
		goBinary:  goBinary,
		goVersion: version,
		common:    newCommon(maxConcurrent, receipts),
	}
}

// New preserves the original constructor for packages that explicitly depend on
// the trusted local runner.
func New(goBinary string, maxConcurrent int, receipts ReceiptIssuer) *Local {
	return NewLocal(goBinary, maxConcurrent, receipts)
}

func (local *Local) Info() Info {
	return Info{
		Mode:         ModeLocal,
		SandboxReady: false,
		GoVersion:    local.goVersion,
		Message:      "Local runner selected. Submitted code executes with the current user's permissions.",
	}
}

func (local *Local) Format(_ context.Context, source string) (string, error) {
	return FormatSource(source)
}

func (local *Local) Run(ctx context.Context, level assessment.Level, source string, testIDs []string) RunResult {
	started := time.Now()
	result, selected, prepared, valid := prepareRun(level, source, testIDs, started)
	if !valid || !local.common.acquire(ctx, &result, started) {
		return result
	}
	defer local.common.release()

	tempDir, err := os.MkdirTemp("", "imperative-assessment-*")
	if err != nil {
		result.RuntimeError = "Could not create an isolated temporary run directory."
		result.FailureKind = FailureInternal
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
		result.FailureKind = FailureInternal
		result.DurationMS = elapsedMS(started)
		return result
	}
	if err := os.WriteFile(harnessPath, []byte(level.BuildHarness(selected)), 0o600); err != nil {
		result.RuntimeError = "Could not prepare the controlled test harness."
		result.FailureKind = FailureInternal
		result.DurationMS = elapsedMS(started)
		return result
	}

	buildCtx, cancelBuild := context.WithTimeout(ctx, CompileTimeout)
	buildStdout, buildStderr, buildErr, buildLimited := runCommand(
		buildCtx, MaxOutputBytes, local.goBinary,
		"build", "-trimpath", "-o", executablePath, solutionPath, harnessPath,
	)
	cancelBuild()
	result.Stdout = buildStdout
	result.Stderr = buildStderr
	if buildErr != nil {
		switch {
		case errors.Is(buildCtx.Err(), context.DeadlineExceeded):
			result.TimedOut = true
			result.CompileError = "Compilation exceeded the 20 second limit."
		case errors.Is(ctx.Err(), context.Canceled):
			result.Stopped = true
		case buildLimited:
			result.FailureKind = FailureOutput
			result.CompileError = "Compiler output exceeded the 256 KiB limit."
		default:
			result.CompileError = cleanCompilerError(buildStderr)
		}
		markAll(result.Results, "compile", result.CompileError)
		result.DurationMS = elapsedMS(started)
		return result
	}

	runCtx, cancelRun := context.WithTimeout(ctx, RuntimeTimeout)
	runStdout, runStderr, runErr, runLimited := runCommand(runCtx, MaxOutputBytes, executablePath)
	cancelRun()
	result.Stdout = stripMarkers(runStdout)
	result.Stderr = runStderr
	applyWireResults(&result, parseResults(runStdout))
	if runErr != nil {
		switch {
		case errors.Is(runCtx.Err(), context.DeadlineExceeded):
			result.TimedOut = true
			result.RuntimeError = "Execution exceeded the 4 second limit and was terminated."
		case errors.Is(ctx.Err(), context.Canceled):
			result.Stopped = true
			result.RuntimeError = "Execution was stopped."
		case runLimited:
			result.FailureKind = FailureOutput
			result.RuntimeError = "Program output exceeded the 256 KiB limit and execution was terminated."
		default:
			result.RuntimeError = strings.TrimSpace(runStderr)
			if result.RuntimeError == "" {
				result.RuntimeError = "The program exited before completing every test."
			}
		}
	}
	local.common.issueReceipt(level, testIDs, &result)
	result.FormattedCode, _ = FormatSource(source)
	result.DurationMS = elapsedMS(started)
	return result
}
