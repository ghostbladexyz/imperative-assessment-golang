package assessment

import (
	"fmt"
	"strings"
)

func baseInstructions(objective, contract, input, output, starter string) Instructions {
	return Instructions{
		Objective:       objective,
		Contract:        contract,
		Input:           input,
		Output:          output,
		StarterNote:     starter,
		Allowed:         []string{"Any Go standard-library package unless specifically restricted."},
		Disallowed:      []string{"Third-party packages", "Network access", "Changing the required API signature"},
		WhitespaceRules: "Expected values use exact JSON or exact text comparison. Line endings are normalized; other whitespace is significant unless the task says otherwise.",
	}
}

func test(id, name, purpose, input, expected string, payload any) VisibleTest {
	return VisibleTest{ID: id, Name: name, Purpose: purpose, Input: input, Expected: expected, payload: payload}
}

func Levels() []Level {
	return []Level{
		levelOne(),
		levelTwo(),
		levelThree(),
		levelFour(),
		levelFive(),
		levelSix(),
		levelSeven(),
		levelEight(),
		levelNine(),
	}
}

func levelOne() Level {
	instructions := baseInstructions(
		"Normalize a line of structured text into reusable tokens.",
		"func NormalizeTokens(input string) []string",
		"An arbitrary UTF-8 string. ASCII letters and digits form tokens; every other rune is a separator.",
		"A non-nil slice of lowercase tokens, in encounter order. Repeated tokens are preserved.",
		"The editor starts nearly blank with the required signature.",
	)
	instructions.Constraints = []string{"Use loops and basic string/rune handling.", "ASCII A-Z becomes a-z.", "An empty result must be []string{}, not nil.", "Do not use regular expressions."}
	instructions.Examples = []Example{{Input: `"Go, GO! 101"`, Output: `["go","go","101"]`}, {Input: `"  one--two  "`, Output: `["one","two"]`}}
	instructions.Documentation = []DocumentationLink{{Label: "strings package", URL: "https://pkg.go.dev/strings"}, {Label: "Unicode rune basics", URL: "https://go.dev/blog/strings"}}
	instructions.Disallowed = append(instructions.Disallowed, "regexp")
	instructions.Hints = []string{"Build the current token one rune at a time.", "Flush a token whenever a rune is not ASCII alphanumeric.", "A small helper can lowercase A-Z by arithmetic, or strings.ToLower can normalize a completed token."}
	instructions.CommonPitfalls = []string{"Returning nil for empty input", "Treating underscore as a token character", "Dropping repeated tokens"}
	tests := []VisibleTest{
		test("l1-empty", "Empty input", "Returns a non-nil empty result.", `""`, `[]`, map[string]any{"input": ""}),
		test("l1-word", "Single word", "Keeps a plain lowercase word.", `"gopher"`, `["gopher"]`, map[string]any{"input": "gopher"}),
		test("l1-case", "Case folding", "Normalizes ASCII uppercase letters.", `"GoLang"`, `["golang"]`, map[string]any{"input": "GoLang"}),
		test("l1-space", "Whitespace separators", "Handles repeated whitespace.", `"  one \t two\nthree "`, `["one","two","three"]`, map[string]any{"input": "  one \t two\nthree "}),
		test("l1-punct", "Punctuation separators", "Splits on punctuation.", `"go,build;test!"`, `["go","build","test"]`, map[string]any{"input": "go,build;test!"}),
		test("l1-digits", "Digits stay in tokens", "Preserves digits within tokens.", `"go1 v2 007"`, `["go1","v2","007"]`, map[string]any{"input": "go1 v2 007"}),
		test("l1-repeat", "Repeated values", "Does not deduplicate.", `"Go go GO"`, `["go","go","go"]`, map[string]any{"input": "Go go GO"}),
		test("l1-symbols", "Only separators", "Returns empty for symbol-only input.", `"---___..."`, `[]`, map[string]any{"input": "---___..."}),
		test("l1-underscore", "Underscore boundary", "Treats underscore as a separator.", `"snake_case"`, `["snake","case"]`, map[string]any{"input": "snake_case"}),
		test("l1-unicode", "Non-ASCII boundary", "Treats non-ASCII runes as separators.", `"caféGo"`, `["caf","go"]`, map[string]any{"input": "caféGo"}),
		test("l1-edges", "Edge punctuation", "Ignores leading and trailing separators.", `".first:last."`, `["first","last"]`, map[string]any{"input": ".first:last."}),
		test("l1-mixed", "Mixed structured text", "Combines all normalization rules.", `"ID=AB-12; id=xy_9"`, `["id","ab","12","id","xy","9"]`, map[string]any{"input": "ID=AB-12; id=xy_9"}),
	}
	return Level{ID: 1, Title: "Token Trail", Topic: "Strings · loops · parsing", Difficulty: "Foundation", Signature: "NormalizeTokens(input string) []string", StarterCode: `func NormalizeTokens(input string) []string {
	// TODO: return lowercase ASCII letter/digit tokens.
	return nil
}
`, Instructions: instructions, Tests: tests, build: buildLevelOne}
}

