package runner

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/assessment"
)

const (
	dockerGoVersion      = "official grader toolchain"
	dockerImage          = "ghcr.io/zone01athens/test-imperative-checkpoint@sha256:6d3501ebb37c4c268177edaa7a0490e601898a25f6426310ad59fae85254b0d7"
	dockerPlatform       = "linux/amd64"
	dockerContainerLabel = "io.github.pleft.imperative-assessment.role=official-grader"
	dockerRunOutputLimit = 4 * 1024 * 1024
	dockerCleanupTimeout = 5 * time.Second
	dockerHostTimeout    = 90 * time.Second
	dockerStartupTimeout = 10 * time.Minute
)

var containerNamePattern = regexp.MustCompile(`^imperative-go-assessment-[a-f0-9]{24}$`)

type CommandResult struct {
	Stdout        string
	Stderr        string
	Err           error
	OutputLimited bool
}

type CommandExecutor interface {
	Run(ctx context.Context, outputLimit int, stdin []byte, name string, args ...string) CommandResult
}

type processExecutor struct{}

func (processExecutor) Run(ctx context.Context, outputLimit int, stdin []byte, name string, args ...string) CommandResult {
	stdout, stderr, err, limited := runCommandInput(ctx, outputLimit, stdin, name, args...)
	return CommandResult{Stdout: stdout, Stderr: stderr, Err: err, OutputLimited: limited}
}

type DockerOptions struct {
	DockerBinary  string
	ProjectRoot   string // Retained for API compatibility; the official image needs no build context.
	MaxConcurrent int
	Receipts      ReceiptIssuer
	Commands      CommandExecutor
	Random        io.Reader
	BuildTimeout  time.Duration
}

type dockerAdapter struct {
	dockerBinary string
	commands     CommandExecutor
	random       io.Reader
}

type dockerResourceMount struct {
	sourcePath string
	targetPath string
}

func NewDocker(ctx context.Context, options DockerOptions) (*Engine, error) {
	if options.DockerBinary == "" {
		options.DockerBinary = "docker"
	}
	if options.Commands == nil {
		options.Commands = processExecutor{}
	}
	if options.Random == nil {
		options.Random = rand.Reader
	}
	if options.BuildTimeout <= 0 {
		options.BuildTimeout = dockerStartupTimeout
	}
	if err := CheckDocker(ctx, options.Commands, options.DockerBinary); err != nil {
		return nil, err
	}
	if err := cleanupStaleContainers(ctx, options.Commands, options.DockerBinary); err != nil {
		return nil, err
	}
	inspect := options.Commands.Run(ctx, 64*1024, nil, options.DockerBinary, "image", "inspect", "--format", "{{.Id}}", dockerImage)
	if inspect.Err != nil {
		pullCtx, cancel := context.WithTimeout(ctx, options.BuildTimeout)
		defer cancel()
		pull := options.Commands.Run(pullCtx, dockerRunOutputLimit, nil, options.DockerBinary, "pull", "--platform", dockerPlatform, dockerImage)
		if pull.Err != nil {
			if errors.Is(pullCtx.Err(), context.DeadlineExceeded) {
				return nil, errors.New("downloading the pinned official grader timed out")
			}
			return nil, errors.New("the pinned official Zone01 grader could not be downloaded")
		}
	}
	return newEngine(&dockerAdapter{dockerBinary: options.DockerBinary, commands: options.Commands, random: options.Random}, options.MaxConcurrent, options.Receipts), nil
}

func CheckDocker(ctx context.Context, commands CommandExecutor, dockerBinary string) error {
	if cli := commands.Run(ctx, 64*1024, nil, dockerBinary, "--version"); cli.Err != nil {
		return errors.New("Docker CLI is unavailable. Install Docker Desktop, then rerun the assessment")
	}
	daemon := commands.Run(ctx, 64*1024, nil, dockerBinary, "info", "--format", "{{.ServerVersion}}")
	if daemon.Err != nil || strings.TrimSpace(daemon.Stdout) == "" {
		return errors.New("Docker Desktop is not running. Start Docker Desktop, wait until it is ready, then rerun the assessment")
	}
	return nil
}

