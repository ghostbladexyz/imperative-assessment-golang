package assessment

import (
	"strings"
	"testing"
)

func TestAuthoringCompilerDerivesExecutableExercise(t *testing.T) {
	t.Parallel()
	spec := practiceSpec{
		id: 7, title: "Double", topic: "Arithmetic", difficulty: "Beginner",
		signature: "Double(value int) int",
		objective: "Double an integer.", input: "One integer.", output: "Twice the integer.",
		constraints: []string{"Return the result."},
		hints:       []string{"Add the value to itself.", "Return the expression.", "Keep the signature."},
		pitfalls:    []string{"Printing instead of returning."},
		starter:     "func Double(value int) int {\n\treturn 0\n}\n",
		inputFields: intField,
		call:        "Double(current.Payload.Value)",
		cases: []practiceCase{
			pc("Positive", "2", "4", map[string]any{"value": 2}),
		},
	}

	exercise := compilePracticeExercise(sourceFoundation, spec)

	if exercise.definitionErr != nil {
		t.Fatal(exercise.definitionErr)
	}
	if exercise.Key != "foundation/7" || exercise.Tests[0].ID != "l7-01" {
		t.Fatalf("unexpected compiled identity: %#v", exercise)
	}
	if len(exercise.Instructions.Examples) != 1 {
		t.Fatalf("got %d examples, want 1 genuine example", len(exercise.Instructions.Examples))
	}
	harness := exercise.BuildHarness(exercise.Tests)
	for _, required := range []string{"type inputCase struct", spec.call, exercise.Tests[0].ID} {
		if !strings.Contains(harness, required) {
			t.Fatalf("harness is missing %q", required)
		}
	}
}

// TestExerciseCompilerRemovesDuplicateCases verifies duplicate assertions do not reach either the test list or teaching examples.
func TestExerciseCompilerRemovesDuplicateCases(t *testing.T) {
	t.Parallel()
	spec := practiceSpec{
		id: 22, title: "Classify", topic: "Conditions", difficulty: "Beginner",
		signature: "Classify(value int) bool",
		objective: "Classify an integer.", input: "One integer.", output: "A classification.",
		starter:     "func Classify(value int) bool {\n\treturn false\n}\n",
		inputFields: intField,
		call:        "Classify(current.Payload.Value)",
		cases: []practiceCase{
			pc("Negative", "-1", "false", map[string]any{"value": -1}),
			pc("Zero", "0", "true", map[string]any{"value": 0}),
			pc("Positive", "1", "true", map[string]any{"value": 1}),
			pc("Duplicate negative", "-1", "false", map[string]any{"value": -1}),
		},
	}

	exercise := compilePracticeExercise(sourcePiscine, spec)

	if len(exercise.Tests) != 3 {
		t.Fatalf("got %d tests, want 3 unique tests", len(exercise.Tests))
	}
	if exercise.Tests[0].Name != "Negative" || exercise.Tests[2].Name != "Positive" {
		t.Fatalf("first occurrence order was not preserved: %#v", exercise.Tests)
	}
	if len(exercise.Instructions.Examples) != 3 {
		t.Fatalf("got %d examples, want 3 unique examples", len(exercise.Instructions.Examples))
	}
}

func TestAuthoringCompilerRejectsSynchronizedFieldDrift(t *testing.T) {
	t.Parallel()
	spec := practiceSpec{
		id: 1, title: "Broken", topic: "Functions", difficulty: "Beginner",
		signature: "Expected() int",
		objective: "Return a value.", input: "No input.", output: "One integer.",
		starter: "func Different() int { return 0 }\n",
		call:    "Expected()",
		cases:   []practiceCase{pc("Case", "", "0", map[string]any{})},
	}
	if err := validatePracticeSpec(sourceFoundation, spec); err == nil ||
		!strings.Contains(err.Error(), "starter does not implement signature") {
		t.Fatalf("unexpected authoring validation error: %v", err)
	}
}

func TestPrintCombinationUsesItsOriginalPrintedOutputContract(t *testing.T) {
	t.Parallel()

	exercise, found := FindExercise("piscine/1014")
	if !found {
		t.Fatal("Print Combination exercise is missing")
	}
	if exercise.Signature != "PrintComb()" {
		t.Fatalf("signature = %q, want PrintComb()", exercise.Signature)
	}
	if !strings.Contains(exercise.StarterCode, `"github.com/01-edu/z01"`) ||
		!strings.Contains(exercise.StarterCode, "z01.PrintRune") {
		t.Fatal("starter does not teach the required z01.PrintRune output contract")
	}
	if !strings.Contains(exercise.BuildHarness(exercise.Tests), "assessmentRunPrinted") {
		t.Fatal("harness does not assess printed output")
	}
}
