package assessment

import (
	"fmt"
	"strings"
)

type advancedSpec struct {
	id          int
	title       string
	topic       string
	signature   string
	objective   string
	input       string
	output      string
	starter     string
	constraints []string
	hints       []string
	pitfalls    []string
	builtins    []string
	packages    []string
	tests       []VisibleTest
	build       func([]VisibleTest) string
}

// advancedLevels builds the original capstone catalogue that follows the imported checkpoint curriculum.
func advancedLevels() []Level {
	specs := []advancedSpec{
		parsePortsExercise(),
		eventLedgerExercise(),
		glyphCanvasExercise(),
		dependencyPlanExercise(),
		widgetHandlerExercise(),
		orderedWorkersExercise(),
		atomicTransferExercise(),
		jsonLinesExercise(),
		roomSchedulerExercise(),
		errorInspectorExercise(),
		compensatingStepsExercise(),
		profileClientExercise(),
		conditionalDocumentExercise(),
		weightedMazeExercise(),
		channelFanInExercise(),
		firstSuccessExercise(),
		overdueInvoicesExercise(),
		bulkOrderExercise(),
	}
	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		levels = append(levels, compileAdvancedExercise(spec))
	}
	return levels
}

// compileAdvancedExercise projects a specialized grader onto the same Level interface used by ordinary exercises.
func compileAdvancedExercise(spec advancedSpec) Level {
	instructions := baseInstructions(
		spec.objective,
		"func "+spec.signature,
		spec.input,
		spec.output,
		"Complete the TODOs in the provided declarations. Keep every public name and signature unchanged.",
	)
	instructions.Constraints = append([]string(nil), spec.constraints...)
	instructions.Hints = append([]string(nil), spec.hints...)
	instructions.CommonPitfalls = append([]string(nil), spec.pitfalls...)
	for index := 0; index < len(spec.tests) && index < maxPracticeExamples; index++ {
		current := spec.tests[index]
		instructions.Examples = append(instructions.Examples, Example{
			Input:  functionName(spec.signature) + "(" + current.Input + ")",
			Output: current.Expected,
		})
	}
	setSourcePolicy(&instructions, spec.builtins, spec.packages)
	level := Level{
		Key: exerciseKey(sourceAdvanced, spec.id), Title: spec.title,
		Track: TrackAdvanced,
		Topic: spec.topic, Difficulty: "Advanced", Stretch: true,
		Signature: spec.signature, StarterCode: spec.starter,
		Instructions: instructions, Tests: spec.tests,
		source: sourceAdvanced, sourceID: spec.id,
		order: curriculumOrder(sourceAdvanced, spec.id), build: spec.build,
	}
	level.definitionErr = validateAdvancedSpec(spec)
	return level
}

// validateAdvancedSpec rejects incomplete definitions before they can break the public catalogue at runtime.
func validateAdvancedSpec(spec advancedSpec) error {
	if spec.id < 1 || spec.title == "" || spec.topic == "" || spec.signature == "" ||
		spec.objective == "" || spec.input == "" || spec.output == "" || spec.starter == "" ||
		spec.build == nil || len(spec.tests) == 0 || len(spec.tests) > 18 {
		return fmt.Errorf("advanced exercise %d is incomplete", spec.id)
	}
	if functionName(spec.signature) == "" {
		return fmt.Errorf("advanced exercise %d has an invalid signature", spec.id)
	}
	if !strings.Contains(spec.starter, "func "+spec.signature) {
		return fmt.Errorf("advanced exercise %d starter does not implement %q", spec.id, spec.signature)
	}
	return nil
}

// advancedTests converts compact authored rows into stable, source-qualified test identities and payload selectors.
func advancedTests(id int, rows ...[4]string) []VisibleTest {
	tests := make([]VisibleTest, 0, len(rows))
	for index, row := range rows {
		tests = append(tests, test(
			fmt.Sprintf("a%d-%02d", id, index+1), row[0], row[1], row[2], row[3],
			map[string]any{"case": index},
		))
	}
	return tests
}

const advancedCaseDeclarations = `type inputCase struct {
	Case int ` + "`json:\"case\"`" + `
}
type wireCase struct {
	ID string ` + "`json:\"id\"`" + `
	Payload inputCase ` + "`json:\"payload\"`" + `
}`

// parsePortsExercise defines the strict parsing and normalization capstone inspired by the photographed CLI task.
func parsePortsExercise() advancedSpec {
	tests := advancedTests(1,
		[4]string{"Singles", "Parses comma-separated values.", `"443, 80, 8080"`, `[80 443 8080]|<nil>`},
		[4]string{"Ranges and duplicates", "Expands ranges, removes duplicates, and sorts.", `"8002-8004,8001,8003"`, `[8001 8002 8003 8004]|<nil>`},
		[4]string{"Boundary ports", "Accepts the complete valid port boundary.", `"65535,1"`, `[1 65535]|<nil>`},
		[4]string{"Descending range", "Rejects a descending range.", `"90-80"`, `[]|invalid port specification`},
		[4]string{"Malformed token", "Rejects empty and non-numeric tokens.", `"80,,http"`, `[]|invalid port specification`},
		[4]string{"Out of range", "Rejects values outside 1..65535.", `"0,65536"`, `[]|invalid port specification`},
	)
	return advancedSpec{
		id: 1, title: "Port Set", topic: "Parsing & data handling · validation · normalization",
		signature: "ParsePorts(spec string) ([]int, error)",
		objective: "Parse a compact port specification into a sorted set of distinct ports.",
		input:     "A comma-separated string of decimal ports and inclusive ranges such as 8000-8003.",
		output:    "A non-nil ascending []int and nil error, or an empty non-nil slice and ErrInvalidPortSpec.",
		constraints: []string{
			"Trim whitespace around comma-separated tokens and range endpoints.",
			"Every port must be in 1..65535; ranges must be ascending and contain at most 1024 ports.",
			"Reject empty tokens, extra hyphens, signs, and non-base-10 text without returning partial data.",
		},
		hints: []string{
			"Split into tokens first, then decide whether each token is a single value or one range.",
			"Use a map as a set and sort only after all validation succeeds.",
			"Return ErrInvalidPortSpec for every invalid form so callers can use errors.Is.",
		},
		pitfalls: []string{"Accepting Atoi signs such as +80", "Expanding an unbounded range", "Returning nil instead of []int{} on failure"},
		builtins: []string{"append", "len", "make"},
		packages: []string{"errors", "sort", "strconv", "strings"},
		starter: `import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

var ErrInvalidPortSpec = errors.New("invalid port specification")

func ParsePorts(spec string) ([]int, error) {
	// TODO: parse, validate, deduplicate, and sort the port set.
	_ = sort.Ints
	_ = strconv.Atoi
	_ = strings.TrimSpace
	return []int{}, ErrInvalidPortSpec
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	inputs := []string{"443, 80, 8080", "8002-8004,8001,8003", "65535,1", "90-80", "80,,http", "0,65536"}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			ports, err := ParsePorts(inputs[current.Payload.Case])
			return fmt.Sprintf("%v|%v", ports, err)
		})
	}`
			return harness(commonImports(), advancedCaseDeclarations, loop, selected)
		},
	}
}

