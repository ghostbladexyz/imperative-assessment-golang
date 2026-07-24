package assessment

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const exerciseCount = 171

type exerciseSource string

const (
	sourceFoundation exerciseSource = "foundation"
	sourcePiscine    exerciseSource = "piscine"
	sourceZone01     exerciseSource = "zone01"
)

type catalogueSnapshot struct {
	levels     []Level
	byKey      map[ExerciseKey]int
	legacyKeys []ExerciseKey
	err        error
}

var currentCatalogue = buildCatalogue()

func baseInstructions(objective, contract, input, output, starter string) Instructions {
	return Instructions{
		Objective:       objective,
		Contract:        contract,
		Input:           input,
		Output:          output,
		StarterNote:     starter,
		Allowed:         []string{"User-defined helper functions and types", "Ordinary operators, loops, conditionals, goroutines, channels, and composite literals required by the task"},
		Disallowed:      []string{"Any package or Go built-in not listed for this exercise", "Third-party packages", "Network access", "Changing the required API signature"},
		WhitespaceRules: "Expected values use exact JSON or exact text comparison. Line endings are normalized; other whitespace is significant unless the task says otherwise.",
	}
}

func setSourcePolicy(instructions *Instructions, builtins, packages []string) {
	instructions.AllowedBuiltins = append([]string{"print", "println"}, builtins...)
	instructions.AllowedPackages = append([]string(nil), packages...)
}

func test(id, name, purpose, input, expected string, payload any) VisibleTest {
	return VisibleTest{ID: id, Name: name, Purpose: purpose, Input: input, Expected: expected, payload: payload}
}

func Levels() []Level {
	return catalogueLevels()
}

func Validate() error {
	return currentCatalogue.err
}

func buildCatalogue() catalogueSnapshot {
	levels := foundationalLevels()
	imported := append(piscineLevels(), zone01Levels()...)
	sort.SliceStable(imported, func(left, right int) bool {
		return imported[left].order < imported[right].order
	})
	levels = append(levels, imported...)

	snapshot := catalogueSnapshot{
		levels:     levels,
		byKey:      make(map[ExerciseKey]int, len(levels)),
		legacyKeys: legacyV4ExerciseKeys(),
	}
	for index := range snapshot.levels {
		level := &snapshot.levels[index]
		level.ID = index + 1
		level.StarterCode = "package main\n\n" + level.StarterCode
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
		if level.ID != index+1 || level.Key == "" || level.sourceID < 1 ||
			len(level.Tests) == 0 || len(level.Tests) > 18 || level.build == nil {
			return fmt.Errorf("invalid exercise %d definition", level.ID)
		}
		if level.Key != exerciseKey(level.source, level.sourceID) {
			return fmt.Errorf("exercise %d has invalid key %q", level.ID, level.Key)
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
	if len(snapshot.legacyKeys) != exerciseCount {
		return fmt.Errorf("schema-v4 migration has %d keys, want %d", len(snapshot.legacyKeys), exerciseCount)
	}
	seenLegacyKeys := make(map[ExerciseKey]struct{}, len(snapshot.legacyKeys))
	for position, key := range snapshot.legacyKeys {
		if _, exists := snapshot.byKey[key]; !exists {
			return fmt.Errorf("schema-v4 position %d references unknown key %q", position+1, key)
		}
		if _, exists := seenLegacyKeys[key]; exists {
			return fmt.Errorf("schema-v4 migration repeats key %q", key)
		}
		seenLegacyKeys[key] = struct{}{}
	}
	return nil
}

func legacyV4ExerciseKeys() []ExerciseKey {
	keys := make([]ExerciseKey, 0, exerciseCount)
	appendRange := func(source exerciseSource, first, last int) {
		for sourceID := first; sourceID <= last; sourceID++ {
			keys = append(keys, exerciseKey(source, sourceID))
		}
	}

	appendRange(sourceFoundation, 1, 21)
	appendRange(sourcePiscine, 1001, 1018)
	appendRange(sourceZone01, 22, 22)
	appendRange(sourcePiscine, 1019, 1028)
	appendRange(sourceZone01, 23, 29)
	appendRange(sourcePiscine, 1029, 1032)
	appendRange(sourcePiscine, 1033, 1056)
	appendRange(sourceZone01, 30, 38)
	appendRange(sourcePiscine, 1057, 1065)
	appendRange(sourcePiscine, 1066, 1075)
	appendRange(sourceZone01, 39, 49)
	appendRange(sourcePiscine, 1076, 1088)
	appendRange(sourcePiscine, 1089, 1090)
	appendRange(sourceZone01, 50, 61)
	appendRange(sourcePiscine, 1091, 1106)
	appendRange(sourceZone01, 62, 65)
	return keys
}

func exerciseKey(source exerciseSource, sourceID int) ExerciseKey {
	return ExerciseKey(string(source) + "/" + strconv.Itoa(sourceID))
}

func curriculumOrder(source exerciseSource, sourceID int) int {
	switch source {
	case sourceFoundation:
		return 0
	case sourcePiscine:
		switch {
		case sourceID <= 1018:
			return 10
		case sourceID <= 1028:
			return 20
		case sourceID <= 1032:
			return 30
		case sourceID <= 1056:
			return 35
		case sourceID <= 1065:
			return 45
		case sourceID <= 1075:
			return 50
		case sourceID <= 1088:
			return 60
		case sourceID <= 1090:
			return 65
		default:
			return 75
		}
	case sourceZone01:
		switch {
		case sourceID == 22:
			return 15
		case sourceID <= 29:
			return 25
		case sourceID <= 38:
			return 40
		case sourceID <= 49:
			return 55
		case sourceID <= 61:
			return 70
		default:
			return 80
		}
	default:
		return 1_000
	}
}

func catalogueLevels() []Level {
	return cloneLevels(currentCatalogue.levels)
}

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

func catalogueLegacyKey(position int) (ExerciseKey, bool) {
	if position < 1 || position > len(currentCatalogue.legacyKeys) {
		return "", false
	}
	return currentCatalogue.legacyKeys[position-1], true
}

func catalogueLegacyKeys() []ExerciseKey {
	return append([]ExerciseKey(nil), currentCatalogue.legacyKeys...)
}

func cloneLevels(levels []Level) []Level {
	cloned := make([]Level, len(levels))
	for index, level := range levels {
		cloned[index] = cloneLevel(level)
	}
	return cloned
}

func cloneLevel(level Level) Level {
	level.Tests = append([]VisibleTest(nil), level.Tests...)
	level.Instructions.Constraints = append([]string(nil), level.Instructions.Constraints...)
	level.Instructions.Examples = append([]Example(nil), level.Instructions.Examples...)
	level.Instructions.Documentation = append([]DocumentationLink(nil), level.Instructions.Documentation...)
	level.Instructions.AllowedBuiltins = append([]string(nil), level.Instructions.AllowedBuiltins...)
	level.Instructions.AllowedPackages = append([]string(nil), level.Instructions.AllowedPackages...)
	level.Instructions.Allowed = append([]string(nil), level.Instructions.Allowed...)
	level.Instructions.Disallowed = append([]string(nil), level.Instructions.Disallowed...)
	level.Instructions.CommonPitfalls = append([]string(nil), level.Instructions.CommonPitfalls...)
	level.Instructions.Hints = append([]string(nil), level.Instructions.Hints...)
	return level
}
