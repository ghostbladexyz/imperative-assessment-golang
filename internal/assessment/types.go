package assessment

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
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Topic        string        `json:"topic"`
	Difficulty   string        `json:"difficulty"`
	Stretch      bool          `json:"stretch"`
	Signature    string        `json:"signature"`
	StarterCode  string        `json:"starterCode"`
	Instructions Instructions  `json:"instructions"`
	Tests        []VisibleTest `json:"tests"`
	build        func([]VisibleTest) string
}

func (level Level) BuildHarness(tests []VisibleTest) string {
	return level.build(tests)
}

func PublicLevels() []Level {
	levels := Levels()
	for index := range levels {
		levels[index].build = nil
	}
	return levels
}

func SelectTests(level Level, ids []string) ([]VisibleTest, bool) {
	if len(ids) == 0 {
		return level.Tests, true
	}
	wanted := make(map[string]bool, len(ids))
	for _, id := range ids {
		wanted[id] = true
	}
	selected := make([]VisibleTest, 0, len(ids))
	for _, test := range level.Tests {
		if wanted[test.ID] {
			selected = append(selected, test)
			delete(wanted, test.ID)
		}
	}
	return selected, len(wanted) == 0
}

func FindLevel(id int) (Level, bool) {
	for _, level := range Levels() {
		if level.ID == id {
			return level, true
		}
	}
	return Level{}, false
}