func buildLevelOne(tests []VisibleTest) string {
	return harness(commonImports(), `type inputCase struct {
	Input string `+"`json:\"input\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, _ := json.Marshal(NormalizeTokens(current.Payload.Input))
			return string(value)
		})
	}`, tests)
}

func levelTwo() Level {
	instructions := baseInstructions(
		"Count meaningful labels while preserving the order in which each label first appears.",
		"func TallyWords(values []string) []WordCount",
		"A slice of strings. Trim surrounding Unicode whitespace and ignore empty trimmed values. Matching is case-sensitive.",
		"A non-nil []WordCount ordered by first appearance, with each trimmed word and its total count.",
		"A partial implementation supplies WordCount and a TODO function.",
	)
	instructions.Constraints = []string{"Preserve first-seen order.", "Trim with strings.TrimSpace.", "Case-sensitive keys.", "Return []WordCount{} for no meaningful values."}
	instructions.Examples = []Example{{Input: `[" red ","blue","red"]`, Output: `[{"word":"red","count":2},{"word":"blue","count":1}]`}, {Input: `["","  ","Go","go"]`, Output: `[{"word":"Go","count":1},{"word":"go","count":1}]`}}
	instructions.Documentation = []DocumentationLink{{Label: "Maps", URL: "https://go.dev/blog/maps"}, {Label: "strings.TrimSpace", URL: "https://pkg.go.dev/strings#TrimSpace"}}
	instructions.Hints = []string{"A map can remember where a word lives in the result slice.", "Append only on first sight; otherwise increment the existing element.", "Initialize the result with make([]WordCount, 0) so empty input is not nil."}
	instructions.CommonPitfalls = []string{"Iterating over the map for output", "Counting blank strings", "Normalizing case when the contract says not to"}
	tests := []VisibleTest{
		test("l2-empty", "Empty slice", "Produces a non-nil empty result.", `[]`, `[]`, map[string]any{"values": []string{}}),
		test("l2-one", "One label", "Counts a single value.", `["red"]`, `[{"word":"red","count":1}]`, map[string]any{"values": []string{"red"}}),
		test("l2-repeat", "Repeated label", "Increments an existing count.", `["red","red","red"]`, `[{"word":"red","count":3}]`, map[string]any{"values": []string{"red", "red", "red"}}),
		test("l2-order", "Stable first order", "Keeps first appearance order.", `["b","a","b","c","a"]`, `[{"word":"b","count":2},{"word":"a","count":2},{"word":"c","count":1}]`, map[string]any{"values": []string{"b", "a", "b", "c", "a"}}),
		test("l2-trim", "Trim surrounding space", "Uses trimmed labels as keys.", `[" red","red "," red "]`, `[{"word":"red","count":3}]`, map[string]any{"values": []string{" red", "red ", " red "}}),
		test("l2-blanks", "Ignore blanks", "Skips empty and whitespace-only entries.", `[""," ","x","\t"]`, `[{"word":"x","count":1}]`, map[string]any{"values": []string{"", " ", "x", "\t"}}),
		test("l2-case", "Case sensitivity", "Keeps differently cased labels distinct.", `["Go","go","GO"]`, `[{"word":"Go","count":1},{"word":"go","count":1},{"word":"GO","count":1}]`, map[string]any{"values": []string{"Go", "go", "GO"}}),
		test("l2-unicode-space", "Unicode whitespace", "Trims Unicode whitespace.", `["\u00a0tea\u00a0","tea"]`, `[{"word":"tea","count":2}]`, map[string]any{"values": []string{"\u00a0tea\u00a0", "tea"}}),
		test("l2-punctuation", "Punctuation is data", "Does not tokenize label contents.", `["a-b","a","a-b"]`, `[{"word":"a-b","count":2},{"word":"a","count":1}]`, map[string]any{"values": []string{"a-b", "a", "a-b"}}),
		test("l2-numbers", "Numeric-looking labels", "Treats numbers as ordinary strings.", `["01","1","01"]`, `[{"word":"01","count":2},{"word":"1","count":1}]`, map[string]any{"values": []string{"01", "1", "01"}}),
		test("l2-late-repeat", "Late repeat", "Updates without moving an item.", `["a","b","c","a"]`, `[{"word":"a","count":2},{"word":"b","count":1},{"word":"c","count":1}]`, map[string]any{"values": []string{"a", "b", "c", "a"}}),
		test("l2-many", "Several groups", "Handles a denser stable tally.", `["x"," y ","z","x","y","x"]`, `[{"word":"x","count":3},{"word":"y","count":2},{"word":"z","count":1}]`, map[string]any{"values": []string{"x", " y ", "z", "x", "y", "x"}}),
	}
	return Level{ID: 2, Title: "First Seen Ledger", Topic: "Slices · maps · stable output", Difficulty: "Foundation+", Signature: "TallyWords(values []string) []WordCount", StarterCode: `import "strings"

type WordCount struct {
	Word  string ` + "`json:\"word\"`" + `
	Count int    ` + "`json:\"count\"`" + `
}

func TallyWords(values []string) []WordCount {
	// TODO: trim, ignore blanks, count, and preserve first appearance.
	_ = strings.TrimSpace
	return nil
}
`, Instructions: instructions, Tests: tests, build: buildLevelTwo}
}

func buildLevelTwo(tests []VisibleTest) string {
	return harness(commonImports(), `type inputCase struct {
	Values []string `+"`json:\"values\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, _ := json.Marshal(TallyWords(current.Payload.Values))
			return string(value)
		})
	}`, tests)
}

func levelThree() Level {
	instructions := baseInstructions(
		"Parse a network port with a stable error contract.",
		"func ParsePort(input string) (int, error)",
		"A string containing optional surrounding whitespace and otherwise only base-10 digits.",
		"The port as an int when it is in 1..65535. Every invalid input must return an error matching ErrInvalidPort through errors.Is.",
		"Imports and the sentinel error are supplied; the parser is unfinished.",
	)
	instructions.Constraints = []string{"Trim surrounding whitespace.", "Reject signs, decimals, embedded whitespace, zero, and values above 65535.", "Wrap or return ErrInvalidPort.", "Do not panic."}
	instructions.Examples = []Example{{Input: `" 8080 "`, Output: `8080, nil`}, {Input: `"0"`, Output: `0, error matching ErrInvalidPort`}}
	instructions.Documentation = []DocumentationLink{{Label: "strconv.Atoi", URL: "https://pkg.go.dev/strconv#Atoi"}, {Label: "errors.Is", URL: "https://pkg.go.dev/errors#Is"}, {Label: "fmt.Errorf wrapping", URL: "https://pkg.go.dev/fmt#Errorf"}}
	instructions.Hints = []string{"Trim before checking the text.", "Check every remaining byte is between '0' and '9' before Atoi.", "Use fmt.Errorf(\"...: %w\", ErrInvalidPort) if you want context while preserving errors.Is."}
	instructions.CommonPitfalls = []string{"Letting strconv accept +80", "Accepting port zero", "Returning an unrelated parse error"}
	tests := []VisibleTest{
		test("l3-min", "Minimum port", "Accepts the lowest valid value.", `"1"`, `ok:1`, map[string]any{"input": "1"}),
		test("l3-common", "Common HTTP port", "Parses an ordinary port.", `"8080"`, `ok:8080`, map[string]any{"input": "8080"}),
		test("l3-max", "Maximum port", "Accepts the highest valid value.", `"65535"`, `ok:65535`, map[string]any{"input": "65535"}),
		test("l3-trim", "Surrounding whitespace", "Trims before parsing.", `" \t443\n"`, `ok:443`, map[string]any{"input": " \t443\n"}),
		test("l3-empty", "Empty input", "Returns the sentinel category.", `""`, `error:ErrInvalidPort`, map[string]any{"input": ""}),
		test("l3-space", "Whitespace only", "Rejects blank input.", `"   "`, `error:ErrInvalidPort`, map[string]any{"input": "   "}),
		test("l3-zero", "Port zero", "Rejects zero.", `"0"`, `error:ErrInvalidPort`, map[string]any{"input": "0"}),
		test("l3-high", "Above range", "Rejects 65536.", `"65536"`, `error:ErrInvalidPort`, map[string]any{"input": "65536"}),
		test("l3-negative", "Negative sign", "Rejects signed text.", `"-1"`, `error:ErrInvalidPort`, map[string]any{"input": "-1"}),
		test("l3-plus", "Positive sign", "Rejects an explicit plus sign.", `"+80"`, `error:ErrInvalidPort`, map[string]any{"input": "+80"}),
		test("l3-decimal", "Decimal syntax", "Rejects non-integers.", `"80.0"`, `error:ErrInvalidPort`, map[string]any{"input": "80.0"}),
		test("l3-inner-space", "Embedded whitespace", "Rejects split digits.", `"8 0"`, `error:ErrInvalidPort`, map[string]any{"input": "8 0"}),
		test("l3-text", "Alphabetic text", "Rejects non-numeric input.", `"http"`, `error:ErrInvalidPort`, map[string]any{"input": "http"}),
		test("l3-overflow", "Integer overflow", "Returns the sentinel rather than leaking strconv errors.", `"999999999999999999999999"`, `error:ErrInvalidPort`, map[string]any{"input": "999999999999999999999999"}),
	}
	return Level{ID: 3, Title: "Safe Harbor", Topic: "Validation · strconv · errors", Difficulty: "Intermediate", Signature: "ParsePort(input string) (int, error)", StarterCode: `import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidPort = errors.New("invalid port")

