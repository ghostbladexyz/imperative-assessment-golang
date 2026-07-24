package assessment

import (
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
		if len(level.Instructions.Examples) < 4 {
			t.Fatalf("exercise %q has fewer than four examples", level.Key)
		}
		if len(level.Instructions.AllowedBuiltins) == 0 ||
			len(level.Instructions.AllowedPackages) == 0 {
			t.Fatalf("exercise %q has no explicit source policy", level.Key)
		}
		if harness := level.BuildHarness(level.Tests[:1]); harness == "" {
			t.Fatalf("exercise %q generated an empty harness", level.Key)
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
