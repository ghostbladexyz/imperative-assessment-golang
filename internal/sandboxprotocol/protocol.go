package sandboxprotocol

import (
	"fmt"
	"time"
)

const (
	StatusSuccess        = "success"
	StatusAssertion      = "assertion"
	StatusCompile        = "compile"
	StatusCompileTimeout = "compile_timeout"
	StatusRuntime        = "runtime"
	StatusRuntimeTimeout = "runtime_timeout"
	StatusStopped        = "stopped"
	StatusOutput         = "output"
	StatusInternal       = "internal"
)

type ExpectedTest struct {
	ID       string `json:"id"`
	Expected string `json:"expected"`
}

type Request struct {
	Source           string         `json:"source"`
	Harness          string         `json:"harness"`
	Tests            []ExpectedTest `json:"tests"`
	CompileTimeoutMS int64          `json:"compileTimeoutMs"`
	RuntimeTimeoutMS int64          `json:"runtimeTimeoutMs"`
	OutputLimitBytes int            `json:"outputLimitBytes"`
}

func (request Request) Validate(maxSourceBytes int) error {
	if request.Source == "" || len(request.Source) > maxSourceBytes {
		return fmt.Errorf("source must contain between 1 and %d bytes", maxSourceBytes)
	}
	if request.Harness == "" || len(request.Harness) > maxSourceBytes {
		return fmt.Errorf("harness must contain between 1 and %d bytes", maxSourceBytes)
	}
	if len(request.Tests) == 0 || len(request.Tests) > 32 {
		return fmt.Errorf("request must contain between 1 and 32 tests")
	}
	seen := make(map[string]struct{}, len(request.Tests))
	for _, test := range request.Tests {
		if test.ID == "" || len(test.ID) > 80 {
			return fmt.Errorf("test identifier is invalid")
		}
		if _, exists := seen[test.ID]; exists {
			return fmt.Errorf("test identifiers must be unique")
		}
		seen[test.ID] = struct{}{}
	}
	if request.CompileTimeoutMS < 1 || request.CompileTimeoutMS > int64((30*time.Second)/time.Millisecond) {
		return fmt.Errorf("compile timeout is invalid")
	}
	if request.RuntimeTimeoutMS < 1 || request.RuntimeTimeoutMS > int64((10*time.Second)/time.Millisecond) {
		return fmt.Errorf("runtime timeout is invalid")
	}
	if request.OutputLimitBytes < 1024 || request.OutputLimitBytes > 256*1024 {
		return fmt.Errorf("output limit is invalid")
	}
	return nil
}

type TestResult struct {
	ID         string  `json:"id"`
	Actual     string  `json:"actual"`
	Failure    string  `json:"failure,omitempty"`
	Stdout     string  `json:"stdout,omitempty"`
	DurationMS float64 `json:"durationMs"`
}

type Response struct {
	Status          string       `json:"status"`
	CompileError    string       `json:"compileError,omitempty"`
	RuntimeError    string       `json:"runtimeError,omitempty"`
	Stdout          string       `json:"stdout"`
	Stderr          string       `json:"stderr"`
	FormattedSource string       `json:"formattedSource,omitempty"`
	Results         []TestResult `json:"results"`
}

func (response Response) Validate(expected []ExpectedTest) error {
	switch response.Status {
	case StatusSuccess, StatusAssertion, StatusCompile, StatusCompileTimeout,
		StatusRuntime, StatusRuntimeTimeout, StatusStopped, StatusOutput, StatusInternal:
	default:
		return fmt.Errorf("unknown sandbox response status")
	}
	expectedIDs := make(map[string]struct{}, len(expected))
	for _, test := range expected {
		expectedIDs[test.ID] = struct{}{}
	}
	seen := make(map[string]struct{}, len(response.Results))
	for _, result := range response.Results {
		if _, exists := expectedIDs[result.ID]; !exists {
			return fmt.Errorf("sandbox response contains an unknown test")
		}
		if _, exists := seen[result.ID]; exists {
			return fmt.Errorf("sandbox response contains a duplicate test")
		}
		seen[result.ID] = struct{}{}
	}
	return nil
}
