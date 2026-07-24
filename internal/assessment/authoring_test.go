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
	if len(exercise.Instructions.Examples) != 4 {
		t.Fatalf("got %d examples, want 4", len(exercise.Instructions.Examples))
	}
	harness := exercise.BuildHarness(exercise.Tests)
	for _, required := range []string{"type inputCase struct", spec.call, exercise.Tests[0].ID} {
		if !strings.Contains(harness, required) {
			t.Fatalf("harness is missing %q", required)
		}
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
