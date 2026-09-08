package runner

import (
	"context"
	"testing"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

type testIssuer struct{}

func (testIssuer) Issue(_ assessment.ExerciseKey, _ string) (string, error) {
	return "verified-receipt", nil
}

type outcomeAdapter struct {
	outcome executionOutcome
	plan    executionPlan
}

func (adapter *outcomeAdapter) Execute(_ context.Context, plan executionPlan) executionOutcome {
	adapter.plan = plan
	return adapter.outcome
}

func (*outcomeAdapter) Info() Info { return Info{Mode: ModeDocker} }

func TestEngineOwnsPreparationAndCompletion(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	wires := make([]wireResult, 0, len(level.Tests))
	for _, current := range level.Tests {
		wires = append(wires, wireResult{ID: current.ID, Actual: "pass"})
	}
	adapter := &outcomeAdapter{outcome: executionOutcome{status: executionSuccess, results: wires}}
	result := newEngine(adapter, 1, testIssuer{}).Run(context.Background(), level, level.StarterCode, nil)
	if !result.Passed || result.Receipt == "" || adapter.plan.sourceHash == "" {
		t.Fatalf("engine did not complete the passing run: %#v", result)
	}
}

func TestReorderedCompleteSuiteCountsAsWholeSuite(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	ids := make([]string, 0, len(level.Tests))
	wires := make([]wireResult, 0, len(level.Tests))
	for index := len(level.Tests) - 1; index >= 0; index-- {
		ids = append(ids, level.Tests[index].ID)
		wires = append(wires, wireResult{ID: level.Tests[index].ID, Actual: "pass"})
	}
	adapter := &outcomeAdapter{outcome: executionOutcome{status: executionSuccess, results: wires}}
	result := newEngine(adapter, 1, testIssuer{}).Run(context.Background(), level, level.StarterCode, ids)
	if !result.Passed || result.Receipt == "" {
		t.Fatalf("reordered suite was not complete: %#v", result)
	}
}

func TestFormatSourceKeepsEditablePackageDeclaration(t *testing.T) {
	formatted, err := FormatSource("\ufeff package main\n\nfunc solve(){ }")
	if err != nil {
		t.Fatal(err)
	}
	if formatted != "package main\n\nfunc solve() {}\n" {
		t.Fatalf("unexpected formatted source %q", formatted)
	}
}

func mustExercise(t *testing.T, key assessment.ExerciseKey) assessment.Level {
	t.Helper()
	level, found := assessment.FindExercise(key)
	if !found {
		t.Fatalf("exercise %q is missing", key)
	}
	return level
}
