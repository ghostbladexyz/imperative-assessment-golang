package assessment

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDefinitionsAreComplete(t *testing.T) {
	t.Parallel()
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
	for _, level := range Levels() {
		if level.Key == "" || level.Title == "" || level.Signature == "" || level.StarterCode == "" {
			t.Fatalf("exercise %q is missing core metadata", level.Key)
		}
		if !strings.HasPrefix(level.StarterCode, "package main\n") {
			t.Fatalf("exercise %q starter code must include an editable package declaration", level.Key)
		}
		if len(level.Instructions.Hints) < 3 {
			t.Fatalf("exercise %q has fewer than three hints", level.Key)
		}
		if len(level.Instructions.Examples) == 0 || len(level.Instructions.Examples) > maxPracticeExamples {
			t.Fatalf("exercise %q has %d examples, want between 1 and %d", level.Key, len(level.Instructions.Examples), maxPracticeExamples)
		}
		exerciseCallPrefix := functionName(level.Signature) + "("
		seenExamples := make(map[Example]struct{}, len(level.Instructions.Examples))
		for _, example := range level.Instructions.Examples {
			if !strings.HasPrefix(example.Input, exerciseCallPrefix) {
				t.Fatalf("exercise %q includes an example from another function: %#v", level.Key, example)
			}
			if _, exists := seenExamples[example]; exists {
				t.Fatalf("exercise %q repeats example %#v", level.Key, example)
			}
			seenExamples[example] = struct{}{}
		}
		if (len(level.Instructions.AllowedBuiltins) == 0 && level.answerMode != printedAnswer) ||
			len(level.Instructions.AllowedPackages) == 0 {
			t.Fatalf("exercise %q has no explicit source policy", level.Key)
		}
		if harness := level.BuildHarness(level.Tests[:1]); harness == "" {
			t.Fatalf("exercise %q generated an empty harness", level.Key)
		}
	}
}

// TestEveryExerciseHasCleanTestCoverage audits duplicate and synthetic rows across the complete catalogue.
func TestEveryExerciseHasCleanTestCoverage(t *testing.T) {
	t.Parallel()
	for _, level := range Levels() {
		seenNames := make(map[string]struct{}, len(level.Tests))
		seenRows := make(map[string]struct{}, len(level.Tests))
		for _, current := range level.Tests {
			if strings.Contains(current.Name, "stability pass") ||
				strings.Contains(current.Purpose, "stability pass") ||
				strings.Contains(current.Input, "stability pass") {
				t.Errorf("exercise %q contains a synthetic stability row %#v", level.Key, current)
			}
			if _, exists := seenNames[current.Name]; exists {
				t.Errorf("exercise %q repeats test name %q", level.Key, current.Name)
			}
			seenNames[current.Name] = struct{}{}
			row := current.Input + "\x00" + current.Expected
			if _, exists := seenRows[row]; exists {
				t.Errorf("exercise %q repeats test input and expectation %q", level.Key, current.Input)
			}
			seenRows[row] = struct{}{}
		}
	}
}

func TestPublicLevelsDoNotExposeHarnesses(t *testing.T) {
	t.Parallel()
	for _, level := range PublicLevels() {
		if level.build != nil {
			t.Fatalf("level %d exposed its harness builder", level.ID)
		}
	}
}

func TestPublicLevelsSerializeFrontendAllowlistsAsArrays(t *testing.T) {
	t.Parallel()
	for _, level := range PublicLevels() {
		if level.Instructions.AllowedBuiltins == nil {
			t.Errorf("exercise %q serializes allowedBuiltins as null", level.Key)
		}
		if level.Instructions.AllowedPackages == nil {
			t.Errorf("exercise %q serializes allowedPackages as null", level.Key)
		}
	}
}

func TestSelectTestsPreservesRequestedExecutionOrder(t *testing.T) {
	t.Parallel()
	level := Levels()[0]
	first, second, third := level.Tests[0], level.Tests[1], level.Tests[2]

	selected, valid := SelectTests(level, []string{third.ID, first.ID, second.ID})

	if !valid {
		t.Fatal("expected known test identifiers to be accepted")
	}
	for index, want := range []string{third.ID, first.ID, second.ID} {
		if selected[index].ID != want {
			t.Fatalf("selected test %d = %q, want %q", index, selected[index].ID, want)
		}
	}
}

func TestZone01CatalogueRemainsComplete(t *testing.T) {
	t.Parallel()
	for sourceID := 22; sourceID <= 65; sourceID++ {
		key := exerciseKey(sourceZone01, sourceID)
		if _, found := FindExercise(key); !found {
			t.Errorf("missing Zone01 exercise %q", key)
		}
	}
}

func TestProgressionStartsSimpleAndContainsFullCatalogue(t *testing.T) {
	t.Parallel()
	levels := Levels()
	if len(levels) != exerciseCount {
		t.Fatalf("got %d exercises, want %d", len(levels), exerciseCount)
	}
	wantSignatures := []string{
		"Echo(value string) string",
		"Increment(value int) int",
		"IsPositive(value int) bool",
	}
	for index, want := range wantSignatures {
		if levels[index].Signature != want {
			t.Errorf("exercise %d signature = %q, want %q", index+1, levels[index].Signature, want)
		}
		if levels[index].Difficulty != "Beginner" {
			t.Errorf("exercise %d difficulty = %q, want Beginner", index+1, levels[index].Difficulty)
		}
	}
	if levels[21].Title != "Only Z" {
		t.Fatalf("exercise 22 = %q, want Only Z", levels[21].Title)
	}
}