// eventLedgerExercise defines CSV, time, aggregation, and contextual-error behavior as one data pipeline.
func eventLedgerExercise() advancedSpec {
	tests := advancedTests(2,
		[4]string{"Aggregate", "Aggregates users and keeps their latest timestamp.", `"2026-01-01T10:00:00Z,ana,5\n2026-01-01T11:00:00Z,bob,2\n2026-01-01T12:00:00Z,ana,-1"`, `[{"User":"ana","Total":4,"Last":"2026-01-01T12:00:00Z"},{"User":"bob","Total":2,"Last":"2026-01-01T11:00:00Z"}]|<nil>`},
		[4]string{"Comments and blanks", "Skips blank lines and full-line comments.", `"# export\n\n2026-02-03T00:00:00+02:00,zoe,9"`, `[{"User":"zoe","Total":9,"Last":"2026-02-02T22:00:00Z"}]|<nil>`},
		[4]string{"Quoted user", "Uses CSV quoting rather than splitting blindly.", `"2026-01-01T00:00:00Z,\"lee, jr\",7"`, `[{"User":"lee, jr","Total":7,"Last":"2026-01-01T00:00:00Z"}]|<nil>`},
		[4]string{"Bad delta", "Reports the one-based source line and no partial summary.", `"2026-01-01T00:00:00Z,ana,1\n2026-01-02T00:00:00Z,ana,nope"`, `[]|line 2: invalid event`},
		[4]string{"Bad field count", "Rejects records that do not contain exactly three fields.", `"2026-01-01T00:00:00Z,ana"`, `[]|line 1: invalid event`},
	)
	return advancedSpec{
		id: 2, title: "Event Ledger", topic: "Parsing & data handling · CSV · time · aggregation",
		signature: "SummarizeEvents(input string) ([]UserSummary, error)",
		objective: "Parse an exported event ledger and produce deterministic per-user totals.",
		input:     "UTF-8 CSV records timestamp,user,delta; blank lines and lines beginning with # are ignored.",
		output:    "User summaries sorted by User, with UTC RFC3339 timestamps, or an empty slice and a line-numbered error.",
		constraints: []string{
			"Parse each non-comment line as one CSV record with exactly three fields.",
			"User must be non-empty after trimming; timestamp must be RFC3339; delta must fit in int64.",
			"Wrap ErrInvalidEvent as `line N: invalid event` and discard all partial summaries on failure.",
		},
		hints: []string{
			"Track the original line number before filtering blank and comment lines.",
			"Aggregate in a map, then copy values to a slice and sort by User.",
			"Compare timestamps with Time.After and format the chosen value in UTC.",
		},
		pitfalls: []string{"Using strings.Split for quoted CSV", "Sorting input lines instead of final users", "Returning partial aggregates after an invalid row"},
		builtins: []string{"append", "len", "make"},
		packages: []string{"encoding/csv", "errors", "fmt", "io", "sort", "strconv", "strings", "time"},
		starter: `import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidEvent = errors.New("invalid event")

type UserSummary struct {
	User  string
	Total int64
	Last  string
}

func SummarizeEvents(input string) ([]UserSummary, error) {
	// TODO: parse valid records, aggregate them, and return stable output.
	_ = csv.NewReader
	_ = fmt.Errorf
	_ = io.EOF
	_ = sort.Slice
	_ = strconv.ParseInt
	_ = strings.TrimSpace
	_ = time.RFC3339
	return []UserSummary{}, ErrInvalidEvent
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	inputs := []string{
		"2026-01-01T10:00:00Z,ana,5\n2026-01-01T11:00:00Z,bob,2\n2026-01-01T12:00:00Z,ana,-1",
		"# export\n\n2026-02-03T00:00:00+02:00,zoe,9",
		"2026-01-01T00:00:00Z,\"lee, jr\",7",
		"2026-01-01T00:00:00Z,ana,1\n2026-01-02T00:00:00Z,ana,nope",
		"2026-01-01T00:00:00Z,ana",
	}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, err := SummarizeEvents(inputs[current.Payload.Case])
			encoded, _ := json.Marshal(value)
			return fmt.Sprintf("%s|%v", encoded, err)
		})
	}`
			return harness(commonImports(), advancedCaseDeclarations, loop, selected)
		},
	}
}

