package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
	"github.com/pleft/imperative-assessment-golang/internal/receipts"
	"github.com/pleft/imperative-assessment-golang/internal/runner"
)

type fakeRunner struct {
	info runner.Info
	runs int
}

func (fake *fakeRunner) Run(_ context.Context, level assessment.Level, _ string, _ []string) runner.RunResult {
	fake.runs++
	return runner.RunResult{ExerciseKey: level.Key, LevelID: level.ID}
}

func (fake *fakeRunner) Format(_ context.Context, source string) (string, error) {
	return source, nil
}

func (fake *fakeRunner) Info() runner.Info {
	return fake.info
}

func TestRunnerSelectionDefaultsToDocker(t *testing.T) {
	if defaultRunnerMode != "docker" {
		t.Fatalf("default mode is %q", defaultRunnerMode)
	}
	for input, expected := range map[string]runner.Mode{
		"docker": runner.ModeDocker,
		"DOCKER": runner.ModeDocker,
		"local":  runner.ModeLocal,
	} {
		actual, err := parseRunnerMode(input)
		if err != nil || actual != expected {
			t.Fatalf("parseRunnerMode(%q) = %q, %v", input, actual, err)
		}
	}
	if _, err := parseRunnerMode("browser"); err == nil {
		t.Fatal("expected invalid runner selection error")
	}
}

func TestHealthReportsSafeRunnerConfiguration(t *testing.T) {
	execution := &fakeRunner{info: runner.Info{
		Mode:         runner.ModeDocker,
		SandboxReady: true,
		GoVersion:    "go1.23.12",
		DockerImage:  "imperative-go-assessment-runner:abc123",
		Message:      "Docker sandbox ready.",
	}}
	handler, err := routes(&api{runner: execution})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["runnerMode"] != "docker" || body["sandboxReady"] != true ||
		body["goVersion"] != "go1.23.12" ||
		body["dockerImage"] != "imperative-go-assessment-runner:abc123" {
		t.Fatalf("unexpected health response %#v", body)
	}
	for _, sensitive := range []string{`C:\`, "/Users/", "npipe:", "docker.sock"} {
		if strings.Contains(response.Body.String(), sensitive) {
			t.Fatalf("health response exposed sensitive value %q", sensitive)
		}
	}

	localHandler, err := routes(&api{runner: &fakeRunner{info: runner.Info{
		Mode: runner.ModeLocal, GoVersion: "go1.25.4", Message: "Local runner selected.",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	localResponse := httptest.NewRecorder()
	localHandler.ServeHTTP(localResponse, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if strings.Contains(localResponse.Body.String(), "dockerImage") {
		t.Fatalf("local health response should omit Docker image: %s", localResponse.Body.String())
	}
}

func TestBrowserCannotSelectRunnerMode(t *testing.T) {
	execution := &fakeRunner{info: runner.Info{Mode: runner.ModeDocker}}
	handler, err := routes(&api{runner: execution})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/run",
		strings.NewReader(`{"exerciseKey":"foundation/1","code":"func NormalizeTokens(string) []string { return nil }","testIds":[],"runner":"local"}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}
	if execution.runs != 0 {
		t.Fatal("runner was invoked despite a browser-controlled mode field")
	}
}

func TestCataloguePublishesProgressIdentityManifest(t *testing.T) {
	handler, err := routes(&api{runner: &fakeRunner{}})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/levels", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}
	var body struct {
		Levels                []assessment.Level `json:"levels"`
		ProgressSchemaVersion int                `json:"progressSchemaVersion"`
		LegacyProgress        struct {
			SchemaVersion int                      `json:"schemaVersion"`
			ExerciseKeys  []assessment.ExerciseKey `json:"exerciseKeys"`
		} `json:"legacyProgress"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.ProgressSchemaVersion != progressSchemaVersion ||
		body.LegacyProgress.SchemaVersion != 4 ||
		len(body.LegacyProgress.ExerciseKeys) != len(body.Levels) ||
		body.Levels[0].Key != "foundation/1" ||
		body.LegacyProgress.ExerciseKeys[0] != body.Levels[0].Key {
		t.Fatalf("unexpected catalogue manifest: %#v", body)
	}
}

func TestReceiptValidationReturnsExerciseKeys(t *testing.T) {
	manager, err := receipts.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := manager.Issue("foundation/1", "source")
	if err != nil {
		t.Fatal(err)
	}
	handler, err := routes(&api{runner: &fakeRunner{}, receipts: manager})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/receipts/validate",
		strings.NewReader(`{"receipts":{"foundation/1":"`+encoded+`"}}`),
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK ||
		!strings.Contains(response.Body.String(), `"validExerciseKeys":["foundation/1"]`) {
		t.Fatalf("unexpected receipt validation: %d %s", response.Code, response.Body.String())
	}
}