func ParsePort(input string) (int, error) {
	// TODO: validate the syntax, parse the value, and preserve ErrInvalidPort.
	_, _ = strconv.Atoi, strings.TrimSpace
	return 0, ErrInvalidPort
}
`, Instructions: instructions, Tests: tests, build: buildLevelThree}
}

func buildLevelThree(tests []VisibleTest) string {
	return harness(commonImports("errors"), `type inputCase struct {
	Input string `+"`json:\"input\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, err := ParsePort(current.Payload.Input)
			if err == nil {
				return fmt.Sprintf("ok:%d", value)
			}
			if errors.Is(err, ErrInvalidPort) {
				return "error:ErrInvalidPort"
			}
			return "error:other"
		})
	}`, tests)
}

func levelFour() Level {
	instructions := baseInstructions(
		"Aggregate a compact inventory feed into stable SKU totals.",
		"func SummarizeInventory(r io.Reader) ([]ItemTotal, error)",
		"UTF-8 lines in sku,quantity,unit_price_cents form. Blank lines and lines whose trimmed form starts with # are ignored.",
		"Totals sorted by SKU ascending. Quantity and revenue are summed. Any malformed record returns an error and no partial result is required.",
		"Structs are supplied and the transformation starts with a TODO.",
	)
	instructions.Constraints = []string{"Exactly three comma-separated fields.", "Trim every field.", "SKU must be non-empty; quantity and price must be non-negative integers.", "Sort by SKU.", "Scanner errors must be returned."}
	instructions.Examples = []Example{{Input: `"tea,2,150\ntea,1,150"`, Output: `[{"sku":"tea","quantity":3,"revenueCents":450}]`}, {Input: `"# note\n\npen,4,25"`, Output: `[{"sku":"pen","quantity":4,"revenueCents":100}]`}}
	instructions.Documentation = []DocumentationLink{{Label: "bufio.Scanner", URL: "https://pkg.go.dev/bufio#Scanner"}, {Label: "sort.Slice", URL: "https://pkg.go.dev/sort#Slice"}, {Label: "io.Reader", URL: "https://pkg.go.dev/io#Reader"}}
	instructions.Hints = []string{"Scan one line at a time and keep a map keyed by SKU.", "Store a struct back into the map after updating it.", "After scanning, copy map values into a slice and sort it."}
	instructions.CommonPitfalls = []string{"Forgetting Scanner.Err", "Multiplying after totals have already been aggregated", "Returning map iteration order"}
	tests := []VisibleTest{
		test("l4-empty", "Empty feed", "Returns an empty result.", `""`, `[]`, map[string]any{"text": ""}),
		test("l4-one", "Single record", "Calculates revenue.", `"pen,4,25"`, `[{"sku":"pen","quantity":4,"revenueCents":100}]`, map[string]any{"text": "pen,4,25"}),
		test("l4-aggregate", "Repeated SKU", "Combines repeated records.", `"tea,2,150\ntea,1,175"`, `[{"sku":"tea","quantity":3,"revenueCents":475}]`, map[string]any{"text": "tea,2,150\ntea,1,175"}),
		test("l4-sort", "Sorted output", "Sorts independent SKUs.", `"z,1,1\na,1,2\nm,1,3"`, `[{"sku":"a","quantity":1,"revenueCents":2},{"sku":"m","quantity":1,"revenueCents":3},{"sku":"z","quantity":1,"revenueCents":1}]`, map[string]any{"text": "z,1,1\na,1,2\nm,1,3"}),
		test("l4-comments", "Comments and blanks", "Ignores non-record lines.", `"# stock\n\npen,2,5\n  # end"`, `[{"sku":"pen","quantity":2,"revenueCents":10}]`, map[string]any{"text": "# stock\n\npen,2,5\n  # end"}),
		test("l4-trim", "Trimmed fields", "Trims all three fields.", `" pen , 2 , 5 "`, `[{"sku":"pen","quantity":2,"revenueCents":10}]`, map[string]any{"text": " pen , 2 , 5 "}),
		test("l4-zero", "Zero values", "Allows non-negative zero values.", `"free,3,0\nnone,0,99"`, `[{"sku":"free","quantity":3,"revenueCents":0},{"sku":"none","quantity":0,"revenueCents":0}]`, map[string]any{"text": "free,3,0\nnone,0,99"}),
		test("l4-fields", "Wrong field count", "Rejects missing fields.", `"pen,2"`, `error`, map[string]any{"text": "pen,2"}),
		test("l4-empty-sku", "Empty SKU", "Rejects a blank identifier.", `" ,2,5"`, `error`, map[string]any{"text": " ,2,5"}),
		test("l4-bad-qty", "Malformed quantity", "Rejects non-numeric quantity.", `"pen,two,5"`, `error`, map[string]any{"text": "pen,two,5"}),
		test("l4-negative", "Negative value", "Rejects negative numbers.", `"pen,-1,5"`, `error`, map[string]any{"text": "pen,-1,5"}),
		test("l4-bad-price", "Malformed price", "Rejects a decimal price.", `"pen,2,1.5"`, `error`, map[string]any{"text": "pen,2,1.5"}),
		test("l4-mixed-error", "Error after valid row", "Does not silently ignore later damage.", `"pen,1,5\nbroken\ntea,1,5"`, `error`, map[string]any{"text": "pen,1,5\nbroken\ntea,1,5"}),
	}
	return Level{ID: 4, Title: "Stockroom Stream", Topic: "Scanner · structs · aggregation", Difficulty: "Intermediate", Signature: "SummarizeInventory(r io.Reader) ([]ItemTotal, error)", StarterCode: `import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type ItemTotal struct {
	SKU          string ` + "`json:\"sku\"`" + `
	Quantity     int    ` + "`json:\"quantity\"`" + `
	RevenueCents int    ` + "`json:\"revenueCents\"`" + `
}

func SummarizeInventory(r io.Reader) ([]ItemTotal, error) {
	// TODO: scan, validate, aggregate, and sort.
	_, _, _, _, _ = bufio.NewScanner, fmt.Errorf, sort.Slice, strconv.Atoi, strings.TrimSpace
	return nil, nil
}
`, Instructions: instructions, Tests: tests, build: buildLevelFour}
}

