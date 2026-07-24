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
		if level.Title == "" || level.Signature == "" || level.StarterCode == "" {
			t.Fatalf("level %d is missing core metadata", level.ID)
		}
		if !strings.HasPrefix(level.StarterCode, "package main\n") {
			t.Fatalf("level %d starter code must include an editable package declaration", level.ID)
		}
		if len(level.Instructions.Hints) < 3 {
			t.Fatalf("level %d has fewer than three hints", level.ID)
		}
		if len(level.Instructions.Examples) < 4 {
			t.Fatalf("level %d has fewer than four examples", level.ID)
		}
		if len(level.Instructions.AllowedBuiltins) == 0 ||
			len(level.Instructions.AllowedPackages) == 0 {
			t.Fatalf("level %d has no explicit source policy", level.ID)
		}
		if harness := level.BuildHarness(level.Tests[:1]); harness == "" {
			t.Fatalf("level %d generated an empty harness", level.ID)
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
	levels := Levels()
	wantTitles := []string{
		"Only A", "Print If Not", "Print If", "Rectangle Perimeter",
		"Count Character", "Check Number", "Retain First Half", "Count Alpha",
		"First Word", "Last Word", "Fish And Chips", "Digit Length",
		"Search Replace", "Repeat Alpha", "Greatest Common Divisor",
		"Camel To Snake Case", "Hash Code", "Third Time Is A Charm", "From To",
		"Is Capitalized", "Find Previous Prime", "Integer To ASCII",
		"Clean String", "Expand String", "We Are Unique", "Zip String",
		"Print Reverse Combo", "Print Memory", "Concat Slice", "Save And Miss",
		"Hidden P", "Word Match", "Intersection", "Union", "Concat Alternate",
		"Chunk", "Reverse String Capitalization", "Can Jump", "Add Prime Sum",
		"Prime Factors", "Fifth And Skip", "Reverse Concat Alternate",
		"Not Decimal", "Slice",
	}
	available := make(map[string]bool, len(levels))
	for _, level := range levels {
		available[level.Title] = true
	}
	for _, want := range wantTitles {
		if !available[want] {
			t.Errorf("missing Zone01 exercise %q", want)
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
		current := importedDifficultyRank(level.Difficulty)
		if current < previous {
			t.Fatalf("%q (%s) appears after a harder exercise", level.Title, level.Difficulty)
		}
		previous = current
	}
}

func TestCatalogueHasNoDuplicateTitles(t *testing.T) {
	t.Parallel()
	seen := make(map[string]int, exerciseCount)
	for _, level := range Levels() {
		if previous, found := seen[level.Title]; found {
			t.Fatalf("duplicate title %q at exercises %d and %d", level.Title, previous, level.ID)
		}
		seen[level.Title] = level.ID
	}
}

func TestKinoz01CatalogueAddsOneHundredAndSixUniqueExercises(t *testing.T) {
	t.Parallel()
	if got := len(piscineLevels()); got != 106 {
		t.Fatalf("got %d kinoz01 exercises, want 106", got)
	}
}
