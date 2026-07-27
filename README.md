# Imperative Go Practice Assessment

A local, browser-based practice assessment with 171 Go exercises, automatic
tests, console feedback, and saved progress. It starts with foundational
problems and progresses through checkpoint-style challenges. Submitted code
runs in a restricted Docker container by default.

This is an unofficial practice tool and is not affiliated with Zone01.
Exercises after the 21-level foundation adapt the public
[Go checkpoint catalogue](https://github.com/software-sappho/.go-checkpoints-solutions/tree/main)
and [Zone01 piscine catalogue](https://github.com/kinoz01/zone01-Piscine).
Duplicate and alternate-solution files are omitted, and the canonical
[01-edu subjects](https://github.com/01-edu/public/tree/master/subjects) define
the contracts and examples.

## Requirements

- Go 1.23 or newer
- Docker Desktop or Docker Engine with a running Linux-container daemon

That is everything needed to run the application. The compiled frontend is
included in the repository, so students do not need Node.js or npm.

## Quick start

```sh
go run ./cmd/server -open
```

The first run builds the pinned sandbox image, starts the server, and opens
[http://127.0.0.1:8080](http://127.0.0.1:8080). Later runs reuse the image while
its inputs remain unchanged. Keep the terminal open while using the app, and
press `Ctrl+C` there to stop it.

At startup, the application checks GitHub for newer commits and prints the
appropriate `git pull` command when an update is available. The check has a
short timeout, is cached for one hour, and never prevents startup. Disable it
with `-check-updates=false`.

Docker mode gives each submission a disposable container with no network,
a read-only root filesystem, dropped capabilities, resource limits, and no host
directory mounts. It is intended for local and classroom use; do not expose the
server as a public code-execution service. See [SECURITY.md](SECURITY.md) for the
full trust model.

## Optional local runner

If Docker is unavailable, trusted code can run directly on the host:

```sh
go run ./cmd/server -runner local -open
```

The local runner executes submissions with your operating-system permissions.
It keeps the same execution and output limits but is not a sandbox.

## Using the assessment

- Work through the exercises in order; passing every visible test unlocks the
  next exercise.
- Use **Test** for feedback, **gofmt** to format the code, and **.go** to download
  the current editor contents.
- Drag tests or use their arrow buttons to choose the order in which they run.
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
| `make run` | Start with the Docker sandbox |
| `make run-local` | Start with the trusted local runner |
| `make check` | Run Go and frontend verification |
| `make frontend-dev` | Start the Vite development server |
| `make frontend-build` | Rebuild the embedded frontend |
| `make test` | Run the Go test suite |
| `make docker-build` | Build the sandbox image directly |
| `make docker-test` | Run the opt-in Docker integration suite |

For live frontend development, run `make run` and `make frontend-dev` in
separate terminals, then open
[http://127.0.0.1:5173](http://127.0.0.1:5173).

The GitHub Actions workflow runs Go formatting, vetting and tests, frontend
linting, tests and build, and a Docker image build on every pull request.

## License

See [LICENSE](LICENSE).
