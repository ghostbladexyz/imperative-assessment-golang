package assessment

import (
	"embed"
	"fmt"
	"strings"
)

const checkpointUpstreamCommit = "9ef45ee4a164b6d9b8ff3c5ad9b97753c3d63297"

// checkpointAssets mirror the pinned upstream commit, with README.md flowchart
// fences annotated for the embedded frontend.
//
//go:embed checkpoint
var checkpointAssets embed.FS

type checkpointSpec struct {
	slug       string
	topic      string
	difficulty string
	signature  string
}

func checkpointLevels() []Level {
	specs := []checkpointSpec{
		{"validate-stack", "Parsing · maps", "Beginner", "valid(args []string) bool"},
		{"safe-sum", "Robust parsing · errors", "Beginner", "sumLine(line string) string"},
		{"tetris", "Grids · coordinates", "Beginner", "canPlace(...); stamp(...)"},
		{"method-routing", "HTTP · routing", "Easy", "rootHandler(...); submitHandler(...)"},
		{"status-matrix", "HTTP · status handling", "Easy", "rootHandler(...); generateHandler(...)"},
		{"wordcount", "Maps · deterministic sorting", "Easy", "orderedWords(counts map[string]int) []string"},
		{"ls-format", "Formatting · time", "Intermediate", "permString(...); fmtTime(...)"},
		{"ascii-render", "Parsing · ASCII art", "Intermediate", "render(input string, banner []string) string"},
		{"reloaded-rev", "Token streams · transformations", "Intermediate", "resolve(toks []Token) []Token"},
		{"reloaded-format", "Token streams · formatting", "Intermediate", "resolve(toks []Token) []Token"},
		{"lemin-path", "Graphs · breadth-first search", "Hard", "shortestPath(c *Colony) []string"},
		{"lemin-why", "Algorithms · scheduling", "Hard", "minTurns(paths []int, ants int) int"},
		{"bounded-fanout", "Concurrency · worker pools", "Hard", "squareAll(nums []int) []int"},
		{"consume-join", "HTTP clients · joins", "Hard", "run(base, target string) (string, error)"},
		{"push-swap", "Algorithms · constrained operations", "Expert", "sort()"},
		{"broadcast", "Concurrency · TCP", "Expert", "broadcast(...); handle(...)"},
		{"reactions", "SQL · transactions · concurrency", "Expert", "setReaction(...); count(...)"},
	}

	levels := make([]Level, 0, len(specs))
	for _, spec := range specs {
		starter, starterErr := checkpointAssets.ReadFile("checkpoint/" + spec.slug + "/starter.go.txt")
		subject, subjectErr := checkpointAssets.ReadFile("checkpoint/" + spec.slug + "/README.md")
		var definitionErr error
		if starterErr != nil || subjectErr != nil {
			definitionErr = fmt.Errorf("load upstream assets: starter=%v subject=%v", starterErr, subjectErr)
		} else if !strings.Contains(string(starter), "package main") || !strings.HasPrefix(string(subject), "# "+spec.slug) {
			definitionErr = fmt.Errorf("upstream assets do not match %q", spec.slug)
		}
		levels = append(levels, Level{
			Key: ExerciseKey("checkpoint/" + spec.slug), Title: spec.slug, Topic: spec.topic,
			Difficulty: spec.difficulty, Signature: spec.signature, StarterCode: string(starter),
			Subject: string(subject), Resources: checkpointResources(spec.slug), Tests: officialChecks(spec.slug),
			definitionErr: definitionErr, unrestricted: true,
			Instructions: Instructions{AllowedBuiltins: []string{}, AllowedPackages: []string{}},
		})
	}
	return levels
}

func checkpointResources(slug string) []ExerciseResource {
	names := map[string][]string{"ascii-render": {"banner.txt"}, "reactions": {"schema.sql"}}[slug]
	resources := make([]ExerciseResource, 0, len(names))
	for _, name := range names {
		content, err := checkpointAssets.ReadFile("checkpoint/" + slug + "/resources/" + name)
		if err == nil {
			resources = append(resources, ExerciseResource{Name: name, Content: string(content)})
		}
	}
	return resources
}