// glyphCanvasExercise defines exact row-wise rendering without copying the photographed fixed ASCII-banner contract.
func glyphCanvasExercise() advancedSpec {
	tests := advancedTests(3,
		[4]string{"One line", "Composes variable-width glyphs row by row.", `"AB", glyphs`, `"/\\[]\n\\/{}\n"|<nil>`},
		[4]string{"Text blocks", "Preserves empty input lines as empty output lines.", `"A\n\nB", glyphs`, `"/\\\n\\/\n\n[]\n{}\n"|<nil>`},
		[4]string{"Unicode key", "Indexes glyphs by rune rather than byte.", `"λA", glyphs`, `"\u003c\u003e/\\\n\u003c \u003e\\/\n"|<nil>`},
		[4]string{"Missing glyph", "Returns no partial rendering when a glyph is absent.", `"AC", glyphs`, `""|missing glyph: U+0043`},
		[4]string{"Height mismatch", "Rejects a glyph set with inconsistent heights.", `"AB", malformed`, `""|inconsistent glyph height`},
	)
	return advancedSpec{
		id: 3, title: "Glyph Canvas", topic: "Parsing & data handling · Unicode · exact rendering",
		signature: "RenderGlyphs(input string, glyphs map[rune][]string) (string, error)",
		objective: "Render multiline text from a supplied, variable-width Unicode glyph map.",
		input:     "Text split by newline and a non-empty map whose glyphs contain equal, non-zero row counts.",
		output:    "Exact row-wise rendering with one final newline per rendered row, or a sentinel-compatible error and no partial output.",
		constraints: []string{
			"Validate the complete glyph map before rendering input.",
			"Iterate input as runes; concatenate glyph rows without separators.",
			"Each empty logical input line contributes exactly one empty output line.",
		},
		hints: []string{
			"Determine and validate the common glyph height in a first pass.",
			"Render one input line at a time, with an inner loop for glyph rows.",
			"Use fmt.Errorf with %w for ErrMissingGlyph and include the U+ code point.",
		},
		pitfalls: []string{"Indexing UTF-8 bytes instead of runes", "Adding spaces between glyphs", "Discovering malformed glyphs after writing partial output"},
		builtins: []string{"len"},
		packages: []string{"errors", "fmt", "strings"},
		starter: `import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingGlyph = errors.New("missing glyph")
	ErrGlyphHeight  = errors.New("inconsistent glyph height")
)

func RenderGlyphs(input string, glyphs map[rune][]string) (string, error) {
	// TODO: validate the map and render every logical input line.
	_ = fmt.Errorf
	_ = strings.Builder{}
	return "", ErrMissingGlyph
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	glyphs := map[rune][]string{'A': {"/\\", "\\/"}, 'B': {"[]", "{}"}, 'λ': {"<>", "< >"}}
	malformed := map[rune][]string{'A': {"/\\", "\\/"}, 'B': {"[]"}}
	inputs := []string{"AB", "A\n\nB", "λA", "AC", "AB"}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			set := glyphs
			if current.Payload.Case == 4 { set = malformed }
			value, err := RenderGlyphs(inputs[current.Payload.Case], set)
			encoded, _ := json.Marshal(value)
			return fmt.Sprintf("%s|%v", encoded, err)
		})
	}`
			return harness(commonImports(), advancedCaseDeclarations, loop, selected)
		},
	}
}

