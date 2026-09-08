package assessment

import (
	"fmt"
	"strings"
)

const exerciseCount = 17

type catalogueSnapshot struct {
	levels []Level
	byKey  map[ExerciseKey]int
	err    error
}

var currentCatalogue = buildCatalogue()

func Levels() []Level { return catalogueLevels() }

func Validate() error { return currentCatalogue.err }

func buildCatalogue() catalogueSnapshot {
	levels := checkpointLevels()
	snapshot := catalogueSnapshot{levels: levels, byKey: make(map[ExerciseKey]int, len(levels))}
	for index := range snapshot.levels {
		level := &snapshot.levels[index]
		level.ID = index + 1
		level.Track = TrackCore
		level.TrackPosition = index + 1
		snapshot.byKey[level.Key] = index
	}
	snapshot.err = validateCatalogue(snapshot)
	return snapshot
}

func validateCatalogue(snapshot catalogueSnapshot) error {
	if len(snapshot.levels) != exerciseCount {
		return fmt.Errorf("expected %d exercises, got %d", exerciseCount, len(snapshot.levels))
	}
	seenKeys := make(map[ExerciseKey]struct{}, exerciseCount)
	seenTitles := make(map[string]ExerciseKey, exerciseCount)
	for index, level := range snapshot.levels {
		if level.definitionErr != nil {
			return fmt.Errorf("exercise %q authoring: %w", level.Key, level.definitionErr)
		}
		if level.ID != index+1 || level.Key == "" || level.Track != TrackCore ||
			level.TrackPosition != index+1 || level.Title == "" || level.Subject == "" ||
			level.StarterCode == "" || len(level.Tests) == 0 || len(level.Tests) > 18 {
			return fmt.Errorf("invalid exercise %d definition", level.ID)
		}
		if _, exists := seenKeys[level.Key]; exists {
			return fmt.Errorf("duplicate exercise key %q", level.Key)
		}
		seenKeys[level.Key] = struct{}{}
		if previous, exists := seenTitles[level.Title]; exists {
			return fmt.Errorf("duplicate exercise title %q for %q and %q", level.Title, previous, level.Key)
		}
		seenTitles[level.Title] = level.Key
		seenTests := make(map[string]struct{}, len(level.Tests))
		for _, current := range level.Tests {
			if strings.TrimSpace(current.ID) == "" {
				return fmt.Errorf("exercise %q has an empty test id", level.Key)
			}
			if _, exists := seenTests[current.ID]; exists {
				return fmt.Errorf("exercise %q has duplicate test id %q", level.Key, current.ID)
			}
			seenTests[current.ID] = struct{}{}
		}
	}
	if len(snapshot.byKey) != len(snapshot.levels) {
		return fmt.Errorf("catalogue key index is incomplete")
	}
	return nil
}

func catalogueLevels() []Level { return cloneLevels(currentCatalogue.levels) }

func catalogueFindPosition(position int) (Level, bool) {
	if position < 1 || position > len(currentCatalogue.levels) {
		return Level{}, false
	}
	return cloneLevel(currentCatalogue.levels[position-1]), true
}

func catalogueFindKey(key ExerciseKey) (Level, bool) {
	index, found := currentCatalogue.byKey[key]
	if !found {
		return Level{}, false
	}
	return cloneLevel(currentCatalogue.levels[index]), true
}

func catalogueLegacyKey(int) (ExerciseKey, bool) { return "", false }

func catalogueLegacyKeys() []ExerciseKey { return []ExerciseKey{} }

func cloneLevels(levels []Level) []Level {
	cloned := make([]Level, len(levels))
	for index, level := range levels {
		cloned[index] = cloneLevel(level)
	}
	return cloned
}

func cloneLevel(level Level) Level {
	level.Tests = append([]VisibleTest(nil), level.Tests...)
	for index := range level.Tests {
		level.Tests[index].labels = append([]string(nil), level.Tests[index].labels...)
	}
	level.Resources = append([]ExerciseResource{}, level.Resources...)
	level.Instructions.Constraints = append([]string(nil), level.Instructions.Constraints...)
	level.Instructions.Examples = append([]Example(nil), level.Instructions.Examples...)
	level.Instructions.Documentation = append([]DocumentationLink(nil), level.Instructions.Documentation...)
	level.Instructions.AllowedBuiltins = append([]string{}, level.Instructions.AllowedBuiltins...)
	level.Instructions.AllowedPackages = append([]string{}, level.Instructions.AllowedPackages...)
	level.Instructions.Allowed = append([]string(nil), level.Instructions.Allowed...)
	level.Instructions.Disallowed = append([]string(nil), level.Instructions.Disallowed...)
	level.Instructions.CommonPitfalls = append([]string(nil), level.Instructions.CommonPitfalls...)
	level.Instructions.Hints = append([]string(nil), level.Instructions.Hints...)
	return level
}