func buildLevelFour(tests []VisibleTest) string {
	return harness(commonImports("strings"), `type inputCase struct {
	Text string `+"`json:\"text\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			value, err := SummarizeInventory(strings.NewReader(current.Payload.Text))
			if err != nil {
				return "error"
			}
			encoded, _ := json.Marshal(value)
			return string(encoded)
		})
	}`, tests)
}

func levelFive() Level {
	instructions := baseInstructions(
		"Implement a small JSON greeting endpoint with deliberate HTTP behavior.",
		"func NewGreetingHandler() http.Handler",
		"POST /greet with Content-Type application/json and exactly one JSON object containing a non-blank name. Unknown JSON fields are invalid.",
		"On success return 200, application/json, and {\"message\":\"Hello, NAME!\"}. Use 400 for invalid input, 405 plus Allow: POST for wrong methods, and 404 for other paths.",
		"An unfinished handler factory is provided.",
	)
	instructions.Constraints = []string{"No external server or network.", "Trim the name before using it.", "Reject trailing JSON values and unknown fields.", "All response bodies must be JSON objects.", "Set Content-Type before WriteHeader."}
	instructions.Examples = []Example{{Input: `POST /greet {"name":"Ada"}`, Output: `200 {"message":"Hello, Ada!"}`}, {Input: `GET /greet`, Output: `405, Allow: POST`}}
	instructions.Documentation = []DocumentationLink{{Label: "net/http", URL: "https://pkg.go.dev/net/http"}, {Label: "httptest", URL: "https://pkg.go.dev/net/http/httptest"}, {Label: "encoding/json Decoder", URL: "https://pkg.go.dev/encoding/json#Decoder"}}
	instructions.Hints = []string{"Check the path, then method, before decoding.", "Use json.Decoder.DisallowUnknownFields.", "Decode once into the request, then ensure a second decode returns io.EOF."}
	instructions.CommonPitfalls = []string{"Writing status before headers", "Accepting whitespace-only names", "Returning plain-text errors"}
	tests := makeHTTPTests()
	return Level{ID: 5, Title: "Hello, Handler", Topic: "HTTP · JSON · validation", Difficulty: "Intermediate+", Signature: "NewGreetingHandler() http.Handler", StarterCode: `import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func NewGreetingHandler() http.Handler {
	// TODO: implement POST /greet.
	_, _, _ = json.NewDecoder, io.EOF, strings.TrimSpace
	return http.NotFoundHandler()
}
`, Instructions: instructions, Tests: tests, build: buildLevelFive}
}

func makeHTTPTests() []VisibleTest {
	type httpPayload map[string]any
	cases := []struct {
		id, name, purpose, method, path, body, expected string
	}{
		{"l5-basic", "Basic greeting", "Returns the documented response.", "POST", "/greet", `{"name":"Ada"}`, `{"status":200,"contentType":"application/json","body":{"message":"Hello, Ada!"}}`},
		{"l5-trim", "Trimmed name", "Trims before building the greeting.", "POST", "/greet", `{"name":"  Lin  "}`, `{"status":200,"contentType":"application/json","body":{"message":"Hello, Lin!"}}`},
		{"l5-unicode", "Unicode name", "Preserves valid UTF-8 text.", "POST", "/greet", `{"name":"Zoë"}`, `{"status":200,"contentType":"application/json","body":{"message":"Hello, Zoë!"}}`},
		{"l5-get", "Wrong method", "Returns 405 and the Allow header.", "GET", "/greet", "", `{"status":405,"contentType":"application/json","allow":"POST"}`},
		{"l5-put", "Another wrong method", "Handles all non-POST methods consistently.", "PUT", "/greet", `{}`, `{"status":405,"contentType":"application/json","allow":"POST"}`},
		{"l5-path", "Unknown path", "Returns 404 outside /greet.", "POST", "/other", `{"name":"Ada"}`, `{"status":404,"contentType":"application/json"}`},
		{"l5-empty", "Empty name", "Rejects an empty value.", "POST", "/greet", `{"name":""}`, `{"status":400,"contentType":"application/json"}`},
		{"l5-blank", "Blank name", "Rejects whitespace-only names.", "POST", "/greet", `{"name":"   "}`, `{"status":400,"contentType":"application/json"}`},
		{"l5-malformed", "Malformed JSON", "Returns 400 for broken JSON.", "POST", "/greet", `{"name":`, `{"status":400,"contentType":"application/json"}`},
		{"l5-missing", "Missing name", "Rejects an object without name.", "POST", "/greet", `{}`, `{"status":400,"contentType":"application/json"}`},
		{"l5-unknown", "Unknown field", "Rejects fields outside the contract.", "POST", "/greet", `{"name":"Ada","admin":true}`, `{"status":400,"contentType":"application/json"}`},
		{"l5-trailing", "Trailing JSON", "Rejects multiple JSON values.", "POST", "/greet", `{"name":"Ada"} {}`, `{"status":400,"contentType":"application/json"}`},
	}
	out := make([]VisibleTest, 0, len(cases))
	for _, item := range cases {
		out = append(out, test(item.id, item.name, item.purpose, fmt.Sprintf("%s %s %s", item.method, item.path, item.body), item.expected, httpPayload{"method": item.method, "path": item.path, "body": item.body}))
	}
	return out
}

func buildLevelFive(tests []VisibleTest) string {
	return harness(commonImports("net/http/httptest", "strings"), `type inputCase struct {
	Method string `+"`json:\"method\"`"+`
	Path string `+"`json:\"path\"`"+`
	Body string `+"`json:\"body\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}
type responseView struct {
	Status int `+"`json:\"status\"`"+`
	ContentType string `+"`json:\"contentType\"`"+`
	Allow string `+"`json:\"allow,omitempty\"`"+`
	Body map[string]any `+"`json:\"body,omitempty\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	handler := NewGreetingHandler()
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			request := httptest.NewRequest(current.Payload.Method, current.Payload.Path, strings.NewReader(current.Payload.Body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			view := responseView{Status: recorder.Code, ContentType: recorder.Header().Get("Content-Type"), Allow: recorder.Header().Get("Allow")}
			if recorder.Code == 200 {
				_ = json.Unmarshal(recorder.Body.Bytes(), &view.Body)
			}
			encoded, _ := json.Marshal(view)
			return string(encoded)
		})
	}`, tests)
}

func levelSix() Level {
	instructions := baseInstructions(
		"Find every unique value pair that reaches a target sum.",
		"func PairSums(values []int, target int) [][2]int",
		"An integer slice and an integer target. The input may be unsorted and contain duplicates.",
		"A non-nil slice of unique pairs [a,b] where a <= b and a+b == target, sorted by a then b.",
		"The signature and a TODO are supplied.",
	)
	instructions.Constraints = []string{"Expected O(n log n) time or O(n) average time.", "Do not mutate the caller's input.", "Pairs are unique by value, not index.", "Return [][2]int{} when there are no pairs."}
	instructions.Examples = []Example{{Input: `[3,1,2,2,4], target 5`, Output: `[[1,4],[2,3]]`}, {Input: `[2,2,2], target 4`, Output: `[[2,2]]`}}
	instructions.Documentation = []DocumentationLink{{Label: "sort.Ints", URL: "https://pkg.go.dev/sort#Ints"}, {Label: "Go slices", URL: "https://go.dev/blog/slices-intro"}}
	instructions.Hints = []string{"Copy and sort the input, then use a left and right pointer.", "When you find a pair, skip every duplicate of both values.", "A pair using the same value requires at least two occurrences, which two pointers naturally enforce."}
	instructions.CommonPitfalls = []string{"Returning duplicate pairs", "Mutating values while sorting", "Returning pairs in discovery order"}
	tests := makePairTests()
	return Level{ID: 6, Title: "Balanced Pairs", Topic: "Algorithms · sorting · two pointers", Difficulty: "Advanced", Signature: "PairSums(values []int, target int) [][2]int", StarterCode: `func PairSums(values []int, target int) [][2]int {
	// TODO: return unique, sorted value pairs without changing values.
	return nil
}
`, Instructions: instructions, Tests: tests, build: buildLevelSix}
}

func makePairTests() []VisibleTest {
	cases := []struct {
		id, name, purpose string
		values            []int
		target            int
		expected          string
	}{
		{"l6-empty", "Empty input", "Returns an empty slice.", []int{}, 1, `[]`},
		{"l6-short", "One value", "Does not reuse one element.", []int{2}, 4, `[]`},
		{"l6-basic", "One pair", "Finds a simple pair.", []int{1, 4}, 5, `[[1,4]]`},
		{"l6-many", "Several pairs", "Returns sorted unique pairs.", []int{3, 1, 2, 2, 4}, 5, `[[1,4],[2,3]]`},
		{"l6-duplicates", "Duplicate pair values", "Emits a value pair once.", []int{1, 4, 1, 4, 1}, 5, `[[1,4]]`},
		{"l6-same", "Equal pair", "Uses two equal elements when available.", []int{2, 2, 2}, 4, `[[2,2]]`},
		{"l6-not-same", "Single midpoint", "Does not reuse one midpoint.", []int{2, 1, 5}, 4, `[]`},
		{"l6-negative", "Negative values", "Supports signed integers.", []int{-3, 7, -1, 5, 1}, 4, `[[-3,7],[-1,5]]`},
		{"l6-zero", "Zero target", "Finds pairs around zero.", []int{-2, 2, 0, 0, -1, 1}, 0, `[[-2,2],[-1,1],[0,0]]`},
		{"l6-none", "No match", "Returns empty when no sum matches.", []int{1, 2, 3}, 99, `[]`},
		{"l6-unsorted", "Strongly unsorted", "Output is independent of input order.", []int{10, -5, 3, 7, 0, 5}, 10, `[[0,10],[3,7]]`},
		{"l6-bound", "Integer boundaries", "Handles large signed values without subtraction overflow tricks.", []int{-1000000, 1000000, 1, -1}, 0, `[[-1000000,1000000],[-1,1]]`},
		{"l6-copy", "Input remains unchanged", "Protects caller-owned order.", []int{4, 1, 3, 2}, 5, `[[1,4],[2,3]]`},
		{"l6-dense", "Dense duplicates", "Skips duplicates on both sides.", []int{1, 1, 2, 2, 3, 3, 4, 4}, 5, `[[1,4],[2,3]]`},
	}
	out := make([]VisibleTest, 0, len(cases))
	for _, item := range cases {
		out = append(out, test(item.id, item.name, item.purpose, fmt.Sprintf("%v, target %d", item.values, item.target), item.expected, map[string]any{"values": item.values, "target": item.target}))
	}
	return out
}

func buildLevelSix(tests []VisibleTest) string {
	return harness(commonImports(), `type inputCase struct {
	Values []int `+"`json:\"values\"`"+`
	Target int `+"`json:\"target\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			before, _ := json.Marshal(current.Payload.Values)
			value, _ := json.Marshal(PairSums(current.Payload.Values, current.Payload.Target))
			after, _ := json.Marshal(current.Payload.Values)
			if string(before) != string(after) {
				return "input mutated"
			}
			return string(value)
		})
	}`, tests)
}

func levelSeven() Level {
	instructions := baseInstructions(
		"Run work with a strict concurrency bound while preserving input order.",
		"func MapOrdered(ctx context.Context, inputs []int, workers int, fn func(context.Context, int) (int, error)) ([]int, error)",
		"A context, input integers, a positive worker limit, and a supplied function.",
		"Results in input order. At most workers calls may run at once. On the first function error or context cancellation, stop scheduling, cancel in-flight work, wait for workers, and return the error.",
		"A sentinel for invalid worker counts and the required function are supplied.",
	)
	instructions.Constraints = []string{"workers <= 0 returns ErrInvalidWorkers.", "Empty input returns []int{} without calling fn.", "Do not leak goroutines.", "Preserve the original function error so errors.Is works.", "Do not spawn one unbounded goroutine per item."}
	instructions.Examples = []Example{{Input: `inputs [3,1,2], workers 2, fn doubles`, Output: `[6,2,4]`}, {Input: `workers 0`, Output: `error matching ErrInvalidWorkers`}}
	instructions.Documentation = []DocumentationLink{{Label: "context", URL: "https://pkg.go.dev/context"}, {Label: "sync.WaitGroup", URL: "https://pkg.go.dev/sync#WaitGroup"}, {Label: "Go pipelines", URL: "https://go.dev/blog/pipelines"}}
	instructions.Hints = []string{"Create a child context you can cancel on the first error.", "Send indexed jobs through a channel to a fixed number of workers.", "Close jobs, wait for every worker, then prefer the first worker error over context cancellation."}
	instructions.CommonPitfalls = []string{"Appending results in completion order", "Returning before workers stop", "Losing the original error behind context.Canceled"}
	tests := makeConcurrencyTests()
	return Level{ID: 7, Title: "Ordered Workshop", Topic: "Goroutines · channels · cancellation", Difficulty: "Advanced", Signature: "MapOrdered(ctx context.Context, inputs []int, workers int, fn func(context.Context, int) (int, error)) ([]int, error)", StarterCode: `import (
	"context"
	"errors"
)

