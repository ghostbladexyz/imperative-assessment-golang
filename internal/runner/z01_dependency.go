//go:build dependencies

package runner

// Keep the learner-facing package in go.mod even though only submitted code imports it.
import _ "github.com/01-edu/z01"
