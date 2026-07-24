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

func TestTokenTrailTeachesGoSliceSyntax(t *testing.T) {
	t.Parallel()
	var level Level
	for _, candidate := range Levels() {
		if candidate.Title == "Token Trail" {
			level = candidate
			break
		}
	}
	if level.ID == 0 {
		t.Fatal("Token Trail level is missing")
	}
	if !strings.Contains(level.StarterCode, "result := []string{}") {
		t.Fatal("Token Trail starter must demonstrate a valid empty string slice")
	}
	pitfalls := strings.Join(level.Instructions.CommonPitfalls, "\n")
	for _, syntax := range []string{"[]string{}", "append(result, value)", "string(input[i])"} {
		if !strings.Contains(pitfalls, syntax) {
			t.Errorf("level 1 pitfalls do not explain %q", syntax)
		}
	}
}

func TestProgressionStartsSimpleAndContainsThirtyExercises(t *testing.T) {
	t.Parallel()
	levels := Levels()
	if len(levels) != 30 {
		t.Fatalf("got %d exercises, want 30", len(levels))
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
	if levels[21].Title != "Token Trail" {
		t.Fatalf("exercise 22 = %q, want Token Trail", levels[21].Title)
	}
}
