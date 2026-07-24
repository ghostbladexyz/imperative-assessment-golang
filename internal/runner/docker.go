package runner

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pleft/imperative-assessment-golang/internal/sandboxprotocol"
)

const (
	dockerGoVersion        = "go1.23.12"
	dockerImageRepository  = "imperative-go-assessment-runner"
	dockerBuildOutputLimit = 2 * 1024 * 1024
	dockerRunOutputLimit   = 2 * 1024 * 1024
	dockerCleanupTimeout   = 5 * time.Second
	dockerHostTimeout      = CompileTimeout + RuntimeTimeout + 15*time.Second
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
	ProjectRoot   string
	MaxConcurrent int
	Receipts      ReceiptIssuer
	Commands      CommandExecutor
	Random        io.Reader
	BuildTimeout  time.Duration
}

type dockerAdapter struct {
	dockerBinary string
	image        string
	commands     CommandExecutor
	random       io.Reader
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
		options.BuildTimeout = 10 * time.Minute
	}
	if err := CheckDocker(ctx, options.Commands, options.DockerBinary); err != nil {
		return nil, err
	}
	root := options.ProjectRoot
	if root == "" {
		var err error
		root, err = findProjectRoot()
		if err != nil {
			return nil, err
		}
	}
	hash, err := imageInputHash(root)
	if err != nil {
		return nil, fmt.Errorf("inspect Docker sandbox inputs: %w", err)
	}
	image := dockerImageRepository + ":" + hash
	cached := options.Commands.Run(
		ctx,
		64*1024,
		nil,
		options.DockerBinary,
		"image",
		"inspect",
		"--format",
		"{{.Id}}",
		image,
	)
	if cached.Err != nil {
		buildCtx, cancelBuild := context.WithTimeout(ctx, options.BuildTimeout)
		defer cancelBuild()
		built := options.Commands.Run(
			buildCtx,
			dockerBuildOutputLimit,
			nil,
			options.DockerBinary,
			dockerBuildArgs(root, image)...,
		)
		if built.Err != nil {
			switch {
			case errors.Is(buildCtx.Err(), context.DeadlineExceeded):
				return nil, errors.New("the Docker sandbox image build timed out; restart Docker Desktop and try again")
			case built.OutputLimited:
				return nil, errors.New("the Docker sandbox image build produced too much output; check Docker Desktop and try again")
			default:
				return nil, errors.New("the Docker sandbox image could not be built; verify Docker Desktop has network access for the first build and try again")
			}
		}
	}
	return newEngine(&dockerAdapter{
		dockerBinary: options.DockerBinary,
		image:        image,
		commands:     options.Commands,
		random:       options.Random,
	}, options.MaxConcurrent, options.Receipts), nil
}

func CheckDocker(ctx context.Context, commands CommandExecutor, dockerBinary string) error {
	cli := commands.Run(ctx, 64*1024, nil, dockerBinary, "--version")
	if cli.Err != nil {
		return errors.New("Docker CLI is unavailable. Install Docker Desktop, then rerun the assessment")
	}
	daemon := commands.Run(ctx, 64*1024, nil, dockerBinary, "info", "--format", "{{.ServerVersion}}")
	if daemon.Err != nil || strings.TrimSpace(daemon.Stdout) == "" {
		return errors.New("Docker Desktop is not running. Start Docker Desktop, wait until it is ready, then rerun the assessment")
	}
	return nil
}

func (docker *dockerAdapter) Info() Info {
	return Info{
		Mode:         ModeDocker,
		SandboxReady: true,
		GoVersion:    dockerGoVersion,
		DockerImage:  docker.image,
		Message:      "Docker sandbox ready. Each run uses a fresh, restricted container.",
	}
}

func (docker *dockerAdapter) Execute(ctx context.Context, plan executionPlan) executionOutcome {
	name, err := newContainerName(docker.random)
	if err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "The sandbox could not create a unique run identifier.",
		}
	}
	tests := make([]sandboxprotocol.ExpectedTest, 0, len(plan.tests))
	for _, test := range plan.tests {
		tests = append(tests, sandboxprotocol.ExpectedTest{ID: test.ID, Expected: test.Expected})
	}
	request := sandboxprotocol.Request{
		Source:           plan.prepared,
		Harness:          plan.harness,
		Tests:            tests,
		CompileTimeoutMS: CompileTimeout.Milliseconds(),
		RuntimeTimeoutMS: RuntimeTimeout.Milliseconds(),
		OutputLimitBytes: MaxOutputBytes,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "The sandbox request could not be prepared.",
		}
	}

	runCtx, cancelRun := context.WithTimeout(ctx, dockerHostTimeout)
	commandResult := docker.commands.Run(
		runCtx,
		dockerRunOutputLimit,
		payload,
		docker.dockerBinary,
		dockerRunArgs(name, docker.image)...,
	)
	cancelRun()

	cleanupErr := docker.cleanup(name)
	if cleanupErr != nil {
		return executionOutcome{
			status:       executionCleanup,
			runtimeError: "The sandbox container could not be removed. Close Docker Desktop and remove only the named assessment container before trying again.",
		}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return executionOutcome{status: executionStopped}
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return executionOutcome{
			status:       executionStartup,
			runtimeError: "The sandbox did not finish within its overall time limit and was terminated.",
		}
	}
	if commandResult.OutputLimited {
		return executionOutcome{
			status:       executionOutput,
			runtimeError: "Sandbox output exceeded the response limit and execution was terminated.",
		}
	}

	response, decodeErr := decodeSandboxResponse(commandResult.Stdout)
	if decodeErr != nil {
		message := "The sandbox returned an invalid response. Restart the assessment server and try again."
		if commandResult.Err != nil {
			message = "The sandbox container could not start. Restart Docker Desktop and try again."
		}
		return executionOutcome{status: executionStartup, runtimeError: message}
	}
	if err := response.Validate(tests); err != nil {
		return executionOutcome{
			status: executionInternal, runtimeError: "The sandbox returned an invalid test result.",
		}
	}

	wires := make([]wireResult, 0, len(response.Results))
	for _, item := range response.Results {
		wires = append(wires, wireResult{
			ID: item.ID, Actual: item.Actual, Failure: item.Failure, DurationMS: item.DurationMS,
		})
	}
	outcome := executionOutcome{
		status:          executionSuccess,
		compileError:    response.CompileError,
		runtimeError:    response.RuntimeError,
		stdout:          response.Stdout,
		stderr:          response.Stderr,
		formattedSource: response.FormattedSource,
		results:         wires,
	}
	switch response.Status {
	case sandboxprotocol.StatusCompileTimeout:
		outcome.status = executionCompileTimeout
	case sandboxprotocol.StatusCompile:
		outcome.status = executionCompile
	case sandboxprotocol.StatusRuntimeTimeout:
		outcome.status = executionRuntimeTimeout
	case sandboxprotocol.StatusStopped:
		outcome.status = executionStopped
	case sandboxprotocol.StatusOutput:
		outcome.status = executionOutput
	case sandboxprotocol.StatusInternal:
		outcome.status = executionInternal
		if outcome.runtimeError == "" {
			outcome.runtimeError = "The sandbox could not complete this run."
		}
	case sandboxprotocol.StatusRuntime:
		outcome.status = executionRuntime
	case sandboxprotocol.StatusSuccess, sandboxprotocol.StatusAssertion:
	}
	return outcome
}

