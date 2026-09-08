package assessment

import "strings"

type ExerciseKey string

type AssessmentTrack string

const (
	TrackCore     AssessmentTrack = "core"
	TrackAdvanced AssessmentTrack = "advanced"
)

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
	labels   []string
}

type ExerciseResource struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Level struct {
	Key           ExerciseKey        `json:"key"`
	ID            int                `json:"id"`
	Track         AssessmentTrack    `json:"track"`
	TrackPosition int                `json:"trackPosition"`
	Title         string             `json:"title"`
	Topic         string             `json:"topic"`
	Difficulty    string             `json:"difficulty"`
	Stretch       bool               `json:"stretch"`
	Signature     string             `json:"signature"`
	StarterCode   string             `json:"starterCode"`
	Subject       string             `json:"subject"`
	Resources     []ExerciseResource `json:"resources"`
	Instructions  Instructions       `json:"instructions"`
	Tests         []VisibleTest      `json:"tests"`
	unrestricted  bool
	definitionErr error
}

func PublicLevels() []Level {
	return catalogueLevels()
}

// MatchesOfficialLabel reports whether an output label emitted by the pinned
// Zone01 grader belongs to this public check. A trailing * is a prefix match.
func (test VisibleTest) MatchesOfficialLabel(label string) bool {
	for _, candidate := range test.labels {
		if candidate == label {
			return true
		}
		if strings.HasSuffix(candidate, "*") && strings.HasPrefix(label, strings.TrimSuffix(candidate, "*")) {
			return true
		}
	}
	return false
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