// dependencyPlanExercise defines deterministic graph validation and topological ordering.
func dependencyPlanExercise() advancedSpec {
	tests := advancedTests(4,
		[4]string{"Deterministic order", "Uses lexical order whenever several tasks are ready.", `tasks`, `[fetch compile lint test deploy]|<nil>`},
		[4]string{"Diamond", "Emits every task once across shared dependencies.", `diamond`, `[a b c d]|<nil>`},
		[4]string{"Missing dependency", "Rejects references to tasks that do not exist.", `missing`, `[]|invalid dependency graph: missing task x`},
		[4]string{"Duplicate task", "Rejects duplicate task names.", `duplicate`, `[]|invalid dependency graph: duplicate task a`},
		[4]string{"Cycle", "Detects cycles instead of returning a partial order.", `cycle`, `[]|dependency cycle`},
	)
	return advancedSpec{
		id: 4, title: "Dependency Plan", topic: "Algorithms · graphs · deterministic topological sorting",
		signature: "BuildOrder(tasks []Task) ([]string, error)",
		objective: "Produce a deterministic build order for a dependency graph.",
		input:     "Named tasks with zero or more named dependencies; task names and each dependency list may arrive unsorted.",
		output:    "Every task exactly once, choosing the lexicographically smallest ready task, or a classified validation/cycle error.",
		constraints: []string{
			"Reject empty names, duplicate task names, duplicate dependencies, self-dependencies, and missing tasks as ErrInvalidGraph.",
			"Return ErrDependencyCycle when valid declarations contain a cycle.",
			"Return an empty non-nil slice on every error and do not mutate the input.",
		},
		hints: []string{
			"Validate names while building an indegree map and reverse adjacency list.",
			"Kahn's algorithm exposes the ready set after each completed task.",
			"Keep the ready set sorted, or use a small min-heap.",
		},
		pitfalls: []string{"Depending on map iteration order", "Returning a partial order for a cycle", "Counting a duplicate dependency twice"},
		builtins: []string{"append", "len", "make"},
		packages: []string{"errors", "fmt", "sort"},
		starter: `import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrInvalidGraph    = errors.New("invalid dependency graph")
	ErrDependencyCycle = errors.New("dependency cycle")
)

type Task struct {
	Name      string
	DependsOn []string
}

func BuildOrder(tasks []Task) ([]string, error) {
	// TODO: validate the graph and perform a deterministic topological sort.
	_ = fmt.Errorf
	_ = sort.Strings
	return []string{}, ErrInvalidGraph
}
`,
		tests: tests,
		build: func(selected []VisibleTest) string {
			loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	cases := [][]Task{
		{{Name:"deploy",DependsOn:[]string{"test"}},{Name:"test",DependsOn:[]string{"compile"}},{Name:"lint"},{Name:"compile",DependsOn:[]string{"fetch"}},{Name:"fetch"}},
		{{Name:"d",DependsOn:[]string{"b","c"}},{Name:"c",DependsOn:[]string{"a"}},{Name:"b",DependsOn:[]string{"a"}},{Name:"a"}},
		{{Name:"a",DependsOn:[]string{"x"}}},
		{{Name:"a"},{Name:"a"}},
		{{Name:"a",DependsOn:[]string{"b"}},{Name:"b",DependsOn:[]string{"a"}}},
	}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, err := BuildOrder(cases[current.Payload.Case])
			return fmt.Sprintf("%v|%v", value, err)
		})
	}`
			return harness(commonImports(), advancedCaseDeclarations, loop, selected)
		},
	}
}

// widgetHandlerExercise defines the HTTP contract while the harness supplies an isolated storage adapter.
func widgetHandlerExercise() advancedSpec {
	tests := advancedTests(5,
		[4]string{"GET widget", "Returns a JSON representation for an existing widget.", `GET /widgets/a`, `200|application/json|{"id":"a","name":"Alpha"}`},
		[4]string{"PUT widget", "Validates JSON and stores a replacement.", `PUT /widgets/a {"name":"Beta"}`, `204||`},
		[4]string{"Missing widget", "Maps ErrWidgetNotFound to 404.", `GET /widgets/missing`, `404||`},
		[4]string{"Strict JSON", "Rejects unknown fields and trailing JSON.", `PUT /widgets/a {"name":"Beta","extra":true}`, `400||`},
		[4]string{"Conflict", "Maps ErrWidgetConflict to 409.", `PUT /widgets/locked {"name":"Beta"}`, `409||`},
		[4]string{"Panic recovery", "Returns 500 and remains usable after a store panic.", `GET /widgets/panic; GET /widgets/a`, `500,200`},
		[4]string{"Method and path", "Rejects unsupported methods and unrelated paths.", `POST /widgets/a; GET /other`, `405,404`},
	)
	return advancedSpec{
		id: 5, title: "Widget Gateway", topic: "HTTP · JSON · status mapping · recovery",
		signature: "NewWidgetHandler(store WidgetStore) http.Handler",
		objective: "Build a resilient JSON HTTP handler around a supplied storage interface.",
		input:     "GET or PUT requests for /widgets/{id}; PUT accepts exactly one JSON object with a non-empty name.",
		output:    "The specified HTTP status, headers, and JSON body while keeping the handler usable after failures.",
		constraints: []string{
			"GET success returns 200, Content-Type application/json, and the stored Widget as JSON.",
			"PUT success returns 204; malformed/unknown/trailing JSON or blank names return 400.",
			"Map ErrWidgetNotFound to 404, ErrWidgetConflict to 409, other errors and panics to 500, unsupported methods to 405, and unrelated paths to 404.",
		},
		hints: []string{
			"Return one http.HandlerFunc and validate the path before dispatching on Method.",
			"Use json.Decoder.DisallowUnknownFields and verify that a second Decode reaches io.EOF.",
			"Put a recover defer at the outermost request scope and avoid writing success headers early.",
		},
		pitfalls: []string{"Writing 200 before the store operation succeeds", "Accepting /widgets/ with an empty ID", "Recovering a panic but leaving the response as 200"},
		builtins: []string{"recover"},
		packages: []string{"encoding/json", "errors", "io", "net/http", "strings"},
		starter: `import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var (
	ErrWidgetNotFound = errors.New("widget not found")
	ErrWidgetConflict = errors.New("widget conflict")
)

type Widget struct {
	ID   string ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

type WidgetStore interface {
	Get(id string) (Widget, error)
	Put(widget Widget) error
}

func NewWidgetHandler(store WidgetStore) http.Handler {
	// TODO: implement routing, strict JSON, status mapping, and panic recovery.
	_ = json.NewDecoder
	_ = io.EOF
	_ = strings.TrimSpace
	return http.NotFoundHandler()
}
`,
		tests: tests,
		build: buildWidgetHarness,
	}
}

// buildWidgetHarness supplies requests and a fake WidgetStore so the exercise needs no network access.
func buildWidgetHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `

type assessmentWidgetStore struct {
	items map[string]Widget
}

func (store *assessmentWidgetStore) Get(id string) (Widget, error) {
	switch id {
	case "panic": panic("storage panic")
	case "missing": return Widget{}, ErrWidgetNotFound
	}
	value, found := store.items[id]
	if !found { return Widget{}, ErrWidgetNotFound }
	return value, nil
}

func (store *assessmentWidgetStore) Put(widget Widget) error {
	if widget.ID == "locked" { return ErrWidgetConflict }
	if widget.ID == "panic" { panic("storage panic") }
	store.items[widget.ID] = widget
	return nil
}

func assessmentHTTP(handler http.Handler, method, path, body string) (int, string, string) {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	contentType := recorder.Header().Get("Content-Type")
	if marker := strings.IndexByte(contentType, ';'); marker >= 0 { contentType = contentType[:marker] }
	return recorder.Code, contentType, strings.TrimSpace(recorder.Body.String())
}`
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			store := &assessmentWidgetStore{items: map[string]Widget{"a": {ID:"a", Name:"Alpha"}}}
			handler := NewWidgetHandler(store)
			switch current.Payload.Case {
			case 0:
				status, contentType, body := assessmentHTTP(handler, "GET", "/widgets/a", "")
				return fmt.Sprintf("%d|%s|%s", status, contentType, body)
			case 1:
				status, contentType, body := assessmentHTTP(handler, "PUT", "/widgets/a", ` + "`{\"name\":\"Beta\"}`" + `)
				if store.items["a"].Name != "Beta" { return "stored wrong widget" }
				return fmt.Sprintf("%d|%s|%s", status, contentType, body)
			case 2:
				status, _, _ := assessmentHTTP(handler, "GET", "/widgets/missing", "")
				return fmt.Sprintf("%d||", status)
			case 3:
				status, _, _ := assessmentHTTP(handler, "PUT", "/widgets/a", ` + "`{\"name\":\"Beta\",\"extra\":true}`" + `)
				return fmt.Sprintf("%d||", status)
			case 4:
				status, _, _ := assessmentHTTP(handler, "PUT", "/widgets/locked", ` + "`{\"name\":\"Beta\"}`" + `)
				return fmt.Sprintf("%d||", status)
			case 5:
				first, _, _ := assessmentHTTP(handler, "GET", "/widgets/panic", "")
				second, _, _ := assessmentHTTP(handler, "GET", "/widgets/a", "")
				return fmt.Sprintf("%d,%d", first, second)
			default:
				first, _, _ := assessmentHTTP(handler, "POST", "/widgets/a", "")
				second, _, _ := assessmentHTTP(handler, "GET", "/other", "")
				return fmt.Sprintf("%d,%d", first, second)
			}
		})
	}`
	return harness(commonImports("net/http", "net/http/httptest", "strings"), declarations, loop, selected)
}

