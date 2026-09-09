package runner

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
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

func TestDecodeOfficialOutcomeRejectsMissingSingleAggregateResult(t *testing.T) {
	level := mustExercise(t, "checkpoint/push-swap")
	for _, envelope := range []string{"false", "true"} {
		t.Run("Ok="+envelope, func(t *testing.T) {
			raw := "Exercise: push-swap\nRESULT: FAIL\n{\"Ok\":" + envelope + ",\"Output\":\"details\"}\n"
			outcome, err := decodeOfficialOutcome(level, raw, "")
			if err == nil {
				t.Fatal("missing aggregate result was accepted")
			}
			if len(outcome.results) != 0 {
				t.Fatalf("missing aggregate result fabricated checks: %#v", outcome.results)
			}
		})
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

func TestPartialOfficialOutcomeNeedsFailureEvidenceForNotRunChecks(t *testing.T) {
	level := mustExercise(t, "checkpoint/validate-stack")
	for _, test := range []struct {
		name             string
		raw              string
		wantStoppedAfter bool
		wantMissing      string
	}{
		{
			name:             "passing envelope leaves missing checks in error",
			raw:              "Exercise: validate-stack\nRESULT: PASS\n  PASS tp4/distinct_ok [accepted]\n{\"Ok\":true,\"Output\":\"details\"}\n",
			wantStoppedAfter: false,
			wantMissing:      "runtime",
		},
		{
			name:             "failed envelope marks later checks not run",
			raw:              "Exercise: validate-stack\nRESULT: FAIL\n  FAIL tp4/duplicate — wrong exit\n{\"Ok\":false,\"Output\":\"details\"}\n",
			wantStoppedAfter: true,
			wantMissing:      "not_run",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			outcome, err := decodeOfficialOutcome(level, test.raw, "")
			if err != nil {
				t.Fatal(err)
			}
			if outcome.stoppedAfterFailure != test.wantStoppedAfter {
				t.Fatalf("stoppedAfterFailure = %v, want %v", outcome.stoppedAfterFailure, test.wantStoppedAfter)
			}
			result := newEngine(&outcomeAdapter{outcome: outcome}, 1, nil).Run(context.Background(), level, level.StarterCode)
			if result.Results[2].Status != test.wantMissing {
				t.Fatalf("missing check status = %q, want %q: %#v", result.Results[2].Status, test.wantMissing, result.Results)
			}
		})
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

func TestDockerRunUsesPinnedAMD64Platform(t *testing.T) {
	args := dockerRunArgs("imperative-go-assessment-0123456789abcdef01234567", "checkpoint/validate-stack", `C:\tmp\main.go`, nil)
	for index, argument := range args {
		if argument != "--platform" {
			continue
		}
		if index+1 >= len(args) {
			t.Fatal("Docker run platform option has no value")
		}
		if args[index+1] != dockerPlatform {
			t.Fatalf("Docker run platform = %q, want %q", args[index+1], dockerPlatform)
		}
		return
	}
	t.Fatal("Docker run did not specify a platform")
}

func TestDockerRunUsesWritableGoCache(t *testing.T) {
	args := dockerRunArgs("imperative-go-assessment-0123456789abcdef01234567", "checkpoint/validate-stack", `C:\tmp\main.go`, nil)
	for index, argument := range args {
		if argument != "--env" {
			continue
		}
		if index+1 < len(args) && args[index+1] == "GOCACHE=/tmp/go-build" {
			return
		}
	}
	t.Fatal("Docker run did not override GOCACHE with a writable tmpfs path")
}

func TestStageOfficialResourcesAreReadableByLearner(t *testing.T) {
	mounts, err := stageOfficialResources(t.TempDir(), "checkpoint/ascii-render", []assessment.ExerciseResource{{
		Name: "banner.txt", Content: "banner",
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(mounts) != 1 {
		t.Fatalf("staged %d resources, want 1", len(mounts))
	}
	info, err := os.Stat(mounts[0].sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o044 != 0o044 {
		t.Fatalf("resource permissions = %o, want group and other read access", info.Mode().Perm())
	}
}

func TestOfficialStartupMessageExplainsMissingAMD64Emulation(t *testing.T) {
	want := "The official grader requires Docker linux/amd64 emulation. Enable amd64 emulation in Docker, then rerun the assessment."
	got := officialStartupMessage("", "standard_init_linux.go: exec format error")
	if got != want {
		t.Fatalf("startup message = %q, want %q", got, want)
	}
}