var ErrInvalidWorkers = errors.New("workers must be positive")

func MapOrdered(
	ctx context.Context,
	inputs []int,
	workers int,
	fn func(context.Context, int) (int, error),
) ([]int, error) {
	// TODO: build a bounded worker pool with ordered results.
	return nil, ErrInvalidWorkers
}
`, Instructions: instructions, Tests: tests, build: buildLevelSeven}
}

func makeConcurrencyTests() []VisibleTest {
	cases := []struct {
		id, name, purpose, expected string
		payload                     map[string]any
	}{
		{"l7-empty", "Empty input", "Returns an empty non-nil slice.", `[]`, map[string]any{"inputs": []int{}, "workers": 2, "mode": "double"}},
		{"l7-invalid-zero", "Zero workers", "Returns the documented sentinel.", `error:ErrInvalidWorkers`, map[string]any{"inputs": []int{1}, "workers": 0, "mode": "double"}},
		{"l7-invalid-negative", "Negative workers", "Rejects any negative bound.", `error:ErrInvalidWorkers`, map[string]any{"inputs": []int{1}, "workers": -2, "mode": "double"}},
		{"l7-one", "Single worker", "Works sequentially.", `[2,4,6]`, map[string]any{"inputs": []int{1, 2, 3}, "workers": 1, "mode": "double"}},
		{"l7-order", "Completion order differs", "Still returns input order.", `[10,2,6,4]`, map[string]any{"inputs": []int{5, 1, 3, 2}, "workers": 3, "mode": "skew"}},
		{"l7-many-workers", "More workers than jobs", "Handles an oversized valid bound.", `[4,8]`, map[string]any{"inputs": []int{2, 4}, "workers": 10, "mode": "double"}},
		{"l7-bound-one", "Bound of one", "Never exceeds one active call.", `bound:ok`, map[string]any{"inputs": []int{1, 2, 3, 4}, "workers": 1, "mode": "bound"}},
		{"l7-bound-three", "Bound of three", "Never exceeds the supplied bound.", `bound:ok`, map[string]any{"inputs": []int{1, 2, 3, 4, 5, 6}, "workers": 3, "mode": "bound"}},
		{"l7-error", "Function error", "Preserves the worker error.", `error:worker`, map[string]any{"inputs": []int{1, 2, 3, 4}, "workers": 2, "mode": "error", "errorAt": 3}},
		{"l7-context", "Cancelled context", "Returns cancellation promptly.", `error:context`, map[string]any{"inputs": []int{1, 2}, "workers": 2, "mode": "cancelled"}},
		{"l7-values", "Negative values", "Treats inputs as opaque work items.", `[-4,0,10]`, map[string]any{"inputs": []int{-2, 0, 5}, "workers": 2, "mode": "double"}},
		{"l7-deterministic", "Repeated values", "Retains every input position.", `[6,6,6,6]`, map[string]any{"inputs": []int{3, 3, 3, 3}, "workers": 2, "mode": "skew"}},
	}
	out := make([]VisibleTest, 0, len(cases))
	for _, item := range cases {
		out = append(out, test(item.id, item.name, item.purpose, fmt.Sprint(item.payload), item.expected, item.payload))
	}
	return out
}

func buildLevelSeven(tests []VisibleTest) string {
	return harness(commonImports("context", "errors", "sync", "sync/atomic"), `var assessmentWorkerError = errors.New("worker failed")
