package assessment

type ExerciseKey string

type DocumentationLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Example struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Instructions struct {
	Objective       string              `json:"objective"`
	Contract        string              `json:"contract"`
	Input           string              `json:"input"`
	Output          string              `json:"output"`
	Constraints     []string            `json:"constraints"`
	Examples        []Example           `json:"examples"`
	Documentation   []DocumentationLink `json:"documentation"`
	AllowedBuiltins []string            `json:"allowedBuiltins"`
	AllowedPackages []string            `json:"allowedPackages"`
	Allowed         []string            `json:"allowed"`
	Disallowed      []string            `json:"disallowed"`
	StarterNote     string              `json:"starterNote"`
	WhitespaceRules string              `json:"whitespaceRules"`
	CommonPitfalls  []string            `json:"commonPitfalls"`
	Hints           []string            `json:"hints"`
}

type VisibleTest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Purpose  string `json:"purpose"`
	Input    string `json:"input"`
	Expected string `json:"expected"`
	payload  any
}

type Level struct {
	Key           ExerciseKey   `json:"key"`
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	Topic         string        `json:"topic"`
	Difficulty    string        `json:"difficulty"`
	Stretch       bool          `json:"stretch"`
	Signature     string        `json:"signature"`
	StarterCode   string        `json:"starterCode"`
	Instructions  Instructions  `json:"instructions"`
	Tests         []VisibleTest `json:"tests"`
	build         func([]VisibleTest) string
	source        exerciseSource
	sourceID      int
	order         int
	definitionErr error
}

func (level Level) BuildHarness(tests []VisibleTest) string {
	return level.build(tests)
}

func PublicLevels() []Level {
	levels := catalogueLevels()
	for index := range levels {
		levels[index].build = nil
	}
	return levels
}

func SelectTests(level Level, ids []string) ([]VisibleTest, bool) {
	if len(ids) == 0 {
		return level.Tests, true
	}
	available := make(map[string]VisibleTest, len(level.Tests))
	for _, test := range level.Tests {
		available[test.ID] = test
	}
	selected := make([]VisibleTest, 0, len(ids))
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		test, found := available[id]
		if !found || seen[id] {
			return nil, false
		}
		selected = append(selected, test)
		seen[id] = true
	}
	return selected, true
}

func FindLevel(id int) (Level, bool) {
	return catalogueFindPosition(id)
}

func FindExercise(key ExerciseKey) (Level, bool) {
	return catalogueFindKey(key)
}

func LegacyExerciseKey(position int) (ExerciseKey, bool) {
	return catalogueLegacyKey(position)
}

func LegacyExerciseKeys() []ExerciseKey {
	return catalogueLegacyKeys()
}
