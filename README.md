# Imperative Checkpoint Practice Assessment

[![CI](https://github.com/terry-xyz/imperative-assessment-golang/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/terry-xyz/imperative-assessment-golang/actions/workflows/ci.yml)
[![Latest release](https://img.shields.io/github/v/release/terry-xyz/imperative-assessment-golang?display_name=tag)](https://github.com/terry-xyz/imperative-assessment-golang/releases)
[![License](https://img.shields.io/github/license/terry-xyz/imperative-assessment-golang)](https://github.com/terry-xyz/imperative-assessment-golang/blob/main/LICENSE)
[![Go 1.23+](https://img.shields.io/badge/go-1.23%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/doc/install)
[![Go Report Card](https://goreportcard.com/badge/github.com/terry-xyz/imperative-assessment-golang)](https://goreportcard.com/report/github.com/terry-xyz/imperative-assessment-golang)
[![Makefile](https://img.shields.io/badge/build-Makefile-427819?logo=gnu&logoColor=white)](https://github.com/terry-xyz/imperative-assessment-golang/blob/main/Makefile)

A local, browser-based practice assessment with the 17 Imperative Checkpoint
exercises, automatic tests, console feedback, and saved progress. The exercises
are arranged from easiest to hardest, and submitted code is checked by the
pinned official Zone01 grader in a restricted Docker container.

This is an unofficial practice tool and is not affiliated with Zone01.
The exercise subjects, starter programs, and supplied resources are copied
verbatim from the
[Zone01 Athens local tester](https://github.com/LeKoutz/zone01-checkpoint-imperative-local-tester).

## Requirements

- Go 1.23 or newer
- Docker Desktop or Docker Engine with a running Linux-container daemon

The pinned grader image targets `linux/amd64`. On ARM64 hosts, Docker must
provide amd64 emulation; otherwise the application reports that emulation must
be enabled before a run can start.

That is everything needed to run the application. The compiled frontend is
included in the repository, so students do not need Node.js or npm.

## Quick start

```sh
go run ./cmd/server -open
```

The first run downloads the digest-pinned official grader image, starts the
server, and opens [http://127.0.0.1:8080](http://127.0.0.1:8080). Later runs
reuse that exact image. Keep the terminal open while using the app, and press
`Ctrl+C` there to stop it.

At startup, the application checks GitHub for newer commits and prints the
appropriate `git pull` command when an update is available. The check has a
short timeout, is cached for one hour, and never prevents startup. Disable it
with `-check-updates=false`.

Each submission receives a disposable container with no network, a read-only
root filesystem, resource limits, and read-only bind mounts for the submitted
`main.go` and any exercise-supplied resources. The official entrypoint drops
privileges before executing learner code. This remains intended for local and
classroom use; do not expose
the server as a public code-execution service. See [SECURITY.md](SECURITY.md).

## Using the assessment

- Work through the exercises in order; passing every official check unlocks the
  next exercise.
- Use **Test** for feedback, **gofmt** to format the code, and **.go** to download
  the current editor contents.
- Tests are shown in the fixed order used by the official grader.
  Program output from `z01.PrintRune` or `fmt.Print*` appears in the console.
- Progress and editor layout are saved in the browser on the current device.
- Drag the panel dividers to resize the instructions, editor, tests, and console.

## Development

The React/TypeScript source is already built and embedded under
`internal/web/dist`. Node.js and npm are only needed when changing that source
or running its checks. Install the locked dependencies with:

```sh
npm --prefix web ci
```

GNU Make is optional, but provides shortcuts:

```sh
make check
```

Common targets:

| Target | Purpose |
| --- | --- |
| `make run` | Start with the digest-pinned official Docker grader |
| `make check` | Run Go and frontend verification |
| `make frontend-dev` | Start the Vite development server |
| `make frontend-build` | Rebuild the embedded frontend |
| `make test` | Run the Go test suite |
| `make docker-test` | Run the opt-in Docker integration suite |

For live frontend development, run `make run` and `make frontend-dev` in
separate terminals, then open
[http://127.0.0.1:5173](http://127.0.0.1:5173).

The GitHub Actions workflow runs Go formatting, vetting and tests plus frontend
linting, tests, and build on every pull request.

## License

See [LICENSE](LICENSE).
