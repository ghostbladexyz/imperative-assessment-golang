package assessment

import (
	"strings"
	"testing"
)

func TestSourcePolicyAcceptsDocumentedBuiltinsAndPackages(t *testing.T) {
	t.Parallel()
	level := mustFindExercise(t, "zone01/50")
	source := `package main

import "fmt"

func ConcatSlice(left, right []int) []int {
	fmt.Println("joining")
	values := []int{}
	values = append(values, left...)
	values = append(values, right...)
	return values
}`
	if err := ValidateSourcePolicy(level, source); err != nil {
		t.Fatalf("documented source was rejected: %v", err)
	}
}

func TestSourcePolicyRejectsUnlistedBuiltin(t *testing.T) {
	t.Parallel()
	level := mustFindExercise(t, "zone01/50")
	source := `package main
func ConcatSlice(left, right []int) []int {
	values := []int{}
	copy(values, values)
	return values
}`
	err := ValidateSourcePolicy(level, source)
	if err == nil || !strings.Contains(err.Error(), `built-in "copy" is not allowed`) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}

func TestSourcePolicyRejectsAliasedUnlistedBuiltin(t *testing.T) {
	t.Parallel()
	level := mustFindExercise(t, "zone01/50")
	source := `package main
func ConcatSlice(left, right []int) []int {
	forbidden := copy
	values := []int{}
	forbidden(values, values)
	return values
}`
	err := ValidateSourcePolicy(level, source)
	if err == nil || !strings.Contains(err.Error(), `built-in "copy" is not allowed`) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}

func TestSourcePolicyRejectsUnlistedPackage(t *testing.T) {
	t.Parallel()
	level := mustFindExercise(t, "zone01/50")
	source := `package main
import "regexp"
func ConcatSlice(left, right []int) []int {
	_ = regexp.MustCompile(".")
	return left
}`
	err := ValidateSourcePolicy(level, source)
	if err == nil || !strings.Contains(err.Error(), `package "regexp" is not allowed`) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}

func TestPrintedOutputPolicyRequiresZ01(t *testing.T) {
	t.Parallel()
	level := mustFindExercise(t, "piscine/1014")
	z01Source := `package main
import "github.com/01-edu/z01"
func PrintComb() { z01.PrintRune('0') }`
	if err := ValidateSourcePolicy(level, z01Source); err != nil {
		t.Fatalf("z01 output was rejected: %v", err)
	}

	fmtSource := `package main
import "fmt"
func PrintComb() { fmt.Print("0") }`
	err := ValidateSourcePolicy(level, fmtSource)
	if err == nil || !strings.Contains(err.Error(), `package "fmt" is not allowed`) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}

func mustFindExercise(t *testing.T, key ExerciseKey) Level {
	t.Helper()
	level, found := FindExercise(key)
	if !found {
		t.Fatalf("missing exercise %q", key)
	}
	return level
}