func (docker *dockerAdapter) Info() Info {
	return Info{Mode: ModeDocker, SandboxReady: true, GoVersion: dockerGoVersion, DockerImage: dockerImage,
		Message: "Pinned official Zone01 grader ready. Each run uses a fresh restricted container."}
}

func (docker *dockerAdapter) Execute(ctx context.Context, plan executionPlan) executionOutcome {
	name, err := newContainerName(docker.random)
	if err != nil {
		return executionOutcome{status: executionInternal, runtimeError: "The grader could not create a unique run identifier."}
	}
	directory, err := os.MkdirTemp("", "imperative-official-grader-")
	if err != nil {
		return executionOutcome{status: executionInternal, runtimeError: "The grader could not prepare the submitted source."}
	}
	defer os.RemoveAll(directory)
	sourcePath := filepath.Join(directory, "main.go")
	if err := os.WriteFile(sourcePath, []byte(plan.prepared), 0o600); err != nil {
		return executionOutcome{status: executionInternal, runtimeError: "The grader could not prepare the submitted source."}
	}
	resourceMounts, err := stageOfficialResources(directory, string(plan.level.Key), plan.level.Resources)
	if err != nil {
		return executionOutcome{status: executionInternal, runtimeError: "The grader could not prepare the exercise resources."}
	}

	runCtx, cancel := context.WithTimeout(ctx, dockerHostTimeout)
	command := docker.commands.Run(runCtx, dockerRunOutputLimit, nil, docker.dockerBinary,
		dockerRunArgs(name, string(plan.level.Key), sourcePath, resourceMounts)...)
	cancel()
	cleanupErr := docker.cleanup(name)
	if cleanupErr != nil {
		return executionOutcome{status: executionCleanup, runtimeError: "The official grader container could not be removed."}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return executionOutcome{status: executionStopped}
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return executionOutcome{status: executionRuntimeTimeout, runtimeError: "The official grader exceeded its 90 second host limit."}
	}
	if command.OutputLimited {
		return executionOutcome{status: executionOutput, runtimeError: "Official grader output exceeded the response limit."}
	}
	outcome, err := decodeOfficialOutcome(plan.level, command.Stdout, command.Stderr)
	if err != nil {
		message := "The official grader returned an invalid response."
		if command.Err != nil {
			message = officialStartupMessage(command.Stdout, command.Stderr)
		}
		return executionOutcome{status: executionStartup, runtimeError: message, stdout: command.Stdout, stderr: command.Stderr}
	}
	return outcome
}

func officialStartupMessage(stdout, stderr string) string {
	combined := strings.ToLower(stdout + "\n" + stderr)
	if strings.Contains(combined, "exec format error") ||
		strings.Contains(combined, "no matching manifest") ||
		(strings.Contains(combined, "platform") && strings.Contains(combined, "emulat")) {
		return "The official grader requires Docker linux/amd64 emulation. Enable amd64 emulation in Docker, then rerun the assessment."
	}
	return "The official grader container could not start."
}

type officialEnvelope struct {
	OK     bool   `json:"Ok"`
	Output string `json:"Output"`
}

