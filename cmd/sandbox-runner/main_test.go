package main

import (
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
