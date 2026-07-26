package runner

import "testing"

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
