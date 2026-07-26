package updatecheck

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestNoticeReportsAvailableUpdate verifies that an upstream-only change produces a safe pull command.
func TestNoticeReportsAvailableUpdate(t *testing.T) {
	checker, requests, closeServer := testChecker(t, `{"status":"ahead","ahead_by":2,"behind_by":0}`)
	defer closeServer()

	notice := checker.Notice(context.Background())
	if !strings.Contains(notice, "2 commits behind") ||
		!strings.Contains(notice, "git pull --ff-only") {
		t.Fatalf("unexpected notice %q", notice)
	}
	if *requests != 1 {
		t.Fatalf("requests = %d, want 1", *requests)
	}
}

// TestNoticeReportsDivergedCheckout verifies that local commits do not receive an unsafe fast-forward instruction.
func TestNoticeReportsDivergedCheckout(t *testing.T) {
	checker, _, closeServer := testChecker(t, `{"status":"diverged","ahead_by":3,"behind_by":1}`)
	defer closeServer()

	notice := checker.Notice(context.Background())
	if !strings.Contains(notice, "checkout has diverged") ||
		!strings.Contains(notice, "git fetch origin") ||
		strings.Contains(notice, "git pull --ff-only") {
		t.Fatalf("unexpected notice %q", notice)
	}
}

// TestNoticeCachesSuccessfulComparison verifies that repeated launches for one revision reuse the GitHub response.
func TestNoticeCachesSuccessfulComparison(t *testing.T) {
	checker, requests, closeServer := testChecker(t, `{"status":"ahead","ahead_by":1,"behind_by":0}`)
	defer closeServer()

	first := checker.Notice(context.Background())
	second := checker.Notice(context.Background())
	if first == "" || second != first {
		t.Fatalf("notices = %q and %q", first, second)
	}
	if *requests != 1 {
		t.Fatalf("requests = %d, want 1", *requests)
	}
}

// TestNoticeIgnoresUnavailableOrInvalidResponses verifies that update checks remain best-effort.
func TestNoticeIgnoresUnavailableOrInvalidResponses(t *testing.T) {
	for name, response := range map[string]struct {
		status int
		body   string
	}{
		"server error":     {status: http.StatusServiceUnavailable, body: `{}`},
		"invalid JSON":     {status: http.StatusOK, body: `{`},
		"invalid response": {status: http.StatusOK, body: `{"status":"unknown"}`},
	} {
		t.Run(name, func(t *testing.T) {
			checker, _, closeServer := testCheckerWithStatus(t, response.status, response.body)
			defer closeServer()
			if notice := checker.Notice(context.Background()); notice != "" {
				t.Fatalf("notice = %q, want empty", notice)
			}
		})
	}
}

// TestNoticeSkipsBuildsWithoutRevision verifies that source archives and unsupported builds do not contact GitHub.
func TestNoticeSkipsBuildsWithoutRevision(t *testing.T) {
	checker, requests, closeServer := testChecker(t, `{"status":"ahead","ahead_by":1,"behind_by":0}`)
	defer closeServer()
	checker.revision = func() (string, bool) { return "", false }

	if notice := checker.Notice(context.Background()); notice != "" {
		t.Fatalf("notice = %q, want empty", notice)
	}
	if *requests != 0 {
		t.Fatalf("requests = %d, want 0", *requests)
	}
}

// testChecker creates a checker backed by a local HTTP server.
func testChecker(t *testing.T, body string) (*Checker, *int32, func()) {
	t.Helper()
	return testCheckerWithStatus(t, http.StatusOK, body)
}

// testCheckerWithStatus creates a checker whose GitHub adapter returns status and body.
func testCheckerWithStatus(t *testing.T, status int, body string) (*Checker, *int32, func()) {
	t.Helper()
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		atomic.AddInt32(&requests, 1)
		if request.Header.Get("User-Agent") != "imperative-go-assessment" {
			t.Errorf("User-Agent = %q", request.Header.Get("User-Agent"))
		}
		if !strings.Contains(request.URL.Path, "/compare/abc123...main") {
			t.Errorf("path = %q", request.URL.Path)
		}
		writer.WriteHeader(status)
		_, _ = fmt.Fprint(writer, body)
	}))

	checker := New("owner", "repository", "main", t.TempDir())
	checker.apiBaseURL = server.URL
	checker.revision = func() (string, bool) { return "abc123", true }
	checker.now = func() time.Time { return time.Date(2026, 7, 26, 12, 0, 0, 0, time.UTC) }
	return checker, &requests, server.Close
}