// orderedWorkersExercise defines bounded concurrency, stable output, and cancellation as observable behavior.
func orderedWorkersExercise() advancedSpec {
	tests := advancedTests(6,
		[4]string{"Stable order", "Returns results in input order despite out-of-order completion.", `ctx, []int{5,1,3}, 3, transform`, `[10 2 6]|<nil>`},
		[4]string{"Actually concurrent", "Starts work concurrently when capacity permits.", `ctx, []int{1,2}, 2, barrier`, `[1 2]|<nil>`},
		[4]string{"Worker bound", "Never exceeds the requested worker count.", `ctx, six values, 2, transform`, `ordered|max<=2`},
		[4]string{"Transform failure", "Cancels remaining work and preserves the original error.", `ctx, []int{1,2,3,4}, 2, failing`, `[]|boom`},
		[4]string{"Invalid workers", "Rejects a non-positive worker count without invoking transform.", `ctx, []int{1}, 0, transform`, `[]|invalid worker count|calls=0`},
		[4]string{"Empty input", "Returns an empty result without starting workers.", `ctx, []int{}, 4, transform`, `[]|<nil>|calls=0`},
	)
	return advancedSpec{
		id: 6, title: "Ordered Workers", topic: "Concurrency · cancellation · bounded parallelism",
		signature: "MapConcurrent(ctx context.Context, values []int, workers int, transform func(context.Context, int) (int, error)) ([]int, error)",
		objective: "Apply a fallible operation concurrently while preserving input order and bounding active work.",
		input:     "A context, values, a positive worker limit, and a context-aware transform function.",
		output:    "Ordered results after complete success, or an empty non-nil slice and the triggering/context error.",
		constraints: []string{
			"At most workers calls to transform may be active at once, and successful multi-item calls must make concurrent progress.",
			"On the first transform error, cancel internal work, stop scheduling new values, wait for started calls, and return that error.",
			"Respect an already-cancelled parent context and never invoke transform for empty input or invalid workers.",
		},
		hints: []string{
			"Send indexed jobs to a fixed number of worker goroutines.",
			"Store each success at its original index and report only the first failure through a buffered channel or sync.Once.",
			"Derive a child context, cancel it on failure, close jobs once, and wait for all workers before returning.",
		},
		pitfalls: []string{"Appending results in completion order", "Launching one goroutine per value", "Returning while worker goroutines still use local channels"},
		builtins: []string{"append", "cap", "close", "len", "make"},
		packages: []string{"context", "errors", "sync"},
		starter: `import (
	"context"
	"errors"
	"sync"
)

var ErrInvalidWorkerCount = errors.New("invalid worker count")

func MapConcurrent(ctx context.Context, values []int, workers int, transform func(context.Context, int) (int, error)) ([]int, error) {
	// TODO: coordinate a bounded worker pool with cancellation and stable output.
	_ = sync.WaitGroup{}
	return []int{}, ErrInvalidWorkerCount
}
`,
		tests: tests,
		build: buildWorkersHarness,
	}
}

