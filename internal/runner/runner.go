package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/format"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

const (
	MaxSourceBytes = 192 * 1024
	MaxOutputBytes = 256 * 1024

	CompileTimeout = 20 * time.Second
	RuntimeTimeout = 4 * time.Second
)

type Mode string

const (
	ModeDocker Mode = "docker"
	ModeLocal  Mode = "local"
)

const (
	FailureCapacity = "capacity"
	FailureCleanup  = "cleanup"
	FailureInternal = "internal"
	FailureOutput   = "output"
	FailureStartup  = "startup"
)

type executionStatus string

const (
	executionSuccess        executionStatus = "success"
	executionCompile        executionStatus = "compile"
	executionCompileTimeout executionStatus = "compile_timeout"
	executionRuntime        executionStatus = "runtime"
	executionRuntimeTimeout executionStatus = "runtime_timeout"
	executionStopped        executionStatus = "stopped"
	executionOutput         executionStatus = "output"
	executionInternal       executionStatus = "internal"
	executionStartup        executionStatus = "startup"
	executionCleanup        executionStatus = "cleanup"
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
	FailureKind   string       `json:"failureKind,omitempty"`
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

type Info struct {
	Mode         Mode   `json:"runnerMode"`
	SandboxReady bool   `json:"sandboxReady"`
	GoVersion    string `json:"goVersion"`
	DockerImage  string `json:"dockerImage,omitempty"`
	Message      string `json:"message"`
}

type Service interface {
	Run(context.Context, assessment.Level, string, []string) RunResult
	Format(context.Context, string) (string, error)
	Info() Info
}

type ReceiptIssuer interface {
	Issue(int, string) (string, error)
}

type executionPlan struct {
	level      assessment.Level
	tests      []assessment.VisibleTest
	prepared   string
	harness    string
	sourceHash string
}

type executionOutcome struct {
	status          executionStatus
	compileError    string
	runtimeError    string
	failureKind     string
	stdout          string
	stderr          string
	formattedSource string
	results         []wireResult
}

type executionAdapter interface {
	Execute(context.Context, executionPlan) executionOutcome
	Info() Info
}

type Engine struct {
	adapter  executionAdapter
	slots    chan struct{}
	receipts ReceiptIssuer
}

func newEngine(adapter executionAdapter, maxConcurrent int, receipts ReceiptIssuer) *Engine {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	return &Engine{
		adapter:  adapter,
		slots:    make(chan struct{}, maxConcurrent),
		receipts: receipts,
	}
}

func (engine *Engine) Info() Info {
	return engine.adapter.Info()
}

func (engine *Engine) Format(_ context.Context, source string) (string, error) {
	return FormatSource(source)
}

func (engine *Engine) Run(ctx context.Context, level assessment.Level, source string, testIDs []string) RunResult {
	started := time.Now()
	result, plan, valid := prepareRun(level, source, testIDs, started)
	if !valid || !engine.acquire(ctx, &result, started) {
		return result
	}
	defer engine.release()

	outcome := engine.adapter.Execute(ctx, plan)
	return engine.complete(plan, outcome, result, len(testIDs) == 0, started)
}

func (engine *Engine) acquire(ctx context.Context, result *RunResult, started time.Time) bool {
	select {
	case engine.slots <- struct{}{}:
		return true
	case <-ctx.Done():
		result.Stopped = true
		result.DurationMS = elapsedMS(started)
		return false
	default:
		result.FailureKind = FailureCapacity
		result.RuntimeError = "The runner is temporarily at capacity. Wait for the current run to finish, then try again."
		result.DurationMS = elapsedMS(started)
		return false
	}
}

func (engine *Engine) release() {
	<-engine.slots
}

func (engine *Engine) complete(
	plan executionPlan,
	outcome executionOutcome,
	result RunResult,
	wholeSuite bool,
	started time.Time,
) RunResult {
	result.Stdout = outcome.stdout
	result.Stderr = outcome.stderr
	result.CompileError = outcome.compileError
	result.RuntimeError = outcome.runtimeError
	result.FailureKind = outcome.failureKind
	applyWireResults(&result, outcome.results)

	switch outcome.status {
	case executionCompileTimeout:
		result.TimedOut = true
		if result.CompileError == "" {
			result.CompileError = "Compilation exceeded the 20 second limit."
		}
		markAll(result.Results, "compile", result.CompileError)
	case executionCompile:
		result.CompileError = cleanCompilerError(result.CompileError)
		markAll(result.Results, "compile", result.CompileError)
	case executionRuntimeTimeout:
		result.TimedOut = true
		if result.RuntimeError == "" {
			result.RuntimeError = "Execution exceeded the 4 second limit and was terminated."
		}
	case executionStopped:
		result.Stopped = true
		if result.RuntimeError == "" {
			result.RuntimeError = "Execution was stopped."
		}
	case executionOutput:
		result.FailureKind = FailureOutput
		if result.CompileError != "" {
			markAll(result.Results, "compile", result.CompileError)
		}
	case executionInternal:
		result.FailureKind = FailureInternal
	case executionStartup:
		result.FailureKind = FailureStartup
	case executionCleanup:
		result.FailureKind = FailureCleanup
	case executionRuntime, executionSuccess:
	}

	if outcome.formattedSource != "" {
		withoutPackage := packagePattern.ReplaceAllString(outcome.formattedSource, "")
		result.FormattedCode = strings.TrimSpace(withoutPackage) + "\n"
	} else if outcome.status != executionCompile && outcome.status != executionCompileTimeout {
		formatted, _ := FormatSource(plan.prepared)
		withoutPackage := packagePattern.ReplaceAllString(formatted, "")
		result.FormattedCode = strings.TrimSpace(withoutPackage) + "\n"
	}

	result.Passed = outcome.status == executionSuccess &&
		wholeSuite &&
		result.PassedCount == len(plan.level.Tests) &&
		result.CompileError == "" &&
		result.RuntimeError == "" &&
		result.FailureKind == "" &&
		!result.TimedOut &&
		!result.Stopped
	if result.Passed && engine.receipts != nil {
		result.Receipt, _ = engine.receipts.Issue(plan.level.ID, result.SourceHash)
	}
	result.DurationMS = elapsedMS(started)
	return result
}

func PrepareSource(source string) string {
	source = strings.TrimPrefix(source, "\uFEFF")
	source = packagePattern.ReplaceAllString(source, "")
	return "package main\n\n" + strings.TrimSpace(source) + "\n"
}

func FormatSource(source string) (string, error) {
	if len(source) > MaxSourceBytes {
		return "", fmt.Errorf("source exceeds %d bytes", MaxSourceBytes)
	}
	formatted, err := format.Source([]byte(PrepareSource(source)))
	if err != nil {
		return "", fmt.Errorf("%s", cleanCompilerError(err.Error()))
	}
	return string(formatted), nil
}

func sourceHash(prepared string) string {
	sum := sha256.Sum256([]byte(prepared))
	return hex.EncodeToString(sum[:])
}

func prepareRun(level assessment.Level, source string, testIDs []string, started time.Time) (RunResult, executionPlan, bool) {
	result := RunResult{LevelID: level.ID}
	if len(source) > MaxSourceBytes {
		result.CompileError = fmt.Sprintf("Source is too large. The limit is %d KiB.", MaxSourceBytes/1024)
		result.DurationMS = elapsedMS(started)
		return result, executionPlan{}, false
	}
	selected, valid := assessment.SelectTests(level, testIDs)
	if !valid {
		result.CompileError = "The request contains an unknown test identifier."
		result.DurationMS = elapsedMS(started)
		return result, executionPlan{}, false
	}
	result.TotalCount = len(selected)
	for _, current := range selected {
		result.Results = append(result.Results, TestResult{
			ID: current.ID, Name: current.Name, Purpose: current.Purpose,
			Input: current.Input, Expected: current.Expected, Status: "pending",
		})
	}
	prepared := PrepareSource(source)
	result.SourceHash = sourceHash(prepared)
	if err := assessment.ValidateSourcePolicy(level, prepared); err != nil {
		result.CompileError = err.Error()
		markAll(result.Results, "compile", result.CompileError)
		result.DurationMS = elapsedMS(started)
		return result, executionPlan{}, false
	}
	return result, executionPlan{
		level:      level,
		tests:      selected,
		prepared:   prepared,
		harness:    level.BuildHarness(selected),
		sourceHash: result.SourceHash,
	}, true
}

func applyWireResults(result *RunResult, items []wireResult) {
	byID := make(map[string]wireResult, len(items))
	for _, item := range items {
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
}

type wireResult struct {
	ID         string  `json:"id"`
	Actual     string  `json:"actual"`
	Failure    string  `json:"failure"`
	DurationMS float64 `json:"durationMs"`
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
	cleaned := strings.Join(lines, "\n")
	if strings.Contains(cleaned, "undefined: console") {
		cleaned += "\n\nGo has no console.log. Import \"fmt\" and use fmt.Println(...), or use print/println. Runtime output appears in the Console after the program compiles."
	}
	return cleaned
}

func elapsedMS(started time.Time) float64 {
	return float64(time.Since(started).Microseconds()) / 1000
}

func SortedStatuses(results []TestResult) []string {
	statuses := make([]string, 0, len(results))
	for _, result := range results {
		statuses = append(statuses, result.Status)
	}
	sort.Strings(statuses)
	return statuses
}