func TestImportedExercisesIncreaseInDifficulty(t *testing.T) {
	t.Parallel()
	levels := Levels()
	previous := 0
	for _, level := range levels[21:] {
		current := level.order
		if current < previous {
			t.Fatalf("%q (%s) appears after a harder exercise", level.Title, level.Difficulty)
		}
		previous = current
	}
}

func TestCatalogueLookupUsesStableKeysAndFrozenLegacyPositions(t *testing.T) {
	t.Parallel()
	levels := Levels()
	legacy := LegacyExerciseKeys()
	if len(legacy) != len(levels) {
		t.Fatalf("got %d legacy keys, want %d", len(legacy), len(levels))
	}
	for index, level := range levels {
		byKey, found := FindExercise(level.Key)
		if !found || byKey.ID != level.ID {
			t.Fatalf("key lookup for %q returned %#v", level.Key, byKey)
		}
		byPosition, found := FindLevel(index + 1)
		if !found || byPosition.Key != level.Key {
			t.Fatalf("position lookup %d returned %#v", index+1, byPosition)
		}
	}
	for position, want := range map[int]ExerciseKey{
		1: "foundation/1", 22: "piscine/1001", 40: "zone01/22",
		171: "zone01/65",
	} {
		key, found := LegacyExerciseKey(position)
		if !found || key != want {
			t.Fatalf("legacy position %d maps to %q, want %q", position, key, want)
		}
	}
}

func TestKinoz01CatalogueAddsOneHundredAndSixUniqueExercises(t *testing.T) {
	t.Parallel()
	if got := len(piscineLevels()); got != 106 {
		t.Fatalf("got %d kinoz01 exercises, want 106", got)
	}
}

func TestSourceAuditPreservesPrintedOutputExercises(t *testing.T) {
	t.Parallel()
	wanted := map[ExerciseKey]bool{
		"piscine/1001": true, "piscine/1002": true, "piscine/1003": true,
		"piscine/1004": true, "piscine/1005": true, "piscine/1008": true,
		"piscine/1009": true, "piscine/1010": true, "piscine/1011": true,
		"piscine/1012": true, "piscine/1013": true, "piscine/1014": true,
		"piscine/1015": true, "piscine/1016": true, "piscine/1017": true,
		"piscine/1018": true, "piscine/1050": true, "piscine/1051": true,
		"piscine/1060": true, "piscine/1061": true, "piscine/1062": true,
		"piscine/1063": true, "piscine/1064": true, "piscine/1072": true,
		"piscine/1079": true, "piscine/1086": true,
		"zone01/22": true, "zone01/34": true, "zone01/48": true,
		"zone01/49": true, "zone01/53": true,
	}

	for _, level := range Levels() {
		_, shouldPrint := wanted[level.Key]
		if (level.answerMode == printedAnswer) != shouldPrint {
			t.Errorf("exercise %q printed mode = %t, want %t", level.Key, level.answerMode == printedAnswer, shouldPrint)
			continue
		}
		if !shouldPrint {
			continue
		}
		delete(wanted, level.Key)
		if len(level.Instructions.AllowedPackages) != 1 ||
			level.Instructions.AllowedPackages[0] != z01Package {
			t.Errorf("exercise %q packages = %v, want only %q", level.Key, level.Instructions.AllowedPackages, z01Package)
		}
		for _, builtin := range level.Instructions.AllowedBuiltins {
			if builtin == "print" || builtin == "println" {
				t.Errorf("exercise %q allows %q instead of requiring z01.PrintRune", level.Key, builtin)
			}
		}
		if !strings.Contains(level.StarterCode, `"github.com/01-edu/z01"`) ||
			!strings.Contains(level.StarterCode, "z01.PrintRune") {
			t.Errorf("exercise %q starter does not provide z01.PrintRune", level.Key)
		}
		closingParenthesis := strings.LastIndex(level.Signature, ")")
		if closingParenthesis < 0 ||
			strings.TrimSpace(level.Signature[closingParenthesis+1:]) != "" {
			t.Errorf("exercise %q signature still declares a return value: %q", level.Key, level.Signature)
		}
		if strings.Contains(level.StarterCode, "\n\treturn ") {
			t.Errorf("exercise %q starter still returns a value", level.Key)
		}
		for _, current := range level.Tests {
			var printed string
			if err := json.Unmarshal([]byte(current.Expected), &printed); err != nil {
				t.Errorf("exercise %q test %q expected value is not printed text: %q", level.Key, current.Name, current.Expected)
			}
		}
	}
	if len(wanted) != 0 {
		t.Fatalf("source-audited exercises are missing: %v", wanted)
	}
}

func TestCatalogueReturnsDefensiveProjections(t *testing.T) {
	t.Parallel()
	first := Levels()
	first[0].Title = "changed"
	first[0].Tests[0].Name = "changed"
	first[0].Instructions.Hints[0] = "changed"

	second := Levels()
	if second[0].Title == "changed" ||
		second[0].Tests[0].Name == "changed" ||
		second[0].Instructions.Hints[0] == "changed" {
		t.Fatal("caller mutation changed the cached catalogue")
	}
}