// buildWorkersHarness uses barriers and counters to detect sequential work and excess concurrency deterministically.
func buildWorkersHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `

var assessmentBoom = errors.New("boom")

type assessmentCounter struct {
	mu sync.Mutex
	active int
	maximum int
	calls int
}

func (counter *assessmentCounter) enter() {
	counter.mu.Lock()
	counter.calls++
	counter.active++
	if counter.active > counter.maximum { counter.maximum = counter.active }
	counter.mu.Unlock()
}

func (counter *assessmentCounter) leave() {
	counter.mu.Lock()
	counter.active--
	counter.mu.Unlock()
}`
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			counter := &assessmentCounter{}
			switch current.Payload.Case {
			case 0:
				value, err := MapConcurrent(context.Background(), []int{5,1,3}, 3, func(ctx context.Context, item int) (int,error) {
					time.Sleep(time.Duration(6-item) * time.Millisecond)
					return item*2, nil
				})
				return fmt.Sprintf("%v|%v", value, err)
			case 1:
				ready := make(chan struct{}, 2)
				release := make(chan struct{})
				go func(){ <-ready; <-ready; close(release) }()
				ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
				defer cancel()
				value, err := MapConcurrent(ctx, []int{1,2}, 2, func(ctx context.Context, item int) (int,error) {
					ready <- struct{}{}
					select { case <-release: return item,nil; case <-ctx.Done(): return 0,ctx.Err() }
				})
				return fmt.Sprintf("%v|%v", value, err)
			case 2:
				value, err := MapConcurrent(context.Background(), []int{1,2,3,4,5,6}, 2, func(ctx context.Context, item int) (int,error) {
					counter.enter(); defer counter.leave(); time.Sleep(3*time.Millisecond); return item,nil
				})
				if err != nil || fmt.Sprint(value) != "[1 2 3 4 5 6]" || counter.maximum > 2 { return "invalid" }
				return "ordered|max<=2"
			case 3:
				value, err := MapConcurrent(context.Background(), []int{1,2,3,4}, 2, func(ctx context.Context, item int) (int,error) {
					if item == 2 { return 0, assessmentBoom }
					select { case <-ctx.Done(): return 0,ctx.Err(); case <-time.After(20*time.Millisecond): return item,nil }
				})
				if !errors.Is(err, assessmentBoom) { return fmt.Sprintf("%v|wrong error: %v", value, err) }
				return fmt.Sprintf("%v|boom", value)
			case 4:
				value, err := MapConcurrent(context.Background(), []int{1}, 0, func(context.Context,int)(int,error){ counter.calls++; return 1,nil })
				if !errors.Is(err, ErrInvalidWorkerCount) { return fmt.Sprintf("%v|%v|calls=%d", value, err, counter.calls) }
				return fmt.Sprintf("%v|invalid worker count|calls=%d", value, counter.calls)
			default:
				value, err := MapConcurrent(context.Background(), []int{}, 4, func(context.Context,int)(int,error){ counter.calls++; return 1,nil })
				return fmt.Sprintf("%v|%v|calls=%d", value, err, counter.calls)
			}
		})
	}`
	return harness(commonImports("context", "errors", "sync"), declarations, loop, selected)
}

// atomicTransferExercise defines transactional SQL behavior through the standard database/sql seam.
func atomicTransferExercise() advancedSpec {
	tests := advancedTests(7,
		[4]string{"Commit", "Performs both updates and commits in one transaction.", `ctx, db, 1, 2, 30`, `ok|begin,query,debit,credit,commit`},
		[4]string{"Insufficient funds", "Rolls back and returns ErrInsufficientFunds before updates.", `ctx, db, 1, 2, 30`, `insufficient funds|begin,query,rollback`},
		[4]string{"Debit failure", "Rolls back and preserves a debit error.", `ctx, db, 1, 2, 30`, `database failure|begin,query,debit,rollback`},
		[4]string{"Credit failure", "Rolls back after a credit error.", `ctx, db, 1, 2, 30`, `database failure|begin,query,debit,credit,rollback`},
		[4]string{"Commit failure", "Returns a failed commit and relies on transaction semantics.", `ctx, db, 1, 2, 30`, `database failure|begin,query,debit,credit,commit`},
		[4]string{"Invalid transfer", "Rejects invalid arguments before opening a transaction.", `ctx, db, 1, 1, 0`, `invalid transfer|`},
	)
	return advancedSpec{
		id: 7, title: "Atomic Transfer", topic: "SQL · transactions · rollback · error identity",
		signature: "TransferCredits(ctx context.Context, db *sql.DB, fromID int64, toID int64, amount int64) error",
		objective: "Move credits between accounts atomically through database/sql.",
		input:     "A context, database handle, distinct positive account IDs, and a positive integer amount.",
		output:    "nil only after a committed debit and credit; otherwise a sentinel or wrapped database error with correct rollback behavior.",
		constraints: []string{
			"Validate arguments before touching db, then BeginTx with default options.",
			"Read the source balance inside the transaction, return ErrInsufficientFunds when it is below amount, then debit source and credit destination with parameterized statements.",
			"Rollback every path after BeginTx that does not commit; preserve query, exec, and commit errors with %w wrapping.",
		},
		hints: []string{
			"Defer tx.Rollback immediately after a successful BeginTx; it is harmless after Commit.",
			"QueryRowContext(...).Scan(&balance) keeps the read in the transaction.",
			"Check each operation before proceeding and commit only after both updates succeed.",
		},
		pitfalls: []string{"Reading the balance through db instead of tx", "Crediting after a failed debit", "Returning nil when Commit fails"},
		builtins: []string{"new"},
		packages: []string{"context", "database/sql", "errors", "fmt"},
		starter: `import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrInvalidTransfer   = errors.New("invalid transfer")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

func TransferCredits(ctx context.Context, db *sql.DB, fromID int64, toID int64, amount int64) error {
	// TODO: validate, transact, roll back failures, and commit both updates.
	_ = fmt.Errorf
	return ErrInvalidTransfer
}
`,
		tests: tests,
		build: buildSQLHarness,
	}
}

// buildSQLHarness supplies a scripted driver so transaction order and rollback behavior are testable without a database.
func buildSQLHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `

var assessmentDatabaseFailure = errors.New("database failure")
var assessmentCurrentSQL *assessmentSQLScenario

type assessmentSQLScenario struct {
	mode string
	balance int64
	events []string
}

type assessmentSQLDriver struct{}
func (assessmentSQLDriver) Open(string) (driver.Conn, error) {
	return &assessmentSQLConn{scenario: assessmentCurrentSQL}, nil
}

type assessmentSQLConn struct { scenario *assessmentSQLScenario }
func (*assessmentSQLConn) Prepare(string) (driver.Stmt, error) { return nil, errors.New("prepared statements are not supported") }
func (*assessmentSQLConn) Close() error { return nil }
func (connection *assessmentSQLConn) Begin() (driver.Tx, error) { return connection.BeginTx(context.Background(), driver.TxOptions{}) }
func (connection *assessmentSQLConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	connection.scenario.events = append(connection.scenario.events, "begin")
	return &assessmentSQLTx{scenario: connection.scenario}, nil
}
func (connection *assessmentSQLConn) QueryContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	normalized := strings.ToLower(strings.Join(strings.Fields(query), " "))
	connection.scenario.events = append(connection.scenario.events, "query")
	if !strings.Contains(normalized, "select") || !strings.Contains(normalized, "balance") || !strings.Contains(normalized, "account") {
		return nil, assessmentDatabaseFailure
	}
	return &assessmentSQLRows{balance: connection.scenario.balance}, nil
}
func (connection *assessmentSQLConn) ExecContext(_ context.Context, query string, arguments []driver.NamedValue) (driver.Result, error) {
	normalized := strings.ToLower(strings.Join(strings.Fields(query), " "))
	if !strings.Contains(normalized, "update") || !strings.Contains(normalized, "balance") || len(arguments) < 2 {
		return nil, assessmentDatabaseFailure
	}
	isDebit := strings.Contains(normalized, "-")
	if isDebit {
		connection.scenario.events = append(connection.scenario.events, "debit")
		if connection.scenario.mode == "debit" { return nil, assessmentDatabaseFailure }
	} else {
		connection.scenario.events = append(connection.scenario.events, "credit")
		if connection.scenario.mode == "credit" { return nil, assessmentDatabaseFailure }
	}
	return driver.RowsAffected(1), nil
}