func decodeOfficialOutcome(level assessment.Level, stdout, stderr string) (executionOutcome, error) {
	normalized := strings.ReplaceAll(stdout, "\r\n", "\n")
	if strings.Contains(normalized, "Reason: code does not compile") {
		return executionOutcome{
			status: executionCompile, compileError: strings.TrimSpace(normalized),
			stdout: strings.TrimSpace(normalized), stderr: strings.TrimSpace(stderr),
		}, nil
	}
	lines := strings.Split(strings.TrimSpace(normalized), "\n")
	if len(lines) == 0 {
		return executionOutcome{}, errors.New("empty grader response")
	}
	var envelope officialEnvelope
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &envelope); err != nil {
		return executionOutcome{}, err
	}
	human := strings.TrimSpace(strings.Join(lines[:len(lines)-1], "\n"))
	outcome := executionOutcome{status: executionSuccess, stdout: compactOfficialOutput(human), stderr: strings.TrimSpace(stderr)}
	if strings.Contains(envelope.Output, "Reason: code does not compile") || strings.Contains(human, "Reason: code does not compile") {
		outcome.status = executionCompile
		outcome.compileError = strings.TrimSpace(envelope.Output)
		return outcome, nil
	}
	humanLines := strings.Split(human, "\n")
	for lineIndex, line := range humanLines {
		status, label, found := officialResultLine(line)
		if !found {
			continue
		}
		for _, test := range level.Tests {
			if !test.MatchesOfficialLabel(label) {
				continue
			}
			wire := wireResult{ID: test.ID, Actual: "fail", Input: officialObservedInput(humanLines, lineIndex)}
			if status == "PASS" {
				wire.Actual = "pass"
			} else {
				wire.Failure = officialFailureReason(line)
			}
			outcome.results = upsertWire(outcome.results, wire)
			break
		}
	}
	if len(level.Tests) == 1 && len(outcome.results) == 0 {
		wire := wireResult{ID: level.Tests[0].ID, Actual: "fail", Failure: "Official grader suite failed."}
		if envelope.OK {
			wire.Actual, wire.Failure = "pass", ""
		}
		outcome.results = append(outcome.results, wire)
	}
	if len(outcome.results) == 0 {
		return executionOutcome{}, errors.New("grader response contained no recognized checks")
	}
	if !envelope.OK {
		outcome.status = executionRuntime
		allPassed := true
		for _, result := range outcome.results {
			if result.Actual != "pass" || result.Failure != "" {
				allPassed = false
				break
			}
		}
		if allPassed {
			outcome.runtimeError = "The official grader reported a failed suite."
		}
	}
	return outcome, nil
}

func compactOfficialOutput(output string) string {
	var kept []string
	for _, line := range strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if status, label, found := officialResultLine(line); found {
			kept = append(kept, status+" "+label)
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(trimmed, "Exercise:") || strings.HasPrefix(trimmed, "RESULT:") ||
			strings.HasPrefix(trimmed, "Summary:") || strings.HasPrefix(lower, "replay seed:") ||
			strings.HasPrefix(trimmed, "Reason:") || strings.HasPrefix(trimmed, "== total:") {
			kept = append(kept, trimmed)
		}
	}
	return strings.Join(kept, "\n")
}

func officialFailureReason(line string) string {
	trimmed := strings.TrimSpace(line)
	if index := strings.Index(trimmed, " — "); index >= 0 {
		return strings.TrimSpace(trimmed[index+len(" — "):])
	}
	return "official check failed"
}

func officialObservedInput(lines []string, resultLine int) string {
	for index := resultLine + 1; index < len(lines); index++ {
		if _, _, found := officialResultLine(lines[index]); found {
			break
		}
		trimmed := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(trimmed, "Input (") {
			continue
		}
		var values []string
		for index++; index < len(lines); index++ {
			value := strings.TrimSpace(lines[index])
			if !strings.HasPrefix(value, "|") {
				break
			}
			value = strings.TrimSpace(strings.TrimPrefix(value, "|"))
			if value != "(empty)" {
				values = append(values, value)
			}
		}
		if len(values) == 0 {
			return "(empty)"
		}
		return strings.Join(values, "\n")
	}
	return ""
}

func officialResultLine(line string) (status, label string, found bool) {
	trimmed := strings.TrimSpace(line)
	for _, candidate := range []string{"PASS", "FAIL"} {
		if !strings.HasPrefix(trimmed, candidate+" ") {
			continue
		}
		label = strings.TrimSpace(strings.TrimPrefix(trimmed, candidate))
		for _, marker := range []string{" [", " —", " --"} {
			if index := strings.Index(label, marker); index >= 0 {
				label = strings.TrimSpace(label[:index])
			}
		}
		return candidate, label, label != ""
	}
	return "", "", false
}

