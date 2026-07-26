package runner

import (
	"context"
	"fmt"
	"strings"
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

func (*outcomeAdapter) Info() Info {
	return Info{Mode: ModeLocal}
}

func TestEngineOwnsPreparationVerdictsAndCompletion(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	wires := make([]wireResult, 0, len(level.Tests))
	for _, test := range level.Tests {
		wires = append(wires, wireResult{ID: test.ID, Actual: test.Expected})
	}
	adapter := &outcomeAdapter{outcome: executionOutcome{
		status: executionSuccess, results: wires, formattedSource: "package main\n\nfunc CountAlpha(input string) int { return 0 }\n",
	}}
	result := newEngine(adapter, 1, testIssuer{}).Run(
		context.Background(),
		level,
		"func CountAlpha(input string) int { return 0 }",
		nil,
	)
	if !result.Passed || result.Receipt == "" || result.ExerciseKey != level.Key {
		t.Fatalf("engine did not complete the passing run: %#v", result)
	}
	if adapter.plan.harness == "" || adapter.plan.sourceHash == "" {
		t.Fatalf("adapter received an incomplete execution plan: %#v", adapter.plan)
	}

	adapter.outcome = executionOutcome{status: executionCompile, compileError: "/tmp/solution.go:3: broken"}
	result = newEngine(adapter, 1, nil).Run(
		context.Background(),
		level,
		"func CountAlpha(input string) int { return 0 }",
		nil,
	)
	if result.CompileError != "solution.go:3: broken" || result.Results[0].Status != "compile" {
		t.Fatalf("engine did not own the compile verdict: %#v", result)
	}

	adapter.outcome = executionOutcome{status: executionCleanup, runtimeError: "cleanup failed"}
	result = newEngine(adapter, 1, nil).Run(
		context.Background(),
		level,
		"func CountAlpha(input string) int { return 0 }",
		nil,
	)
	if result.FailureKind != FailureCleanup || result.Passed {
		t.Fatalf("engine did not own the cleanup verdict: %#v", result)
	}
}

func TestEngineTreatsAReorderedCompleteSuiteAsTheWholeSuite(t *testing.T) {
	level := mustExercise(t, "foundation/1")
	testIDs := make([]string, 0, len(level.Tests))
	wires := make([]wireResult, 0, len(level.Tests))
	for index := len(level.Tests) - 1; index >= 0; index-- {
		testIDs = append(testIDs, level.Tests[index].ID)
		wires = append(wires, wireResult{
			ID: level.Tests[index].ID, Actual: level.Tests[index].Expected,
		})
	}
	adapter := &outcomeAdapter{outcome: executionOutcome{
		status: executionSuccess, results: wires,
	}}

	result := newEngine(adapter, 1, testIssuer{}).Run(
		context.Background(),
		level,
		"func Echo(value string) string { return value }",
		testIDs,
	)

	if !result.Passed || result.Receipt == "" {
		t.Fatalf("reordered complete suite was not treated as complete: %#v", result)
	}
	for index, want := range testIDs {
		if adapter.plan.tests[index].ID != want || result.Results[index].ID != want {
			t.Fatalf("test %d did not preserve requested order %q", index, want)
		}
	}
}

func TestEngineSelectsOutputForTheFirstFailureOrLastPassingTest(t *testing.T) {
	level := mustExercise(t, "foundation/1")
	passing := make([]wireResult, 0, len(level.Tests))
	for index, test := range level.Tests {
		passing = append(passing, wireResult{
			ID: test.ID, Actual: test.Expected, Stdout: fmt.Sprintf("output-%d", index+1),
		})
	}

	for _, test := range []struct {
		name    string
		results []wireResult
		want    string
	}{
		{
			name: "first failure",
			results: func() []wireResult {
				items := append([]wireResult(nil), passing...)
				items[1].Actual = `"wrong"`
				return items
			}(),
			want: "output-2",
		},
		{name: "last passing test", results: passing, want: "output-5"},
	} {
		t.Run(test.name, func(t *testing.T) {
			adapter := &outcomeAdapter{outcome: executionOutcome{
				status: executionSuccess, stdout: "all outputs", results: test.results,
			}}
			result := newEngine(adapter, 1, nil).Run(
				context.Background(),
				level,
				"func Echo(value string) string { return value }",
				nil,
			)
			if result.Stdout != test.want {
				t.Fatalf("stdout = %q, want %q", result.Stdout, test.want)
			}
		})
	}
}

func TestRunnerAcceptsCorrectAndRejectsIncorrectCode(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	localRunner := New("go", 1, testIssuer{})
	correct := `func CountAlpha(input string) int {
	count := 0
	for _, value := range input {
		if (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z') {
			count++
		}
	}
	return count
}`
	result := localRunner.Run(context.Background(), level, correct, nil)
	if !result.Passed || result.PassedCount != len(level.Tests) || result.Receipt == "" {
		t.Fatalf("correct solution failed: %#v", result)
	}

	incorrect := `func CountAlpha(input string) int { return 0 }`
	result = localRunner.Run(context.Background(), level, incorrect, nil)
	if result.Passed || result.PassedCount == len(level.Tests) {
		t.Fatal("incorrect solution passed")
	}
}

func TestRunnerReportsCompilerErrors(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`func CountAlpha(input string) int { definitely not go }`,
		nil,
	)
	if result.CompileError == "" {
		t.Fatalf("expected compiler error, got %#v", result)
	}
}

