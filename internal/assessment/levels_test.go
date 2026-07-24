package assessment

import "testing"

func TestDefinitionsAreComplete(t *testing.T) {
	t.Parallel()
	if err := Validate(); err != nil {
		t.Fatal(err)
	}
	for _, level := range Levels() {
		if level.Title == "" || level.Signature == "" || level.StarterCode == "" {
			t.Fatalf("level %d is missing core metadata", level.ID)
		}
		if len(level.Instructions.Hints) < 3 {
			t.Fatalf("level %d has fewer than three hints", level.ID)
		}
		if len(level.Instructions.Examples) < 2 {
			t.Fatalf("level %d has fewer than two examples", level.ID)
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
