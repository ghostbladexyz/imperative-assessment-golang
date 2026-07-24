package assessment

import (
	"strings"
	"testing"
)

func TestSourcePolicyAcceptsDocumentedBuiltinsAndPackages(t *testing.T) {
	t.Parallel()
	level, _ := FindLevel(22)
	source := `package main

import "strings"

func helper(value string) string { return strings.ToLower(value) }
func NormalizeTokens(input string) []string {
	values := make([]string, 0)
	if len(input) > 0 {
		values = append(values, helper(input))
		println(input)
	}
	return values
}`
	if err := ValidateSourcePolicy(level, source); err != nil {
		t.Fatalf("documented source was rejected: %v", err)
	}
}

func TestSourcePolicyRejectsUnlistedBuiltin(t *testing.T) {
	t.Parallel()
	level, _ := FindLevel(22)
	source := `package main
func NormalizeTokens(input string) []string {
	values := []string{input}
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
	level, _ := FindLevel(22)
	source := `package main
func NormalizeTokens(input string) []string {
	forbidden := copy
	values := []string{input}
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
	level, _ := FindLevel(22)
	source := `package main
import "regexp"
func NormalizeTokens(input string) []string {
	return []string{regexp.MustCompile(".").FindString(input)}
}`
	err := ValidateSourcePolicy(level, source)
	if err == nil || !strings.Contains(err.Error(), `package "regexp" is not allowed`) {
		t.Fatalf("unexpected policy error: %v", err)
	}
}
