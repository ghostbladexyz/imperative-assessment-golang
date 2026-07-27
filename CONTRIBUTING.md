# Contributing to Imperative Go Practice Assessment

Thank you for helping improve the assessment. Contributions may include exercise corrections, runner and sandbox improvements, frontend changes, tests, and documentation.

## Getting Started

1. Fork the repository and clone your fork:

   ```sh
   git clone https://github.com/<your-user-name>/imperative-assessment-golang.git
   cd imperative-assessment-golang
   ```

2. Create a focused branch:

   ```sh
   git checkout -b fix/short-description
   ```

3. Make your changes and add or update tests for changed behavior.

## Development Requirements

- Go 1.23 or newer
- Docker Desktop or Docker Engine for the default sandbox and Docker checks
- Node.js 24 and npm when changing the frontend
- GNU Make is optional; every target maps to standard Go, npm, or Docker commands

Install the locked frontend dependencies when working under `web/`:

```sh
npm --prefix web ci
```

## Project Structure

- `cmd/server/` starts the local assessment server.
- `cmd/sandbox-runner/` contains the process used inside the Docker sandbox.
- `internal/assessment/` defines the exercise catalogue, authoring rules, and generated test harnesses.
- `internal/runner/` executes submissions through the local or Docker runner.
- `internal/web/` embeds the compiled frontend.
- `web/` contains the React and TypeScript frontend source.
- `docker/runner.Dockerfile` defines the pinned sandbox image.
- `.github/workflows/ci.yml` is the authoritative continuous-integration workflow.

## Code Standards

- Format Go changes with `gofmt`.
- Follow existing Go and TypeScript conventions.
- Keep changes focused and avoid unrelated formatting or refactoring.
- Add regression tests for bug fixes and tests for new behavior.
- Update existing documentation when user-facing commands or requirements change.
- If files under `web/` change, run the frontend build and commit the corresponding generated files under `internal/web/dist/`.

## Testing

Run the same main verification workflow used during development:

```sh
make check
```

Individual checks can also be run directly:

```sh
gofmt -w cmd internal
go vet ./...
go test -count=1 ./...
npm --prefix web run lint
npm --prefix web test
npm --prefix web run build
```

Changes affecting the sandbox should also build its image:

```sh
make docker-build
```

The opt-in Docker integration test is available with:

```sh
make docker-test
```

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org/) with a concise description:

```text
feat: add a new capability
fix(internal/assessment): correct an exercise case
docs(CONTRIBUTING.md): clarify the review workflow
test(internal/runner): cover a runner failure
```

Keep each commit coherent and include generated frontend assets in the same commit as their source changes.

## Submitting Changes

1. Rebase or merge the latest `main` into your branch.
2. Run the relevant checks, including `make check` when the frontend dependencies are installed.
3. Push the branch to your fork.
4. Open a pull request against `main`.
5. Explain the problem, the chosen solution, and how you verified it.

Pull requests must pass CI, remain within their stated scope, and include appropriate tests and documentation.

## Reporting Issues

Use the [issue tracker](https://github.com/ghostbladexyz/imperative-assessment-golang/issues) for reproducible bugs, exercise corrections, and feature proposals. Follow [SECURITY.md](SECURITY.md) when reporting a security concern.
