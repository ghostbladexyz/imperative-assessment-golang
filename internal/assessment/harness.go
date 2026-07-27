package assessment

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func harness(imports []string, declarations, loop string, tests []VisibleTest) string {
	payloads := make([]any, 0, len(tests))
	for _, test := range tests {
		payloads = append(payloads, map[string]any{
			"id":       test.ID,
			"expected": test.Expected,
			"payload":  test.payload,
		})
	}
	data, _ := json.Marshal(payloads)
	var quotedImports strings.Builder
	for _, item := range imports {
		fmt.Fprintf(&quotedImports, "\t%q\n", item)
	}
	return fmt.Sprintf(`package main

import (
%s)

const assessmentPrefix = "__IMPERATIVE_ASSESSMENT_RESULT__"

type assessmentResult struct {
	ID         string  `+"`json:\"id\"`"+`
	Actual     string  `+"`json:\"actual\"`"+`
	Failure    string  `+"`json:\"failure,omitempty\"`"+`
	DurationMS float64 `+"`json:\"durationMs\"`"+`
}

func assessmentRun(id string, fn func() string) {
	started := time.Now()
	result := assessmentResult{ID: id}
	func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				result.Failure = fmt.Sprintf("panic: %%v", recovered)
			}
		}()
		result.Actual = fn()
	}()
	result.DurationMS = float64(time.Since(started).Microseconds()) / 1000
	encoded, _ := json.Marshal(result)
	fmt.Println(assessmentPrefix + string(encoded))
}

func assessmentRunPrinted(id string, fn func()) {
	started := time.Now()
	result := assessmentResult{ID: id}
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		result.Failure = fmt.Sprintf("capture output: %%v", err)
	} else {
		os.Stdout = writer
		var captured bytes.Buffer
		copied := make(chan error, 1)
		go func() {
			_, copyErr := io.Copy(&captured, reader)
			copied <- copyErr
		}()
		func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					result.Failure = fmt.Sprintf("panic: %%v", recovered)
				}
			}()
			fn()
		}()
		_ = writer.Close()
		os.Stdout = original
		if copyErr := <-copied; copyErr != nil && result.Failure == "" {
			result.Failure = fmt.Sprintf("capture output: %%v", copyErr)
		}
		_ = reader.Close()
		fmt.Print(captured.String())
		encodedActual, _ := json.Marshal(captured.String())
		result.Actual = string(encodedActual)
	}
	result.DurationMS = float64(time.Since(started).Microseconds()) / 1000
	encoded, _ := json.Marshal(result)
	fmt.Println(assessmentPrefix + string(encoded))
}

%s

func main() {
	raw := %s
%s
}
`, quotedImports.String(), declarations, strconv.Quote(string(data)), loop)
}

func commonImports(extra ...string) []string {
	return append([]string{"bytes", "encoding/json", "fmt", "io", "os", "time"}, extra...)
}

func jsonString(value any) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}
