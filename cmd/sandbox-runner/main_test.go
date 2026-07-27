package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadRequestRejectsMalformedUnknownAndTrailingJSON(t *testing.T) {
	for _, input := range []string{
		`{`,
		`{"source":"package main","harness":"package main","tests":[{"id":"x","expected":""}],"compileTimeoutMs":1000,"runtimeTimeoutMs":1000,"outputLimitBytes":1024,"unknown":true}`,
		`{"source":"package main","harness":"package main","tests":[{"id":"x","expected":""}],"compileTimeoutMs":1000,"runtimeTimeoutMs":1000,"outputLimitBytes":1024} {}`,
	} {
		if _, err := readRequest(strings.NewReader(input)); err == nil {
			t.Fatalf("expected request error for %q", input)
		}
	}
}

func TestResultProtocolSurvivesUnterminatedStudentOutput(t *testing.T) {
	stdout := "Hello__IMPERATIVE_ASSESSMENT_RESULT__{\"id\":\"l1-empty\",\"actual\":\"[]\",\"durationMs\":0.035}\n" +
		"Again__IMPERATIVE_ASSESSMENT_RESULT__{\"id\":\"l1-word\",\"actual\":\"[]\",\"durationMs\":0.001}\n"

	results := parseResults(stdout)
	if len(results) != 2 || results[0].ID != "l1-empty" || results[1].ID != "l1-word" {
		t.Fatalf("protocol results were not recovered: %#v", results)
	}
	if results[0].Stdout != "Hello" || results[1].Stdout != "Again" {
		t.Fatalf("student output was not associated with its test: %#v", results)
	}
	if output := stripMarkers(stdout); output != "Hello\nAgain" {
		t.Fatalf("internal protocol leaked into stdout: %q", output)
	}
}

func TestSubmissionModuleProvidesZ01Offline(t *testing.T) {
	workspace := t.TempDir()
	if err := writeSubmissionModule(workspace); err != nil {
		t.Fatal(err)
	}
	module, err := os.ReadFile(filepath.Join(workspace, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(module)
	if !strings.Contains(content, "github.com/01-edu/z01 v0.1.0") ||
		!strings.Contains(content, "replace github.com/01-edu/z01 => /opt/z01") {
		t.Errorf("go.mod does not pin the bundled z01 module: %s", content)
	}
}