func TestRunnerExplainsGoDebugOutputAfterConsoleLogError(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`func CountAlpha(input string) int {
			console.log(input)
			return 0
		}`,
		[]string{"l29-02"},
	)
	if !strings.Contains(result.CompileError, "Go has no console.log") ||
		!strings.Contains(result.CompileError, "fmt.Println") {
		t.Fatalf("expected Go-specific debug guidance, got %q", result.CompileError)
	}
}

func TestRunnerCapturesStudentStandardOutput(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`import "fmt"
		func CountAlpha(input string) int {
			fmt.Println("debug:", input)
			return 0
		}`,
		[]string{"l29-01"},
	)
	if !strings.Contains(result.Stdout, "debug:") {
		t.Fatalf("expected student stdout, got %#v", result)
	}
}

func TestRunnerTerminatesTimeout(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`func CountAlpha(input string) int { for {} }`,
		[]string{"l29-02"},
	)
	if !result.TimedOut {
		t.Fatalf("expected timeout, got %#v", result)
	}
}

func TestEveryStarterCompilesWithItsHarness(t *testing.T) {
	localRunner := New("go", 1, nil)
	for _, level := range assessment.Levels() {
		result := localRunner.Run(
			context.Background(),
			level,
			level.StarterCode,
			[]string{level.Tests[0].ID},
		)
		if result.CompileError != "" {
			t.Fatalf("level %d starter did not compile: %s", level.ID, result.CompileError)
		}
	}
}

