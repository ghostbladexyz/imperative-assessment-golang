package updatecheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

const defaultCacheLifetime = time.Hour

type compareResult struct {
	Status   string `json:"status"`
	AheadBy  int    `json:"ahead_by"`
	BehindBy int    `json:"behind_by"`
}

type cachedResult struct {
	Revision  string        `json:"revision"`
	CheckedAt time.Time     `json:"checkedAt"`
	Compare   compareResult `json:"compare"`
}

// Checker reports whether the running checkout is behind its upstream branch.
type Checker struct {
	owner         string
	repository    string
	branch        string
	cachePath     string
	client        *http.Client
	now           func() time.Time
	revision      func() (string, bool)
	cacheLifetime time.Duration
	apiBaseURL    string
}

// New creates a checker for a public GitHub repository and stores its cache in dataDir.
func New(owner, repository, branch, dataDir string) *Checker {
	return &Checker{
		owner:         owner,
		repository:    repository,
		branch:        branch,
		cachePath:     filepath.Join(dataDir, "update-check.json"),
		client:        http.DefaultClient,
		now:           time.Now,
		revision:      currentRevision,
		cacheLifetime: defaultCacheLifetime,
		apiBaseURL:    "https://api.github.com",
	}
}

// Notice returns a user-facing update message, or an empty string when no action is needed.
func (checker *Checker) Notice(ctx context.Context) string {
	revision, ok := checker.revision()
	if !ok {
		return ""
	}

	result, ok := checker.readCache(revision)
	if !ok {
		result, ok = checker.compare(ctx, revision)
		if !ok {
			return ""
		}
		checker.writeCache(revision, result)
	}

	switch {
	case result.AheadBy > 0 && result.BehindBy == 0:
		return fmt.Sprintf(
			"A newer version is available (%d %s behind).\n\nUpdate with:\n  git pull --ff-only",
			result.AheadBy,
			pluralize(result.AheadBy, "commit", "commits"),
		)
	case result.AheadBy > 0 && result.BehindBy > 0:
		return fmt.Sprintf(
			"The remote has %d new %s, but this checkout has diverged.\n\nReview the branches before updating:\n  git fetch origin\n  git log --oneline --left-right HEAD...origin/%s",
			result.AheadBy,
			pluralize(result.AheadBy, "commit", "commits"),
			checker.branch,
		)
	default:
		return ""
	}
}

// compare asks GitHub how the configured branch relates to revision.
func (checker *Checker) compare(ctx context.Context, revision string) (compareResult, bool) {
	endpoint := fmt.Sprintf(
		"%s/repos/%s/%s/compare/%s...%s?per_page=1",
		checker.apiBaseURL,
		url.PathEscape(checker.owner),
		url.PathEscape(checker.repository),
		url.PathEscape(revision),
		url.PathEscape(checker.branch),
	)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return compareResult{}, false
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("User-Agent", "imperative-go-assessment")

	response, err := checker.client.Do(request)
	if err != nil {
		return compareResult{}, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return compareResult{}, false
	}

	var result compareResult
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return compareResult{}, false
	}
	switch result.Status {
	case "identical", "ahead", "behind", "diverged":
		return result, true
	default:
		return compareResult{}, false
	}
}

// readCache reuses only a recent result for the exact revision that is running.
func (checker *Checker) readCache(revision string) (compareResult, bool) {
	data, err := os.ReadFile(checker.cachePath)
	if err != nil {
		return compareResult{}, false
	}
	var cached cachedResult
	if json.Unmarshal(data, &cached) != nil ||
		cached.Revision != revision ||
		checker.now().Sub(cached.CheckedAt) < 0 ||
		checker.now().Sub(cached.CheckedAt) > checker.cacheLifetime {
		return compareResult{}, false
	}
	return cached.Compare, true
}

// writeCache records successful checks so repeated launches do not exhaust GitHub's public rate limit.
func (checker *Checker) writeCache(revision string, result compareResult) {
	data, err := json.Marshal(cachedResult{
		Revision:  revision,
		CheckedAt: checker.now(),
		Compare:   result,
	})
	if err != nil {
		return
	}
	_ = os.WriteFile(checker.cachePath, data, 0o600)
}

// currentRevision reads the source commit embedded by the Go toolchain.
func currentRevision() (string, bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", false
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && setting.Value != "" {
			return setting.Value, true
		}
	}
	return "", false
}

// pluralize chooses the label that agrees with count.
func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}