type inputCase struct {
	Inputs []int `+"`json:\"inputs\"`"+`
	Workers int `+"`json:\"workers\"`"+`
	Mode string `+"`json:\"mode\"`"+`
	ErrorAt int `+"`json:\"errorAt\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			ctx := context.Background()
			if current.Payload.Mode == "cancelled" {
				cancelled, cancel := context.WithCancel(ctx)
				cancel()
				ctx = cancelled
			}
			var active int32
			var exceeded int32
			var gate sync.Mutex
			fn := func(callCtx context.Context, value int) (int, error) {
				now := atomic.AddInt32(&active, 1)
				if now > int32(current.Payload.Workers) && current.Payload.Workers > 0 {
					atomic.StoreInt32(&exceeded, 1)
				}
				defer atomic.AddInt32(&active, -1)
				if current.Payload.Mode == "cancelled" {
					<-callCtx.Done()
					return 0, callCtx.Err()
				}
				if current.Payload.Mode == "error" && value == current.Payload.ErrorAt {
					return 0, assessmentWorkerError
				}
				if current.Payload.Mode == "skew" {
					gate.Lock()
					for spin := 0; spin < (6-value)*300; spin++ {}
					gate.Unlock()
				}
				return value * 2, nil
			}
			values, err := MapOrdered(ctx, current.Payload.Inputs, current.Payload.Workers, fn)
			if errors.Is(err, ErrInvalidWorkers) { return "error:ErrInvalidWorkers" }
			if errors.Is(err, assessmentWorkerError) { return "error:worker" }
			if errors.Is(err, context.Canceled) { return "error:context" }
			if err != nil { return "error:other" }
			if current.Payload.Mode == "bound" {
				if atomic.LoadInt32(&exceeded) != 0 { return "bound:exceeded" }
				return "bound:ok"
			}
			encoded, _ := json.Marshal(values)
			return string(encoded)
		})
	}`, tests)
}

func levelEight() Level {
	instructions := baseInstructions(
		"Query account totals through a small database boundary that behaves like database/sql.",
		"func LoadAccountTotals(ctx context.Context, db Queryer, minCents int64) ([]AccountTotal, error)",
		"A context, supplied Queryer, and minimum total. Execute the exact documented parameterized query and scan id, name, total_cents rows.",
		"Rows in database order, with query/scan/iteration/close errors propagated. Always close rows.",
		"Database-shaped interfaces and result structs are supplied. No database server or third-party driver is needed.",
	)
	instructions.Constraints = []string{"Query: SELECT id, name, total_cents FROM account_totals WHERE total_cents >= ? ORDER BY id", "Pass minCents as the sole argument; never interpolate it.", "defer rows.Close immediately after a successful query.", "Check rows.Err after iteration.", "Return []AccountTotal{} for zero rows."}
	instructions.Examples = []Example{{Input: `minCents 500, rows (1,"Ada",900)`, Output: `[{"id":1,"name":"Ada","totalCents":900}]`}, {Input: `query returns error`, Output: `that error`}}
	instructions.Documentation = []DocumentationLink{{Label: "database/sql querying", URL: "https://go.dev/doc/database/querying"}, {Label: "Rows.Next", URL: "https://pkg.go.dev/database/sql#Rows.Next"}, {Label: "Context", URL: "https://pkg.go.dev/context"}}
	instructions.Hints = []string{"Call QueryContext with the constant query and minCents as a separate argument.", "Loop while rows.Next and scan into a fresh AccountTotal.", "Close with defer, then check rows.Err after the loop."}
	instructions.CommonPitfalls = []string{"Formatting minCents into SQL", "Forgetting Close", "Ignoring Scan or Rows.Err errors"}
	tests := makeSQLTests()
	return Level{ID: 8, Title: "Account Rollup", Topic: "SQL boundary · scanning · errors", Difficulty: "Stretch", Stretch: true, Signature: "LoadAccountTotals(ctx context.Context, db Queryer, minCents int64) ([]AccountTotal, error)", StarterCode: `import "context"

type RowSet interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close() error
}

type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (RowSet, error)
}

type AccountTotal struct {
	ID         int64  ` + "`json:\"id\"`" + `
	Name       string ` + "`json:\"name\"`" + `
	TotalCents int64  ` + "`json:\"totalCents\"`" + `
}

func LoadAccountTotals(ctx context.Context, db Queryer, minCents int64) ([]AccountTotal, error) {
	// TODO: perform the parameterized query, scan rows, and propagate errors.
	return nil, nil
}
`, Instructions: instructions, Tests: tests, build: buildLevelEight}
}

