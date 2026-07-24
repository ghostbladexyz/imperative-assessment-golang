# Imperative Go Practice Assessment

A complete local practice environment for a six-hour, nine-level Imperative Go assessment. It combines a browser-based Go editor with a Docker-sandboxed runner, visible automated tests, sequential progression, a persistent timer, and a final 7/9 readiness summary.

The questions are original practice exercises inspired by common Go curriculum areas. They are not copied from a Zone01 assessment.

## What is included

- Nine levels with 12–14 visible, deterministic tests each
- A protected `package main` line and Go-aware CodeMirror editor
- Docker-sandboxed compilation and execution by default
- An explicit trusted local runner for environments where isolation is not required
- Readable compiler errors, separate compile/runtime timeouts, and bounded output
- Per-test input, expected output, actual output, duration, status, and diffs
- Sequential unlocking plus an explicitly labelled practice mode
- Six-hour timer with pause/reset and reload persistence
- Local code saving, attempts, time per level, hints, and recent test history
- Server-signed pass receipts that are revalidated when progress is imported
- JSON progress backup and restore
- Light/dark themes, reduced motion, keyboard access, and responsive layouts
- A final summary with readiness status and suggested review topics

Levels 1–7 form the required path. Levels 8–9 are optional SQL-boundary and integrated stretch exercises.

## Prerequisites

For the standard Docker-sandboxed experience:

