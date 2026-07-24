package assessment

import (
	"fmt"
	"sort"
	"strings"
)

const exerciseCount = 171

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
	levels := foundationalLevels()
	imported := append(piscineLevels(), zone01Levels()...)
	sort.SliceStable(imported, func(left, right int) bool {
		return importedDifficultyRank(imported[left].Difficulty) < importedDifficultyRank(imported[right].Difficulty)
	})
	levels = append(levels, imported...)
	for index := range levels {
		levels[index].ID = index + 1
		levels[index].StarterCode = "package main\n\n" + levels[index].StarterCode
	}
	return levels
}

func importedDifficultyRank(difficulty string) int {
	ranks := map[string]int{
		"Piscine Quest 2":  10,
		"Checkpoint 5%":    15,
		"Piscine Quest 3":  20,
		"Checkpoint 10%":   25,
		"Piscine Quest 4":  30,
		"Piscine Quest 5":  35,
		"Checkpoint 20%":   40,
		"Piscine Quest 6":  45,
		"Piscine Quest 7":  50,
		"Checkpoint 35%":   55,
		"Piscine Quest 8":  60,
		"Piscine Quest 9":  65,
		"Checkpoint 50%":   70,
		"Piscine Advanced": 75,
		"Checkpoint 65%":   80,
	}
	if rank, found := ranks[difficulty]; found {
		return rank
	}
	return 1_000
}

func Validate() error {
	levels := Levels()
	if len(levels) != exerciseCount {
		return fmt.Errorf("expected %d levels, got %d", exerciseCount, len(levels))
	}
	for index, level := range levels {
		if level.ID != index+1 || len(level.Tests) == 0 || len(level.Tests) > 18 {
			return fmt.Errorf("invalid level %d definition", level.ID)
		}
		seen := map[string]bool{}
		for _, current := range level.Tests {
			if strings.TrimSpace(current.ID) == "" || seen[current.ID] {
				return fmt.Errorf("invalid test id in level %d", level.ID)
			}
			seen[current.ID] = true
		}
	}
	return nil
}
