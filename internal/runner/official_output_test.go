package runner

import (
	"context"
	"strings"
	"testing"
)

func TestDecodeOfficialOutcomeMapsPassesAndFailures(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	human := "Exercise: validate-stack\nRESULT: FAIL\n\n  PASS tp4/distinct_ok [accepted]\n  FAIL tp4/duplicate [rejected] — wrong exit\n        Input (7 bytes; trailing newline: no):\n        | 1 2 2 3\n        Expected stdout (0 bytes):\n        | (empty)\n\nSummary: 1 passed, 1 failed"
	raw := human + "\n{\"Ok\":false,\"Output\":\"details\"}\n"
	outcome, err := decodeOfficialOutcome(level, raw, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(outcome.results) != 2 || outcome.results[0].Actual != "pass" || outcome.results[1].Actual != "fail" {
		t.Fatalf("unexpected mapped results %#v", outcome.results)
	}
	if outcome.results[1].Input != "1 2 2 3" || strings.Contains(outcome.stdout, "Expected stdout") {
		t.Fatalf("official output was not compacted: %#v", outcome)
	}
	if outcome.stdout != "Exercise: validate-stack\nRESULT: FAIL\nPASS tp4/distinct_ok\nFAIL tp4/duplicate\nSummary: 1 passed, 1 failed" {
		t.Fatalf("unexpected compact output %q", outcome.stdout)
	}
}

func TestDecodeOfficialOutcomeMapsAggregateFailure(t *testing.T) {
	level := mustExercise(t, "checkpoint/push-swap")
	raw := "Exercise: push-swap\nRESULT: FAIL\n  FAIL al3/case-93 — too many operations\n{\"Ok\":false,\"Output\":\"details\"}\n"
	outcome, err := decodeOfficialOutcome(level, raw, "")
	if err != nil || len(outcome.results) != 1 || outcome.results[0].Actual != "fail" {
		t.Fatalf("aggregate mapping = %#v, %v", outcome, err)
	}
}

func TestDecodeOfficialOutcomeRejectsFailedAggregate(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	raw := "Exercise: validate-stack\nRESULT: PASS\n  PASS tp4/distinct_ok [accepted]\n{\"Ok\":false,\"Output\":\"details\"}\n"
	outcome, err := decodeOfficialOutcome(level, raw, "")
	if err != nil {
		t.Fatal(err)
	}
	if outcome.status == executionSuccess || outcome.runtimeError == "" {
		t.Fatalf("failed aggregate was accepted: %#v", outcome)
	}
}

func TestFailedAggregateCannotIssueReceipt(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	raw := "Exercise: validate-stack\nRESULT: PASS\n  PASS tp4/distinct_ok [accepted]\n{\"Ok\":false,\"Output\":\"details\"}\n"
	outcome, err := decodeOfficialOutcome(level, raw, "")
	if err != nil {
		t.Fatal(err)
	}
	outcome.results = make([]wireResult, 0, len(level.Tests))
	for _, test := range level.Tests {
		outcome.results = append(outcome.results, wireResult{ID: test.ID, Actual: "pass"})
	}
	result := newEngine(&outcomeAdapter{outcome: outcome}, 1, testIssuer{}).Run(context.Background(), level, level.StarterCode)
	if result.Passed || result.Receipt != "" {
		t.Fatalf("failed aggregate completed as pass: %#v", result)
	}
}

func TestDockerRunUsesOnlyExactSourceMountAndPinnedImage(t *testing.T) {
	args := dockerRunArgs("imperative-go-assessment-0123456789abcdef01234567", "checkpoint/validate-stack", `C:\\tmp\\main.go`, nil)
	joined := strings.Join(args, " ")
	for _, required := range []string{dockerImage, "--network none", "--read-only", "EMIT_JSON=1", "target=/jail/student/validate-stack/main.go,readonly"} {
		if !strings.Contains(joined, required) {
			t.Errorf("missing %q in %s", required, joined)
		}
	}
	if strings.Contains(joined, "docker.sock") || strings.Contains(joined, "target=/workspace") {
		t.Fatalf("unsafe broad mount in %s", joined)
	}
}

func TestDockerRunMountsDeclaredResourcesReadOnly(t *testing.T) {
	args := dockerRunArgs(
		"imperative-go-assessment-0123456789abcdef01234567",
		"checkpoint/ascii-render",
		"C:\\tmp\\main.go",
		[]dockerResourceMount{
			{sourcePath: "C:\\tmp\\banner.txt", targetPath: "/jail/student/ascii-render/resources/banner.txt"},
		},
	)
	var mounts []string
	for index := 0; index < len(args); index++ {
		if args[index] == "--mount" {
			mounts = append(mounts, args[index+1])
		}
	}
	want := []string{
		"type=bind,source=C:\\tmp\\main.go,target=/jail/student/ascii-render/main.go,readonly",
		"type=bind,source=C:\\tmp\\banner.txt,target=/jail/student/ascii-render/resources/banner.txt,readonly",
	}
	if len(mounts) != len(want) {
		t.Fatalf("got %d mounts, want %d: %#v", len(mounts), len(want), mounts)
	}
	for index := range want {
		if mounts[index] != want[index] {
			t.Fatalf("mount %d = %q, want %q", index, mounts[index], want[index])
		}
	}
}