type assessmentSQLTx struct { scenario *assessmentSQLScenario }
func (transaction *assessmentSQLTx) Commit() error {
	transaction.scenario.events = append(transaction.scenario.events, "commit")
	if transaction.scenario.mode == "commit" { return assessmentDatabaseFailure }
	return nil
}
func (transaction *assessmentSQLTx) Rollback() error {
	transaction.scenario.events = append(transaction.scenario.events, "rollback")
	return nil
}

type assessmentSQLRows struct { balance int64; sent bool }
func (*assessmentSQLRows) Columns() []string { return []string{"balance"} }
func (*assessmentSQLRows) Close() error { return nil }
func (rows *assessmentSQLRows) Next(values []driver.Value) error {
	if rows.sent { return io.EOF }
	rows.sent = true
	values[0] = rows.balance
	return nil
}

func init() { sql.Register("imperative-assessment-sql", assessmentSQLDriver{}) }`
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			scenario := &assessmentSQLScenario{balance: 100}
			switch current.Payload.Case {
			case 1: scenario.balance = 5
			case 2: scenario.mode = "debit"
			case 3: scenario.mode = "credit"
			case 4: scenario.mode = "commit"
			}
			assessmentCurrentSQL = scenario
			db, err := sql.Open("imperative-assessment-sql", "")
			if err != nil { return err.Error() }
			defer db.Close()
			fromID, toID, amount := int64(1), int64(2), int64(30)
			if current.Payload.Case == 5 { toID, amount = 1, 0 }
			err = TransferCredits(context.Background(), db, fromID, toID, amount)
			label := "ok"
			switch {
			case errors.Is(err, ErrInvalidTransfer): label = "invalid transfer"
			case errors.Is(err, ErrInsufficientFunds): label = "insufficient funds"
			case errors.Is(err, assessmentDatabaseFailure): label = "database failure"
			case err != nil: label = "wrong error: " + err.Error()
			}
			return label + "|" + strings.Join(scenario.events, ",")
		})
	}`
	return harness(commonImports("context", "database/sql", "database/sql/driver", "errors", "strings"), declarations, loop, selected)
}

// jsonLinesExercise defines bounded streaming decode behavior with partial results and preserved I/O errors.
func jsonLinesExercise() advancedSpec {
	tests := advancedTests(8,
		[4]string{"Valid stream", "Decodes valid lines in order.", `reader, 128`, `[{"id":"a","payload":"one"},{"id":"b","payload":"two"}]|<nil>`},
		[4]string{"Blank lines", "Ignores whitespace-only lines without changing source line numbers.", `reader, 128`, `[{"id":"a","payload":"one"}]|<nil>`},
		[4]string{"Malformed second line", "Returns prior messages and a contextual sentinel error.", `reader, 128`, `[{"id":"a","payload":"one"}]|line 2: invalid message`},
		[4]string{"Unknown field", "Uses strict JSON decoding.", `reader, 128`, `[]|line 1: invalid message`},
		[4]string{"Duplicate ID", "Rejects a repeated message identity.", `reader, 128`, `[{"id":"a","payload":"one"}]|line 2: duplicate message id`},
		[4]string{"Line limit", "Rejects an oversized logical record without Scanner's default limit leaking through.", `reader, 12`, `[]|line 1: message too large`},
		[4]string{"Read failure", "Preserves successfully decoded messages and wraps the reader error.", `failingReader, 128`, `[{"id":"a","payload":"one"}]|read messages: source failed`},
	)
	return advancedSpec{
		id: 8, title: "Message Stream", topic: "Error handling · streaming JSON · partial results",
		signature: "DecodeMessages(reader io.Reader, maxLineBytes int) ([]Message, error)",
		objective: "Decode a bounded JSON-lines stream with strict validation and useful partial failure results.",
		input:     "An io.Reader containing one JSON object per non-blank line and a positive maximum encoded line size.",
		output:    "Messages decoded before success/failure; validation errors identify the physical line and read errors retain their identity.",
		constraints: []string{
			"Each object has exactly non-empty string fields id and payload; IDs must be unique.",
			"Ignore whitespace-only lines but count them when reporting later physical line numbers.",
			"Return ErrMessageTooLarge for a line over maxLineBytes, wrap ErrInvalidMessage or ErrDuplicateMessageID with `line N`, and wrap underlying reads as `read messages: ...`.",
		},
		hints: []string{
			"A bufio.Reader with ReadString or ReadBytes lets you enforce an explicit per-line limit and distinguish EOF from other read errors.",
			"Decode each trimmed line with DisallowUnknownFields and require a second decode to return io.EOF.",
			"Append only after a message is fully valid and its ID is new.",
		},
		pitfalls: []string{"Losing the final line when it has no newline", "Returning an empty slice instead of valid prior messages", "Treating a blank line as a JSON error"},
		builtins: []string{"append", "len", "make"},
		packages: []string{"bufio", "bytes", "encoding/json", "errors", "fmt", "io", "strings"},
		starter: `import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrInvalidMessage     = errors.New("invalid message")
	ErrDuplicateMessageID = errors.New("duplicate message id")
	ErrMessageTooLarge    = errors.New("message too large")
)

type Message struct {
	ID      string ` + "`json:\"id\"`" + `
	Payload string ` + "`json:\"payload\"`" + `
}

func DecodeMessages(reader io.Reader, maxLineBytes int) ([]Message, error) {
	// TODO: decode bounded, strict JSON lines while preserving valid partial results.
	_ = bufio.NewReader
	_ = bytes.NewReader
	_ = json.NewDecoder
	_ = fmt.Errorf
	_ = strings.TrimSpace
	return []Message{}, ErrInvalidMessage
}
`,
		tests: tests,
		build: buildJSONLinesHarness,
	}
}