- [Go](https://go.dev/dl/) 1.23 or newer
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or a compatible Docker Engine
- A running Docker daemon accessible through the `docker` CLI

For changing and rebuilding the frontend:

- Go 1.23 or newer
- Node.js 20.19+ or 22.12+ (Node 24 is supported)
- npm 10 or newer

The sandbox image contains Go 1.23.12. The host Go toolchain only starts the server. The project has been verified with host Go 1.25.4, Docker Desktop 29.2.1, Node.js 24.12.0, and npm 11.6.2.

## Run the assessment

From the repository root:

```sh
go run ./cmd/server -open
```

This is the only command students need. It:

1. Checks the Docker CLI and daemon.
2. Calculates a deterministic image tag from the image inputs.
3. Builds the pinned runner image when that exact tag is not cached.
4. Starts the localhost server only after the sandbox is ready.
5. Opens the browser.

The first launch may take a few minutes while Docker downloads the pinned base image and precompiles the controlled standard-library dependency cache. Later launches reuse the image whenever `.dockerignore`, `docker/runner.Dockerfile`, `cmd/sandbox-runner`, `internal/sandboxprotocol`, and `go.mod` are unchanged. A cached launch works offline; a first build needs network access to retrieve the pinned base image.

If the browser does not open automatically, visit [http://127.0.0.1:8080](http://127.0.0.1:8080).

Explicit Docker selection is also supported:

```sh
go run ./cmd/server -runner docker -open
```

Docker is always the default. If it is unavailable, startup stops with an actionable error and never silently falls back to local execution.

### Trusted local runner

Trusted users may explicitly run code with their own OS permissions:

```sh
go run ./cmd/server -runner local -open
```

Local mode preserves the same temporary directories, request/source limits, compile/runtime deadlines, output caps, concurrency limit, grading, cleanup, cancellation, and signed receipts. It is not a sandbox. Only use it for code you trust.

To use another localhost port:

```sh
go run ./cmd/server -addr 127.0.0.1:9090 -open
```

Keep the service bound to localhost. The browser displays the server-selected mode; it cannot select or override the runner.

## How test execution works

The server depends on a small runner interface implemented by the Docker and local engines:

1. The browser sends a level ID, student code, and selected visible test IDs.
2. The server restores `package main`, selects the server-owned test cases, builds the controlled harness, and enforces the shared concurrency limit.
3. In Docker mode, the server sends one bounded JSON request over standard input to a fresh container. No submitted source is placed in shell arguments, environment variables, or host paths.
4. The fixed `sandbox-runner` validates the request, copies the immutable compiler cache into a bounded tmpfs, creates files in its temporary workspace, formats the source when possible, compiles, and executes with separate deadlines.
5. The entrypoint captures bounded output and returns exactly one structured JSON response. The server remains authoritative for assertions and receipts.
6. The server force-removes the exact per-run container under an independent cleanup deadline, including after success, failure, timeout, output overflow, cancellation, or Docker CLI failure.

Every Docker submission uses a unique disposable container. The web server and UI remain on the host and are never placed inside it. Runs have a 20-second compile limit, a 4-second execution limit, a 256 KiB output limit, a 192 KiB source limit, and a two-run server concurrency limit.

The `TEST` button runs all visible tests. Each case has a run-one button, and **Run failed tests** reruns the current failures. Partial runs help iteration but cannot produce a pass receipt; use `TEST` for the authoritative complete result.

## Docker isolation model

Each submission container is started with:

- No network
- A read-only root filesystem
- All Linux capabilities dropped
- `no-new-privileges`
- 256 MiB memory and swap limit
- One CPU
- A 64-process limit
- A 128-file descriptor limit
- A dedicated non-root UID/GID
- Separate, size-limited tmpfs mounts for temporary files, the compiler cache, and the executable workspace
- `CGO_ENABLED=0`, no module proxy, and a locally pinned toolchain

The repository, home directory, Docker socket, credentials, browser data, and arbitrary host paths are never mounted. The image contains the pinned Go toolchain, the fixed sandbox entrypoint, and a standard-library build cache. It contains no assessment solutions, receipt key, Docker CLI, Docker socket, or user credentials.

Docker significantly improves isolation compared with the local runner, but this project is still intended for local and classroom use. Containers depend on Docker Desktop or Engine, its VM where applicable, the host kernel/runtime, and correct Docker configuration. The Docker daemon is trusted infrastructure.

Do not expose this as an anonymous internet code-execution API. A public multi-tenant service would additionally require per-submission microVMs or equivalent isolation, authentication, rate limiting, separate execution infrastructure, monitoring, abuse controls, and an independent security review. See [`SECURITY.md`](SECURITY.md).

## Docker image lifecycle

The image repository is:

```text
imperative-go-assessment-runner
```

Its tag is the first 20 hexadecimal characters of a SHA-256 digest over normalized input paths and contents. The Dockerfile pins `golang:1.23.12-alpine3.22` by immutable OCI digest, not only by tag.

To update Go:

1. Choose an exact official Alpine-based Go image tag.
2. Obtain and verify its current OCI digest, for example:

   ```sh
   docker buildx imagetools inspect golang:1.23.12-alpine3.22
   ```

3. Update both `FROM` lines in `docker/runner.Dockerfile` with the exact tag and digest.
4. Update `dockerGoVersion` in `internal/runner/docker.go`.
5. Run the unit, integration, and image-build checks below.

Changing an image input automatically produces a new project image tag. Old cached tags are not removed automatically.

List only this project's cached images:

```sh
docker image ls --filter reference=imperative-go-assessment-runner --format "{{.Repository}}:{{.Tag}}"
```

Remove a specific tag copied from that output:

```sh
docker image rm imperative-go-assessment-runner:EXACT_TAG
```

Never use broad Docker prune commands for assessment cleanup.

Verify that no submission containers remain:

```sh
docker ps -a --filter name=imperative-go-assessment- --format "{{.Names}} {{.Status}}"
```

An empty result is expected when no run is active.

## Development

Install the exact frontend dependency lock:

```sh
cd web
npm ci
```

Build the frontend into the Go embed directory:

```sh
cd web
npm run build
cd ..
```

Run the complete application:

```sh
go run ./cmd/server
```

For live frontend development, run these in separate terminals:

```sh
go run ./cmd/server
```

```sh
cd web
npm run dev
```

Open [http://127.0.0.1:5173](http://127.0.0.1:5173). Vite proxies `/api` to the Go server.

### Quality commands

```sh
gofmt -w cmd internal
go vet ./...
go test -count=1 ./...

cd web
npm ci
npm run lint
npm test
npm run build
```

Build the real runner image:

```sh
docker build --file docker/runner.Dockerfile --tag imperative-go-assessment-runner:manual-check .
```

Run the opt-in real-Docker integration suite:

```sh
IMPERATIVE_DOCKER_INTEGRATION=1 go test -count=1 -run TestDockerRunnerIntegration -v ./internal/runner
```

PowerShell:

```powershell
$env:IMPERATIVE_DOCKER_INTEGRATION = "1"
go test -count=1 -run TestDockerRunnerIntegration -v ./internal/runner
```

The integration suite skips with a clear reason unless explicitly enabled. It verifies correct and incorrect submissions, compiler errors, runtime timeout, cancellation cleanup, network isolation, read-only root, absent host mounts, output capping, container removal, and all nine starter/harness combinations.

Build a standalone server executable:

```sh
go build -trimpath -o imperative-assessment ./cmd/server
```

On Windows the output is `imperative-assessment.exe`. Do not commit it.

## Progress and timer storage

Browser progress uses versioned local storage under:

```text
imperative-go-assessment:progress:v1
```

It stores code independently per level, settings, timer state, attempts, time spent, hint usage, and recent run summaries. The server stores a local receipt signing key in the operating system's user configuration directory:

- Windows: `%AppData%\imperative-go-assessment`
- macOS: the normal user config directory returned by Go
- Linux: `$XDG_CONFIG_HOME/imperative-go-assessment` or the Go default

Set `IMPERATIVE_ASSESSMENT_DATA` to override the server data directory.

Use **Export** for a JSON backup and **Import** to restore it. Imports are schema-validated and claimed passes are revalidated against server-issued receipts. Use **Reset level** for starter code, or **Reset all** and type `RESET` to clear the entire browser assessment.

## Keyboard shortcuts and editor controls

- `Ctrl+Enter` / `Cmd+Enter`: run all visible tests
- `Ctrl+F` / `Cmd+F`: find in the editor
- `Ctrl+Z` / `Cmd+Z`: undo
- `Ctrl+Shift+Z` / `Cmd+Shift+Z`: redo
- `Tab` / `Shift+Tab`: indent/outdent

The toolbar also provides copy, `.go` download, `gofmt`, font-size controls, and a starter-code reset. `package main` is displayed as a protected first line and is restored by the server before formatting or execution.

## Project structure

```text
cmd/server/                   Startup, runner selection, HTTP API, localhost UI
cmd/sandbox-runner/           Fixed JSON container entrypoint
docker/runner.Dockerfile      Digest-pinned purpose-built execution image
internal/assessment/          Level metadata, visible cases, harness builders
internal/runner/              Shared contract, local engine, Docker orchestration
internal/sandboxprotocol/     Bounded container request/response schema
internal/receipts/            Signed server-confirmed completion receipts
internal/web/dist/            Checked-in production frontend embedded by Go
web/src/                      React, TypeScript, CodeMirror, state, and UI
.github/workflows/ci.yml      Go, frontend, and Docker image checks
```

## Assessment authoring guide

The central authoring file is [`internal/assessment/levels.go`](internal/assessment/levels.go). Each level has:

- Public metadata and instructions
- Starter code without `package main`
- Visible test metadata and an internal payload
- A harness builder that turns server-owned payloads into controlled Go calls

To edit a question or add a test:

1. Update the corresponding `levelOne` through `levelNine` definition.
2. Keep the function signature, starter code, contract, examples, and tests consistent.
3. Add a focused `VisibleTest` with a unique ID, meaningful name, purpose, visible input, expected output, and internal payload.
4. Update the level-specific harness only if the payload shape or required API changes.
5. If the level introduces a standard-library package not already precompiled, add it to the cache-building `go install` command in `docker/runner.Dockerfile`.
6. Run the unit and real-Docker integration suites.
7. Rebuild the frontend if public metadata or UI behavior changed.

Test payloads and harness builders are not serialized by the levels API, so controlled test code is not exposed in the frontend bundle. All tests required to pass are nevertheless listed visibly in the interface. There are currently no hidden pass requirements.

## Troubleshooting

### Docker is unavailable

Confirm both commands succeed:

```sh
docker --version
docker info
```

If the CLI exists but `docker info` fails, start Docker Desktop and wait until its engine reports ready. On Windows, confirm Docker Desktop is using Linux containers. The server never falls back to local execution automatically.

### The first image build fails

Confirm Docker Desktop has enough disk space and can reach Docker Hub to download the pinned base image. Retry the standard startup command. The browser is not started until the image is ready. For full Docker build output, run the manual image-build command in **Quality commands**.

### Cached/offline launch fails

List the project image tags with the command in **Docker image lifecycle**. If no tag matches the current image inputs, Docker must build a new image and needs the pinned base image locally or network access.

### Port 8080 is already in use

Choose another localhost port:

```sh
go run ./cmd/server -addr 127.0.0.1:9090 -open
```

### A correct-looking answer does not pass

Open the failed test and compare **Expected** and **Actual** exactly. Review the whitespace rule and constraints in the instructions. Run `gofmt`, then use the full `TEST` action rather than only a partial run.

### Progress disappeared or imports do not restore passes

Code and settings live in the current browser profile. Import the most recent JSON backup. Pass receipts are installation-specific; if the server signing key was deleted or the backup came from another installation, the code is kept but imported completion claims are intentionally not trusted. Rerun all tests to issue new receipts.

### The frontend says it has not been built

The checked-in `internal/web/dist` should already contain the app. If it was removed, rebuild it:

```sh
cd web
npm ci
npm run build
```

## License

See [`LICENSE`](LICENSE).
