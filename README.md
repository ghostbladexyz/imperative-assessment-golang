# Imperative Go Practice Assessment

A complete local practice environment for a six-hour, nine-level Imperative Go assessment. It combines a browser-based Go editor with a localhost-only Go runner, visible automated tests, sequential progression, a persistent timer, and a final 7/9 readiness summary.

The questions are original practice exercises inspired by common Go curriculum areas. They are not copied from a Zone01 assessment.

## What is included

- Nine levels with 12–14 visible, deterministic tests each
- A protected `package main` line and Go-aware CodeMirror editor
- Real Go compilation and execution, readable compiler errors, timeouts, and output limits
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

For simply running the checked-in application:

- [Go](https://go.dev/dl/) 1.23 or newer

For changing and rebuilding the frontend:

- Go 1.23 or newer
- Node.js 20.19+ or 22.12+ (Node 24 is also supported)
- npm 10 or newer

The application uses the Go version installed on the student's machine. It has been verified with Go 1.25.4 and Node.js 24.12.0.

## Run the assessment

From the repository root:

```sh
go run ./cmd/server -open
```

If the browser does not open automatically, visit [http://127.0.0.1:8080](http://127.0.0.1:8080).

That is the only command students need. The production frontend is embedded in the Go server, so there is no frontend installation step for ordinary use.

To use another localhost port:

```sh
go run ./cmd/server -addr 127.0.0.1:9090 -open
```

## Development

Install frontend dependencies once:

```sh
cd web
npm install
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
go test ./...

cd web
npm run lint
npm test
npm run build
```

Build a standalone server executable:

```sh
go build -trimpath -o imperative-assessment ./cmd/server
```

On Windows the output is `imperative-assessment.exe`.

## How test execution works

1. The browser sends the current level ID, student code, and selected visible test IDs as JSON.
2. The server restores `package main`, writes the submission and a server-controlled harness into a unique temporary directory, and compiles them with the local Go toolchain.
3. The compiled temporary program executes the visible cases and emits structured result records.
4. The server compares actual and expected values, removes the temporary directory, and returns structured JSON to the browser.
5. A level only receives a signed pass receipt after the server runs every required test successfully.

Runs have a 20-second compile limit, a 4-second execution limit, a 256 KiB output limit, a 192 KiB source limit, and a two-run concurrency limit. Temporary executables are removed after every run.

The `TEST` button runs all visible tests. Each case has a run-one button, and **Run failed tests** reruns the current failures. Partial runs help iteration but cannot produce a pass receipt; use `TEST` for the authoritative complete result.

## Security limitations

This runner is for trusted local practice only. It is **not** a hardened code-execution sandbox:

- Student code executes with the permissions of the local user.
- The timeout, output cap, temporary directories, safe process arguments, request limit, and localhost binding reduce accidental damage; they are not OS-level isolation.
- Do not bind the server to a public interface or expose it to untrusted users.
- Do not deploy this runner as an internet service without a separate, strongly isolated execution system.

The server binds to `127.0.0.1` by default and displays the same warning in the application.

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
cmd/server/                 HTTP server, API validation, localhost startup
internal/assessment/        Level metadata, visible cases, hidden harness builders
internal/runner/            Temporary compilation, execution, limits, result parsing
internal/receipts/          Signed server-confirmed completion receipts
internal/web/dist/          Checked-in production frontend embedded by Go
web/src/                    React, TypeScript, CodeMirror, state, and UI
web/src/storage.ts          Versioned persistence and import validation
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
5. Run `go test ./...`. The definition test requires exactly nine levels and 10–18 tests per level.
6. Run the frontend checks if public metadata or UI behavior changed.

Test payloads and harness builders are not serialized by the levels API, so full controlled test code is not exposed in the frontend bundle. All tests required to pass are nevertheless listed visibly in the interface. There are currently no hidden pass requirements.

## Troubleshooting

### `go` is not recognized

Install Go, restart the terminal, and confirm:

```sh
go version
```

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
npm install
npm run build
```

## License

See [`LICENSE`](LICENSE).