func makeSQLTests() []VisibleTest {
	cases := []struct {
		id, name, purpose, expected string
		payload                     map[string]any
	}{
		{"l8-empty", "No rows", "Returns an empty non-nil slice.", `[]`, map[string]any{"min": 0, "rows": []any{}}},
		{"l8-one", "One row", "Scans every selected column.", `[{"id":1,"name":"Ada","totalCents":900}]`, map[string]any{"min": 500, "rows": []any{[]any{1, "Ada", 900}}}},
		{"l8-many", "Several rows", "Preserves database ordering.", `[{"id":1,"name":"Ada","totalCents":500},{"id":4,"name":"Lin","totalCents":1250}]`, map[string]any{"min": 500, "rows": []any{[]any{1, "Ada", 500}, []any{4, "Lin", 1250}}}},
		{"l8-zero-min", "Zero parameter", "Passes zero as a query argument.", `[]`, map[string]any{"min": 0, "rows": []any{}}},
		{"l8-negative-min", "Negative parameter", "Does not impose undocumented validation.", `[]`, map[string]any{"min": -50, "rows": []any{}}},
		{"l8-query", "Exact query", "Uses the documented stable query.", `query:ok`, map[string]any{"min": 75, "rows": []any{}, "check": "query"}},
		{"l8-arg", "Parameterized argument", "Keeps values outside SQL text.", `arg:ok`, map[string]any{"min": 123456, "rows": []any{}, "check": "arg"}},
		{"l8-query-error", "Query error", "Propagates a query failure.", `error:query`, map[string]any{"min": 0, "queryError": true}},
		{"l8-scan-error", "Scan error", "Stops and returns a scan failure.", `error:scan`, map[string]any{"min": 0, "rows": []any{[]any{1, "Ada", 3}}, "scanErrorAt": 0}},
		{"l8-rows-error", "Iteration error", "Checks Rows.Err after iteration.", `error:rows`, map[string]any{"min": 0, "rows": []any{}, "rowsError": true}},
		{"l8-close-success", "Close on success", "Always closes successful row sets.", `close:ok`, map[string]any{"min": 0, "rows": []any{}, "check": "close"}},
		{"l8-close-error-path", "Close after scan error", "Defers close before scanning.", `close:ok`, map[string]any{"min": 0, "rows": []any{[]any{1, "Ada", 3}}, "scanErrorAt": 0, "check": "close"}},
	}
	out := make([]VisibleTest, 0, len(cases))
	for _, item := range cases {
		out = append(out, test(item.id, item.name, item.purpose, fmt.Sprint(item.payload), item.expected, item.payload))
	}
	return out
}

func buildLevelEight(tests []VisibleTest) string {
	return harness(commonImports("context", "errors"), `var (
	queryFailure = errors.New("query failure")
	scanFailure = errors.New("scan failure")
	rowsFailure = errors.New("rows failure")
)
const expectedQuery = "SELECT id, name, total_cents FROM account_totals WHERE total_cents >= ? ORDER BY id"
type inputCase struct {
	Min int64 `+"`json:\"min\"`"+`
	Rows [][]any `+"`json:\"rows\"`"+`
	Check string `+"`json:\"check\"`"+`
	QueryError bool `+"`json:\"queryError\"`"+`
	ScanErrorAt *int `+"`json:\"scanErrorAt\"`"+`
	RowsError bool `+"`json:\"rowsError\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}
type fakeDB struct { test *inputCase; query string; args []any; rows *fakeRows }
func (db *fakeDB) QueryContext(_ context.Context, query string, args ...any) (RowSet, error) {
	db.query, db.args = query, args
	if db.test.QueryError { return nil, queryFailure }
	db.rows = &fakeRows{data: db.test.Rows, scanErrorAt: db.test.ScanErrorAt, rowsError: db.test.RowsError}
	return db.rows, nil
}
type fakeRows struct { data [][]any; index int; current int; scanErrorAt *int; rowsError bool; closed bool }
func (rows *fakeRows) Next() bool {
	if rows.index >= len(rows.data) { return false }
	rows.current = rows.index
	rows.index++
	return true
}
func (rows *fakeRows) Scan(dest ...any) error {
	if rows.scanErrorAt != nil && rows.current == *rows.scanErrorAt { return scanFailure }
	record := rows.data[rows.current]
	if len(dest) != 3 || len(record) != 3 { return scanFailure }
	id, okID := record[0].(float64); name, okName := record[1].(string); total, okTotal := record[2].(float64)
	if !okID || !okName || !okTotal { return scanFailure }
	*(dest[0].(*int64)) = int64(id)
	*(dest[1].(*string)) = name
	*(dest[2].(*int64)) = int64(total)
	return nil
}
func (rows *fakeRows) Err() error { if rows.rowsError { return rowsFailure }; return nil }
func (rows *fakeRows) Close() error { rows.closed = true; return nil }
`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			db := &fakeDB{test: &current.Payload}
			values, err := LoadAccountTotals(context.Background(), db, current.Payload.Min)
			if current.Payload.Check == "query" {
				if db.query == expectedQuery { return "query:ok" }
				return "query:wrong"
			}
			if current.Payload.Check == "arg" {
				if db.query == expectedQuery && len(db.args) == 1 && db.args[0] == current.Payload.Min { return "arg:ok" }
				return "arg:wrong"
			}
			if current.Payload.Check == "close" {
				if db.rows != nil && db.rows.closed { return "close:ok" }
				return "close:missed"
			}
			if errors.Is(err, queryFailure) { return "error:query" }
			if errors.Is(err, scanFailure) { return "error:scan" }
			if errors.Is(err, rowsFailure) { return "error:rows" }
			if err != nil { return "error:other" }
			encoded, _ := json.Marshal(values)
			return string(encoded)
		})
	}`, tests)
}