// buildJSONLinesHarness adds a failing reader case because byte strings alone cannot exercise underlying I/O errors.
func buildJSONLinesHarness(selected []VisibleTest) string {
	declarations := advancedCaseDeclarations + `

var assessmentSourceFailure = errors.New("source failed")

type assessmentFailingReader struct { data []byte; sent bool }
func (reader *assessmentFailingReader) Read(destination []byte) (int,error) {
	if !reader.sent {
		reader.sent = true
		return copy(destination, reader.data), nil
	}
	return 0, assessmentSourceFailure
}`
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	inputs := []string{
		"{\"id\":\"a\",\"payload\":\"one\"}\n{\"id\":\"b\",\"payload\":\"two\"}",
		"\n  \n{\"id\":\"a\",\"payload\":\"one\"}\n",
		"{\"id\":\"a\",\"payload\":\"one\"}\nnot-json",
		"{\"id\":\"a\",\"payload\":\"one\",\"extra\":true}",
		"{\"id\":\"a\",\"payload\":\"one\"}\n{\"id\":\"a\",\"payload\":\"two\"}",
		"{\"id\":\"abcdef\",\"payload\":\"one\"}",
	}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			var source io.Reader
			limit := 128
			if current.Payload.Case == 6 {
				source = &assessmentFailingReader{data: []byte("{\"id\":\"a\",\"payload\":\"one\"}\n")}
			} else {
				source = strings.NewReader(inputs[current.Payload.Case])
			}
			if current.Payload.Case == 5 { limit = 12 }
			value, err := DecodeMessages(source, limit)
			encoded, _ := json.Marshal(value)
			return fmt.Sprintf("%s|%v", encoded, err)
		})
	}`
	return harness(commonImports("errors", "strings"), declarations, loop, selected)
}

// roomSchedulerExercise defines deterministic minimum-room allocation over half-open intervals.
func roomSchedulerExercise() advancedSpec {
	tests := advancedTests(9,
		[4]string{"Non-overlapping", "Reuses the lowest available room.", `meetings`, `{"a":1,"b":1,"c":1}|<nil>`},
		[4]string{"Overlapping", "Uses the minimum number of rooms.", `meetings`, `{"a":1,"b":2,"c":3}|<nil>`},
		[4]string{"Touching", "Treats an end equal to another start as reusable.", `meetings`, `{"a":1,"b":2,"c":1}|<nil>`},
		[4]string{"Unsorted ties", "Breaks equal starts by end and then ID.", `meetings`, `{"a":2,"b":1,"c":3,"d":1}|<nil>`},
		[4]string{"Invalid interval", "Rejects start greater than or equal to end.", `meetings`, `{}|invalid meeting: bad`},
		[4]string{"Duplicate ID", "Rejects duplicate identities without partial assignments.", `meetings`, `{}|invalid meeting: duplicate a`},
	)
	return advancedSpec{
		id: 9, title: "Room Scheduler", topic: "Algorithms · sorting · priority queues",
		signature: "AssignRooms(meetings []Meeting) (map[string]int, error)",
		objective: "Assign the minimum number of reusable rooms with deterministic room numbers.",
		input:     "Unsorted half-open meetings [Start, End) with unique non-empty IDs and integer times.",
		output:    "A non-nil map from meeting ID to room number, or an empty map and ErrInvalidMeeting with context.",
		constraints: []string{
			"Validate every meeting before assigning any room; Start must be less than End and IDs must be unique/non-empty.",
			"Process meetings by Start, then End, then ID; release every room whose meeting End is <= the next Start.",
			"Always allocate the smallest currently available positive room number.",
		},
		hints: []string{
			"Sort a copy so the caller's slice remains unchanged.",
			"Use one min-heap ordered by end time for occupied rooms and another min-heap of available room numbers.",
			"If no room is available, create the next room number; otherwise pop the smallest available one.",
		},
		pitfalls: []string{"Only releasing one finished room", "Reusing a room when End is greater than Start", "Assigning rooms based on original slice order"},
		builtins: []string{"append", "len", "make"},
		packages: []string{"container/heap", "errors", "fmt", "sort"},
		starter: `import (
	"container/heap"
	"errors"
	"fmt"
	"sort"
)

var ErrInvalidMeeting = errors.New("invalid meeting")

type Meeting struct {
	ID    string
	Start int
	End   int
}

func AssignRooms(meetings []Meeting) (map[string]int, error) {
	// TODO: validate, sort a copy, and assign the smallest reusable room.
	_ = heap.Init
	_ = fmt.Errorf
	_ = sort.Slice
	return map[string]int{}, ErrInvalidMeeting
}
`,
		tests: tests,
		build: buildRoomsHarness,
	}
}

// buildRoomsHarness provides unsorted and invalid meetings so ordering cannot accidentally follow fixture order.
func buildRoomsHarness(selected []VisibleTest) string {
	loop := `var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	cases := [][]Meeting{
		{{ID:"c",Start:5,End:7},{ID:"a",Start:1,End:2},{ID:"b",Start:2,End:4}},
		{{ID:"c",Start:3,End:8},{ID:"a",Start:1,End:6},{ID:"b",Start:2,End:7}},
		{{ID:"c",Start:2,End:5},{ID:"a",Start:0,End:2},{ID:"b",Start:1,End:2}},
		{{ID:"d",Start:3,End:4},{ID:"c",Start:1,End:5},{ID:"a",Start:1,End:3},{ID:"b",Start:1,End:2}},
		{{ID:"bad",Start:4,End:4}},
		{{ID:"a",Start:1,End:2},{ID:"a",Start:3,End:4}},
	}
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, err := AssignRooms(cases[current.Payload.Case])
			encoded, _ := json.Marshal(value)
			return fmt.Sprintf("%s|%v", encoded, err)
		})
	}`
	return harness(commonImports(), advancedCaseDeclarations, loop, selected)
}
