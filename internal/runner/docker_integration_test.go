package runner

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

func TestDockerRunnerIntegration(t *testing.T) {
	if os.Getenv("IMPERATIVE_DOCKER_INTEGRATION") != "1" {
		t.Skip("set IMPERATIVE_DOCKER_INTEGRATION=1 to run tests against a real Docker daemon")
	}
	commands := processExecutor{}
	checkCtx, cancelCheck := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCheck()
	if err := CheckDocker(checkCtx, commands, "docker"); err != nil {
		t.Skipf("real Docker integration skipped: %v", err)
	}
	buildCtx, cancelBuild := context.WithTimeout(context.Background(), 10*time.Minute)
	sandbox, err := NewDocker(buildCtx, DockerOptions{MaxConcurrent: 1, Commands: commands})
	cancelBuild()
	if err != nil {
		t.Fatal(err)
	}
	level := mustExercise(t, "checkpoint/validate-stack")

	t.Run("correct solution passes", func(t *testing.T) {
		source := strings.ReplaceAll(level.StarterCode, "\r\n", "\n")
		source = strings.Replace(source, "\t_ = strconv.Atoi\n\t// TODO: implement valid.\n\treturn true", "\tseen := map[int]bool{}\n\tfor _, raw := range args { value, err := strconv.Atoi(raw); if err != nil || seen[value] { return false }; seen[value] = true }\n\treturn true", 1)
		if source == strings.ReplaceAll(level.StarterCode, "\r\n", "\n") {
			t.Fatal("correct-solution fixture did not replace the TODO")
		}
		result := sandbox.Run(context.Background(), level, source)
		if !result.Passed {
			t.Fatalf("correct solution failed: %#v", result)
		}
	})

	t.Run("starter fails assertions", func(t *testing.T) {
		result := sandbox.Run(context.Background(), level, level.StarterCode)
		if result.Passed {
			t.Fatalf("starter unexpectedly passed: %#v", result)
		}
		if strings.Contains(result.Stdout, "Expected stdout") || strings.Count(result.Stdout, "\n") > len(level.Tests)+5 {
			t.Fatalf("official output was not compacted:\n%s", result.Stdout)
		}
		if result.Results[1].Input != "1 2 2 3" {
			t.Fatalf("concrete failed input was not captured: %#v", result.Results[1])
		}
	})

	t.Run("invalid Go returns compiler error", func(t *testing.T) {
		result := sandbox.Run(context.Background(), level, "package main\nfunc valid([]string) bool { definitely not go }")
		if result.CompileError == "" {
			t.Fatalf("missing compiler error: %#v", result)
		}
	})

	t.Run("every starter reports an official outcome", func(t *testing.T) {
		for _, current := range assessment.Levels() {
			current := current
			t.Run(current.Title, func(t *testing.T) {
				result := sandbox.Run(context.Background(), current, current.StarterCode)
				if len(current.Tests) == 1 && result.CompileError == "" && result.FailureKind == "" &&
					result.RuntimeError == "The official grader did not return a result for the check." && len(result.Results) == 1 &&
					result.Results[0].Status == "runtime" && result.Results[0].Actual == "" {
					return
				}
				if result.FailureKind != "" || result.RuntimeError != "" || result.CompileError != "" {
					t.Fatalf("official grader integration failed: %#v", result)
				}
				recognized := 0
				for _, check := range result.Results {
					if check.Actual != "" {
						recognized++
					}
				}
				if recognized == 0 {
					t.Fatalf("no official checks were recognized: %#v", result)
				}
			})
		}
	})
}