func levelNine() Level {
	instructions := baseInstructions(
		"Build a deterministic JSON batch endpoint around bounded concurrent processing.",
		"func NewBatchHandler(processor Processor, workers int) http.Handler",
		"POST /batch with {\"values\":[string,...]}. Processor is supplied by the caller and may finish out of order or return an error.",
		"200 with {\"results\":[...]} in input order. Return JSON 400 for invalid input/workers, 405 for wrong methods, 404 for paths, and 422 when processing fails.",
		"Processor is supplied and the handler factory starts unfinished.",
	)
	instructions.Constraints = []string{"At most workers processor calls at once.", "Reject workers <= 0 without processing.", "Limit values to 100 items; each value must be non-blank after trimming.", "Preserve result order.", "Cancel processing on error and wait for workers.", "Reject unknown fields and trailing JSON."}
	instructions.Examples = []Example{{Input: `POST /batch {"values":["go","lang"]}, uppercase processor`, Output: `200 {"results":["GO","LANG"]}`}, {Input: `processor fails`, Output: `422 JSON error`}}
	instructions.Documentation = []DocumentationLink{{Label: "HTTP handlers", URL: "https://pkg.go.dev/net/http#Handler"}, {Label: "JSON Decoder", URL: "https://pkg.go.dev/encoding/json#Decoder"}, {Label: "Context cancellation", URL: "https://go.dev/blog/context"}}
	instructions.Hints = []string{"Validate the entire request before starting goroutines.", "Use indexed jobs and a fixed worker pool, as in level 7.", "Create a child context from r.Context and cancel it on the first processor error."}
	instructions.CommonPitfalls = []string{"Writing results in completion order", "Leaking workers after an error", "Starting work before all values are validated"}
	tests := makeIntegratedTests()
	return Level{ID: 9, Title: "Batch Gateway", Topic: "HTTP · parsing · concurrency", Difficulty: "Stretch+", Stretch: true, Signature: "NewBatchHandler(processor Processor, workers int) http.Handler", StarterCode: `import (
	"context"
	"encoding/json"
	"net/http"
)

type Processor func(context.Context, string) (string, error)

func NewBatchHandler(processor Processor, workers int) http.Handler {
	// TODO: validate POST /batch and process values with a bounded worker pool.
	_ = json.NewDecoder
	return http.NotFoundHandler()
}
`, Instructions: instructions, Tests: tests, build: buildLevelNine}
}

func makeIntegratedTests() []VisibleTest {
	cases := []struct {
		id, name, purpose, method, path, body, mode, expected string
		workers                                               int
	}{
		{"l9-basic", "Basic batch", "Processes every item.", "POST", "/batch", `{"values":["go","lang"]}`, "upper", `{"status":200,"body":{"results":["GO","LANG"]}}`, 2},
		{"l9-empty", "Empty batch", "Allows an empty list.", "POST", "/batch", `{"values":[]}`, "upper", `{"status":200,"body":{"results":[]}}`, 2},
		{"l9-order", "Out-of-order completion", "Preserves request order.", "POST", "/batch", `{"values":["slow","a","bb"]}`, "skew", `{"status":200,"body":{"results":["SLOW","A","BB"]}}`, 3},
		{"l9-one-worker", "Single worker", "Supports the minimum positive bound.", "POST", "/batch", `{"values":["a","b"]}`, "upper", `{"status":200,"body":{"results":["A","B"]}}`, 1},
		{"l9-bad-workers", "Invalid workers", "Rejects configuration without work.", "POST", "/batch", `{"values":["a"]}`, "upper", `{"status":400}`, 0},
		{"l9-method", "Wrong method", "Uses 405 for GET.", "GET", "/batch", "", "upper", `{"status":405,"allow":"POST"}`, 2},
		{"l9-path", "Wrong path", "Uses 404 outside /batch.", "POST", "/other", `{"values":[]}`, "upper", `{"status":404}`, 2},
		{"l9-json", "Malformed JSON", "Rejects broken input.", "POST", "/batch", `{"values":`, "upper", `{"status":400}`, 2},
		{"l9-unknown", "Unknown field", "Rejects fields outside the schema.", "POST", "/batch", `{"values":[],"fast":true}`, "upper", `{"status":400}`, 2},
		{"l9-blank", "Blank value", "Validates every item before work.", "POST", "/batch", `{"values":["ok","  "]}`, "upper", `{"status":400}`, 2},
		{"l9-null", "Null values", "Requires a JSON array.", "POST", "/batch", `{"values":null}`, "upper", `{"status":400}`, 2},
		{"l9-error", "Processor error", "Maps processing failures to 422.", "POST", "/batch", `{"values":["ok","fail","later"]}`, "error", `{"status":422}`, 2},
		{"l9-bound", "Concurrency bound", "Never exceeds worker count.", "POST", "/batch", `{"values":["a","b","c","d","e"]}`, "bound", `{"status":200,"body":{"results":["A","B","C","D","E"]},"bound":"ok"}`, 2},
	}
	out := make([]VisibleTest, 0, len(cases))
	for _, item := range cases {
		payload := map[string]any{"method": item.method, "path": item.path, "body": item.body, "mode": item.mode, "workers": item.workers}
		out = append(out, test(item.id, item.name, item.purpose, fmt.Sprintf("%s %s %s", item.method, item.path, item.body), item.expected, payload))
	}
	return out
}

func buildLevelNine(tests []VisibleTest) string {
	return harness(commonImports("context", "errors", "net/http/httptest", "strings", "sync/atomic"), `var processFailure = errors.New("processing failed")
type inputCase struct {
	Method string `+"`json:\"method\"`"+`
	Path string `+"`json:\"path\"`"+`
	Body string `+"`json:\"body\"`"+`
	Mode string `+"`json:\"mode\"`"+`
	Workers int `+"`json:\"workers\"`"+`
}
type wireCase struct {
	ID string `+"`json:\"id\"`"+`
	Payload inputCase `+"`json:\"payload\"`"+`
}
type responseView struct {
	Status int `+"`json:\"status\"`"+`
	Allow string `+"`json:\"allow,omitempty\"`"+`
	Body map[string]any `+"`json:\"body,omitempty\"`"+`
	Bound string `+"`json:\"bound,omitempty\"`"+`
}`, `	var tests []wireCase
	_ = json.Unmarshal([]byte(raw), &tests)
	for _, current := range tests {
		current := current
		assessmentRun(current.ID, func() string {
			var active int32
			var exceeded int32
			processor := func(ctx context.Context, value string) (string, error) {
				now := atomic.AddInt32(&active, 1)
				if now > int32(current.Payload.Workers) && current.Payload.Workers > 0 { atomic.StoreInt32(&exceeded, 1) }
				defer atomic.AddInt32(&active, -1)
				if current.Payload.Mode == "error" && value == "fail" { return "", processFailure }
				if current.Payload.Mode == "skew" { for spin := 0; spin < (8-len(value))*500; spin++ {} }
				return strings.ToUpper(value), nil
			}
			handler := NewBatchHandler(processor, current.Payload.Workers)
			request := httptest.NewRequest(current.Payload.Method, current.Payload.Path, strings.NewReader(current.Payload.Body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			view := responseView{Status: recorder.Code, Allow: recorder.Header().Get("Allow")}
			if recorder.Code == 200 { _ = json.Unmarshal(recorder.Body.Bytes(), &view.Body) }
			if current.Payload.Mode == "bound" {
				if atomic.LoadInt32(&exceeded) == 0 { view.Bound = "ok" } else { view.Bound = "exceeded" }
			}
			encoded, _ := json.Marshal(view)
			return string(encoded)
		})
	}`, tests)
}

func Validate() error {
	levels := Levels()
	if len(levels) != 9 {
		return fmt.Errorf("expected 9 levels, got %d", len(levels))
	}
	for index, level := range levels {
		if level.ID != index+1 || len(level.Tests) < 10 || len(level.Tests) > 18 {
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
