package runner

import (
	"context"
	"strings"
	"testing"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

type testIssuer struct{}

func (testIssuer) Issue(levelID int, sourceHash string) (string, error) {
	return "verified-receipt", nil
}

func TestRunnerAcceptsCorrectAndRejectsIncorrectCode(t *testing.T) {
	level, _ := assessment.FindLevel(1)
	localRunner := New("go", 1, testIssuer{})
	correct := `func NormalizeTokens(input string) []string {
	result := make([]string, 0)
	current := make([]rune, 0)
	flush := func() {
		if len(current) > 0 {
			result = append(result, string(current))
			current = current[:0]
		}
	}
	for _, r := range input {
		if r >= 'A' && r <= 'Z' {
			r += 'a' - 'A'
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			current = append(current, r)
		} else {
			flush()
		}
	}
	flush()
	return result
}`
	result := localRunner.Run(context.Background(), level, correct, nil)
	if !result.Passed || result.PassedCount != len(level.Tests) || result.Receipt == "" {
		t.Fatalf("correct solution failed: %#v", result)
	}

	incorrect := `func NormalizeTokens(input string) []string { return []string{input} }`
	result = localRunner.Run(context.Background(), level, incorrect, nil)
	if result.Passed || result.PassedCount == len(level.Tests) {
		t.Fatal("incorrect solution passed")
	}
}

func TestRunnerReportsCompilerErrors(t *testing.T) {
	level, _ := assessment.FindLevel(1)
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`func NormalizeTokens(input string) []string { definitely not go }`,
		nil,
	)
	if result.CompileError == "" {
		t.Fatalf("expected compiler error, got %#v", result)
	}
}

func TestRunnerTerminatesTimeout(t *testing.T) {
	level, _ := assessment.FindLevel(1)
	result := New("go", 1, nil).Run(
		context.Background(),
		level,
		`func NormalizeTokens(input string) []string { for {} }`,
		[]string{"l1-word"},
	)
	if !result.TimedOut {
		t.Fatalf("expected timeout, got %#v", result)
	}
}

func TestEveryStarterCompilesWithItsHarness(t *testing.T) {
	localRunner := New("go", 1, nil)
	for _, level := range assessment.Levels() {
		result := localRunner.Run(
			context.Background(),
			level,
			level.StarterCode,
			[]string{level.Tests[0].ID},
		)
		if result.CompileError != "" {
			t.Fatalf("level %d starter did not compile: %s", level.ID, result.CompileError)
		}
	}
}

func TestLocalRunnerEnforcesOutputLimit(t *testing.T) {
	level, _ := assessment.FindLevel(1)
	result := NewLocal("go", 1, nil).Run(
		context.Background(),
		level,
		`func NormalizeTokens(input string) []string {
			for i := 0; i < 300000; i++ { print("output") }
			return []string{input}
		}`,
		[]string{"l1-word"},
	)
	if result.FailureKind != FailureOutput {
		t.Fatalf("expected output failure, got %#v", result)
	}
	if len(result.Stderr) > MaxOutputBytes {
		t.Fatalf("stderr was not capped: %d bytes", len(result.Stderr))
	}
}

func TestFormatSourceRestoresAndRemovesPackageDeclaration(t *testing.T) {
	formatted, err := FormatSource("\ufeff package main\n\nfunc solve(){ }")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(formatted, "package main") || formatted != "func solve() {}\n" {
		t.Fatalf("unexpected formatted source %q", formatted)
	}
}
