package assessment

import (
	"strings"
	"testing"
)

func TestSourcePolicyAcceptsDocumentedBuiltinsAndPackages(t *testing.T) {
	t.Parallel()
	level := findLevelByTitle(t, "Concat Slice")
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
	level := findLevelByTitle(t, "Concat Slice")
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
	level := findLevelByTitle(t, "Concat Slice")
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
	level := findLevelByTitle(t, "Concat Slice")
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

func findLevelByTitle(t *testing.T, title string) Level {
	t.Helper()
	for _, level := range Levels() {
		if level.Title == title {
			return level
		}
	}
	t.Fatalf("missing level %q", title)
	return Level{}
}
