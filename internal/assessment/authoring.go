package assessment

import (
	"fmt"
	"strings"
)

type practiceCase struct {
	name     string
	input    string
	expected string
	payload  any
}

type practiceSpec struct {
	id          int
	title       string
	topic       string
	difficulty  string
	signature   string
	objective   string
	input       string
	output      string
	constraints []string
	hints       []string
	pitfalls    []string
	builtins    []string
	packages    []string
	starter     string
	inputFields string
	call        string
	cases       []practiceCase
	answerMode  answerMode
}

type answerMode uint8

const (
	returnedAnswer answerMode = iota
	printedAnswer
)

const maxPracticeExamples = 5

const z01Package = "github.com/01-edu/z01"

func compilePracticeExercise(source exerciseSource, spec practiceSpec) Level {
	spec.cases = deduplicatePracticeCases(spec.cases)
	instructions := baseInstructions(
		spec.objective,
		"func "+spec.signature,
		spec.input,
		spec.output,
		"Complete the function body. Keep the signature unchanged.",
	)
	instructions.Constraints = append([]string(nil), spec.constraints...)
	instructions.Hints = append([]string(nil), spec.hints...)
	instructions.CommonPitfalls = append([]string(nil), spec.pitfalls...)
	instructions.Examples = practiceExamples(spec)
	packages := spec.packages
	if len(packages) == 0 {
		packages = []string{"fmt"}
	}
	setSourcePolicy(&instructions, spec.builtins, packages)

	tests := make([]VisibleTest, 0, len(spec.cases))
	for index, current := range spec.cases {
		tests = append(tests, test(
			fmt.Sprintf("l%d-%02d", spec.id, index+1),
			current.name,
			"Checks "+current.name+".",
			current.input,
			current.expected,
			current.payload,
		))
	}
	inputFields, call := spec.inputFields, spec.call
	return Level{
		Key: exerciseKey(source, spec.id), ID: spec.id, Title: spec.title, Topic: spec.topic,
		Difficulty: spec.difficulty, Signature: spec.signature,
		StarterCode: spec.starter, Instructions: instructions, Tests: tests,
		source: source, sourceID: spec.id, order: curriculumOrder(source, spec.id),
		definitionErr: validatePracticeSpec(source, spec),
		build: func(selected []VisibleTest) string {
			return buildPracticeHarness(inputFields, call, selected)
		},
	}
}

func printedPracticeSpec(spec practiceSpec) practiceSpec {
	closingParenthesis := strings.LastIndex(spec.signature, ")")
	if closingParenthesis >= 0 {
		spec.signature = spec.signature[:closingParenthesis+1]
	}
	spec.answerMode = printedAnswer
	spec.packages = []string{z01Package}
	spec.starter = fmt.Sprintf(`import %q

func %s {
	// TODO: print the required output with z01.PrintRune.
	_ = z01.PrintRune
}
`, z01Package, spec.signature)
	spec.constraints = []string{
		"Keep the required function name and parameters.",
		"Print the exact output with z01.PrintRune; do not return it.",
		"Match the stated edge cases and whitespace exactly.",
	}
	spec.hints = []string{
		"Start with the smallest example, then handle the boundary cases.",
		"Print one rune at a time with z01.PrintRune.",
		"Check separators and final newlines carefully.",
	}
	spec.pitfalls = []string{
		"Returning a value instead of printing it",
		"Using fmt instead of the required z01 package",
		"Printing an extra separator or missing the final newline",
	}
	return spec
}

// deduplicatePracticeCases keeps the first authored assertion because identical input and output exercise the same behavior.
func deduplicatePracticeCases(cases []practiceCase) []practiceCase {
	unique := make([]practiceCase, 0, len(cases))
	seen := make(map[string]struct{}, len(cases))
	for _, current := range cases {
		key := current.input + "\x00" + current.expected
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, current)
	}
	return unique
}

// practiceExamples shows up to five genuine cases from the current spec without padding the panel with duplicates.
func practiceExamples(spec practiceSpec) []Example {
	examples := make([]Example, 0, maxPracticeExamples)
	for index := 0; index < maxPracticeExamples && index < len(spec.cases); index++ {
		current := spec.cases[index]
		input := functionName(spec.signature) + "(" + current.input + ")"
		examples = append(examples, Example{Input: input, Output: current.expected})
	}
	return examples
}

func validatePracticeSpec(source exerciseSource, spec practiceSpec) error {
	if source == "" || spec.id < 1 {
		return fmt.Errorf("source and source id are required")
	}
	for label, value := range map[string]string{
		"title": spec.title, "topic": spec.topic, "difficulty": spec.difficulty,
		"signature": spec.signature, "objective": spec.objective,
		"input": spec.input, "output": spec.output, "starter": spec.starter,
		"call": spec.call,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required", label)
		}
	}
	if functionName(spec.signature) == "" ||
		!strings.Contains(spec.starter, "func "+spec.signature) {
		return fmt.Errorf("starter does not implement signature %q", spec.signature)
	}
	if len(spec.cases) == 0 || len(spec.cases) > 18 {
		return fmt.Errorf("between 1 and 18 cases are required")
	}
	for index, current := range spec.cases {
		if strings.TrimSpace(current.name) == "" || current.payload == nil {
			return fmt.Errorf("case %d requires a name and payload", index+1)
		}
	}
	return nil
}

func buildPracticeHarness(inputFields, call string, tests []VisibleTest) string {
	declarations := `type inputCase struct {
	` + inputFields + `
}
type wireCase struct {
	ID string ` + "`json:\"id\"`" + `
	Payload inputCase ` + "`json:\"payload\"`" + `
}`
	loop := `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value := ` + call + `
			encoded, _ := json.Marshal(value)
			return string(encoded)
		})
	}`
	return harness(commonImports(), declarations, loop, tests)
}

func buildPrintedPracticeHarness(inputFields, call string, tests []VisibleTest) string {
	declarations := `type inputCase struct {
	` + inputFields + `
}
type wireCase struct {
	ID string ` + "`json:\"id\"`" + `
	Payload inputCase ` + "`json:\"payload\"`" + `
}`
	loop := `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRunPrinted(current.ID, func() {
			` + call + `
		})
	}`
	return harness(commonImports(), declarations, loop, tests)
}

func pc(name, input, expected string, payload any) practiceCase {
	return practiceCase{name: name, input: input, expected: expected, payload: payload}
}

func functionName(signature string) string {
	for index, char := range signature {
		if char == '(' {
			return signature[:index]
		}
	}
	return ""
}

const (
	stringField       = "Value string `json:\"value\"`"
	intField          = "Value int `json:\"value\"`"
	twoIntsFields     = "Left int `json:\"left\"`\n\tRight int `json:\"right\"`"
	clampFields       = "Value int `json:\"value\"`\n\tMinimum int `json:\"minimum\"`\n\tMaximum int `json:\"maximum\"`"
	stringCountFields = "Value string `json:\"value\"`\n\tCount int `json:\"count\"`"
	intsField         = "Values []int `json:\"values\"`"
	intsTargetFields  = "Values []int `json:\"values\"`\n\tTarget int `json:\"target\"`"
	stringsField      = "Values []string `json:\"values\"`"
	intsStepsFields   = "Values []int `json:\"values\"`\n\tSteps int `json:\"steps\"`"
)