func (docker *dockerAdapter) cleanup(name string) error {
	if !containerNamePattern.MatchString(name) {
		return errors.New("invalid container name")
	}
	ctx, cancel := context.WithTimeout(context.Background(), dockerCleanupTimeout)
	defer cancel()
	cleanup := docker.commands.Run(ctx, 64*1024, nil, docker.dockerBinary, "rm", "--force", name)
	if cleanup.Err == nil {
		return nil
	}
	combined := cleanup.Stdout + "\n" + cleanup.Stderr
	if strings.Contains(combined, "No such container") || strings.Contains(combined, "not found") {
		return nil
	}
	return errors.New("container cleanup failed")
}

func decodeSandboxResponse(raw string) (sandboxprotocol.Response, error) {
	var response sandboxprotocol.Response
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&response); err != nil {
		return response, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return response, errors.New("sandbox response must contain exactly one JSON value")
	}
	return response, nil
}

func newContainerName(random io.Reader) (string, error) {
	value := make([]byte, 12)
	if _, err := io.ReadFull(random, value); err != nil {
		return "", err
	}
	name := "imperative-go-assessment-" + hex.EncodeToString(value)
	if !containerNamePattern.MatchString(name) {
		return "", errors.New("generated invalid container name")
	}
	return name, nil
}

func dockerRunArgs(name, image string) []string {
	return []string{
		"run",
		"--name", name,
		"--rm",
		"--interactive",
		"--network", "none",
		"--read-only",
		"--cap-drop", "ALL",
		"--security-opt", "no-new-privileges",
		"--memory", "256m",
		"--memory-swap", "256m",
		"--cpus", "1",
		"--pids-limit", "64",
		"--ulimit", "nofile=128:128",
		"--tmpfs", "/tmp:rw,noexec,nosuid,nodev,size=64m,mode=1777",
		"--tmpfs", "/tmp/go-build:rw,noexec,nosuid,nodev,size=96m,uid=65532,gid=65532,mode=0700",
		"--tmpfs", "/workspace:rw,exec,nosuid,nodev,size=32m,uid=65532,gid=65532,mode=0700",
		"--user", "65532:65532",
		"--env", "CGO_ENABLED=0",
		"--env", "HOME=/tmp",
		"--env", "GOCACHE=/tmp/go-build",
		"--env", "GOTMPDIR=/tmp",
		"--env", "GOPATH=/tmp/go",
		"--env", "GOPROXY=off",
		"--env", "GOTOOLCHAIN=local",
		"--env", "GOMAXPROCS=1",
		"--env", "GOMEMLIMIT=128MiB",
		"--env", "GOGC=50",
		"--workdir", "/workspace",
		image,
	}
}

func dockerBuildArgs(root, image string) []string {
	return []string{
		"build",
		"--file", filepath.Join(root, "docker", "runner.Dockerfile"),
		"--tag", image,
		root,
	}
}

func findProjectRoot() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if regularFile(filepath.Join(current, "go.mod")) &&
			regularFile(filepath.Join(current, "docker", "runner.Dockerfile")) {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("could not find the project root containing docker/runner.Dockerfile")
		}
		current = parent
	}
}

func regularFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func imageInputHash(root string) (string, error) {
	inputs := []string{
		".dockerignore",
		"go.mod",
		filepath.Join("docker", "runner.Dockerfile"),
	}
	for _, directory := range []string{
		filepath.Join("cmd", "sandbox-runner"),
		filepath.Join("internal", "sandboxprotocol"),
	} {
		err := filepath.WalkDir(filepath.Join(root, directory), func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") {
				relative, err := filepath.Rel(root, path)
				if err != nil {
					return err
				}
				inputs = append(inputs, relative)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	sort.Strings(inputs)
	hash := sha256.New()
	for _, relative := range inputs {
		content, err := os.ReadFile(filepath.Join(root, relative))
		if err != nil {
			return "", err
		}
		normalized := filepath.ToSlash(relative)
		_, _ = io.WriteString(hash, normalized)
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write(content)
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))[:20], nil
}