func upsertWire(items []wireResult, wire wireResult) []wireResult {
	for index := range items {
		if items[index].ID == wire.ID {
			if wire.Actual == "fail" {
				items[index] = wire
			}
			return items
		}
	}
	return append(items, wire)
}

func stageOfficialResources(directory, exerciseKey string, resources []assessment.ExerciseResource) ([]dockerResourceMount, error) {
	slug := strings.TrimPrefix(exerciseKey, "checkpoint/")
	targetDirectory := "/jail/student/" + slug
	mounts := make([]dockerResourceMount, 0, len(resources))
	for _, resource := range resources {
		if resource.Name == "" || resource.Name == "." || resource.Name == ".." || strings.ContainsAny(resource.Name, "/\\") {
			return nil, fmt.Errorf("invalid official resource name %q", resource.Name)
		}
		sourcePath := filepath.Join(directory, resource.Name)
		if err := os.WriteFile(sourcePath, []byte(resource.Content), 0o600); err != nil {
			return nil, err
		}
		mounts = append(mounts, dockerResourceMount{
			sourcePath: sourcePath,
			targetPath: targetDirectory + "/" + resource.Name,
		})
	}
	return mounts, nil
}

func dockerRunArgs(name, exerciseKey, sourcePath string, resourceMounts []dockerResourceMount) []string {
	slug := strings.TrimPrefix(exerciseKey, "checkpoint/")
	target := "/jail/student/" + slug + "/main.go"
	args := []string{
		"run", "--name", name, "--label", dockerContainerLabel, "--rm", "--pull", "never",
		"--platform", dockerPlatform, "--network", "none", "--ipc", "none", "--read-only", "--log-driver", "none", "--hostname", "grader",
		"--memory", "512m", "--memory-swap", "512m", "--cpus", "1", "--pids-limit", "256",
		"--ulimit", "nofile=256:256", "--ulimit", "core=0:0",
		"--tmpfs", "/tmp:rw,exec,nosuid,nodev,size=256m,mode=1777",
		"--env", "EXERCISE=" + slug, "--env", "FILE=" + slug + "/main.go", "--env", "EMIT_JSON=1",
	}
	args = append(args, "--mount", "type=bind,source="+sourcePath+",target="+target+",readonly")
	for _, resource := range resourceMounts {
		args = append(args, "--mount", "type=bind,source="+resource.sourcePath+",target="+resource.targetPath+",readonly")
	}
	return append(args, dockerImage)
}

func cleanupStaleContainers(ctx context.Context, commands CommandExecutor, dockerBinary string) error {
	found := commands.Run(ctx, 64*1024, nil, dockerBinary, "ps", "--all", "--filter", "label="+dockerContainerLabel,
		"--filter", "status=created", "--filter", "status=exited", "--filter", "status=dead", "--format", "{{.Names}}")
	if found.Err != nil {
		return errors.New("the Docker grader could not check for stale containers")
	}
	cleaner := &dockerAdapter{dockerBinary: dockerBinary, commands: commands}
	for _, name := range strings.Fields(found.Stdout) {
		if !containerNamePattern.MatchString(name) || cleaner.cleanup(name) != nil {
			return errors.New("the Docker grader could not remove a stale container")
		}
	}
	return nil
}

func (docker *dockerAdapter) cleanup(name string) error {
	if !containerNamePattern.MatchString(name) {
		return errors.New("invalid container name")
	}
	ctx, cancel := context.WithTimeout(context.Background(), dockerCleanupTimeout)
	defer cancel()
	result := docker.commands.Run(ctx, 64*1024, nil, docker.dockerBinary, "rm", "--force", name)
	if result.Err == nil || strings.Contains(result.Stdout+result.Stderr, "No such container") {
		return nil
	}
	return errors.New("container cleanup failed")
}

func newContainerName(random io.Reader) (string, error) {
	value := make([]byte, 12)
	if _, err := io.ReadFull(random, value); err != nil {
		return "", err
	}
	name := "imperative-go-assessment-" + hex.EncodeToString(value)
	if !containerNamePattern.MatchString(name) {
		return "", fmt.Errorf("generated invalid container name")
	}
	return name, nil
}