func TestFoundationalReferenceSolutionsPass(t *testing.T) {
	solutions := map[int]string{
		1:  `func Echo(value string) string { return value }`,
		2:  `func Increment(value int) int { return value + 1 }`,
		3:  `func IsPositive(value int) bool { return value > 0 }`,
		4:  `func MaxInt(left, right int) int { if left > right { return left }; return right }`,
		5:  `func Abs(value int) int { if value < 0 { return -value }; return value }`,
		6:  `func Clamp(value, minimum, maximum int) int { if value < minimum { return minimum }; if value > maximum { return maximum }; return value }`,
		7:  `func ByteCount(value string) int { return len(value) }`,
		8:  `func FirstByte(value string) string { if len(value) == 0 { return "" }; return string(value[0]) }`,
		9:  `func Repeat(value string, count int) string { result := ""; for i := 0; i < count; i++ { result += value }; return result }`,
		10: `func Sum(values []int) int { total := 0; for _, value := range values { total += value }; return total }`,
		11: `func CountValue(values []int, target int) int { count := 0; for _, value := range values { if value == target { count++ } }; return count }`,
		12: `func Contains(values []int, target int) bool { for _, value := range values { if value == target { return true } }; return false }`,
		13: `func Reverse(value string) string { result := ""; for i := len(value)-1; i >= 0; i-- { result += string(value[i]) }; return result }`,
		14: `func IsPalindrome(value string) bool { for left, right := 0, len(value)-1; left < right; left, right = left+1, right-1 { if value[left] != value[right] { return false } }; return true }`,
		15: `func FilterEven(values []int) []int { result := []int{}; for _, value := range values { if value%2 == 0 { result = append(result, value) } }; return result }`,
		16: `func Minimum(values []int) int { if len(values) == 0 { return 0 }; result := values[0]; for _, value := range values[1:] { if value < result { result = value } }; return result }`,
		17: `func UniqueStrings(values []string) []string { result := []string{}; seen := make(map[string]bool); for _, value := range values { if !seen[value] { seen[value] = true; result = append(result, value) } }; return result }`,
		18: `func Lengths(values []string) []int { result := []int{}; for _, value := range values { result = append(result, len(value)) }; return result }`,
		19: `func Frequencies(values []string) map[string]int { result := map[string]int{}; for _, value := range values { result[value]++ }; return result }`,
		20: `func RotateLeft(values []int, steps int) []int { if len(values) == 0 { return []int{} }; steps %= len(values); result := []int{}; result = append(result, values[steps:]...); result = append(result, values[:steps]...); return result }`,
		21: `func BalancedBrackets(value string) bool { stack := []byte{}; pairs := map[byte]byte{')':'(', ']':'[', '}':'{'}; for i := 0; i < len(value); i++ { char := value[i]; if char == '(' || char == '[' || char == '{' { stack = append(stack, char); continue }; if opening, closing := pairs[char]; closing { if len(stack) == 0 || stack[len(stack)-1] != opening { return false }; stack = stack[:len(stack)-1] } }; return len(stack) == 0 }`,
	}
	localRunner := New("go", 2, nil)
	for levelID := 1; levelID <= len(solutions); levelID++ {
		key := assessment.ExerciseKey(fmt.Sprintf("foundation/%d", levelID))
		level, found := assessment.FindExercise(key)
		if !found {
			t.Fatalf("foundational exercise %q is missing", key)
		}
		result := localRunner.Run(
			context.Background(),
			level,
			solutions[levelID],
			nil,
		)
		if !result.Passed {
			t.Fatalf(
				"exercise %d reference failed (%d/%d): compile=%q runtime=%q results=%#v",
				levelID,
				result.PassedCount,
				result.TotalCount,
				result.CompileError,
				result.RuntimeError,
				result.Results,
			)
		}
	}
}

func TestLocalRunnerEnforcesOutputLimit(t *testing.T) {
	level := mustExercise(t, "zone01/29")
	result := NewLocal("go", 1, nil).Run(
		context.Background(),
		level,
		`func CountAlpha(input string) int {
			for i := 0; i < 300000; i++ { print("output") }
			return 0
		}`,
		[]string{"l29-02"},
	)
	if result.FailureKind != FailureOutput {
		t.Fatalf("expected output failure, got %#v", result)
	}
	if len(result.Stderr) > MaxOutputBytes {
		t.Fatalf("stderr was not capped: %d bytes", len(result.Stderr))
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

func TestFormatSourceKeepsEditablePackageDeclaration(t *testing.T) {
	formatted, err := FormatSource("\ufeff package main\n\nfunc solve(){ }")
	if err != nil {
		t.Fatal(err)
	}
	if formatted != "package main\n\nfunc solve() {}\n" {
		t.Fatalf("unexpected formatted source %q", formatted)
	}
}
