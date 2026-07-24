# Imperative Go Practice Assessment

A local, browser-based practice assessment with 30 Go exercises, automatic
tests, console feedback, and saved progress. Submitted code runs in a restricted
Docker container by default.

The exercises are original practice material, not copied from an official
Zone01 assessment.

## Requirements

- Go 1.23 or newer
- Docker Desktop or Docker Engine with a running Linux-container daemon

GNU Make is recommended for the command shortcuts below. Node.js 20.19+ or
22.12+ and npm are only required for frontend development.

## Run with Docker

```sh
make run
```

The first run builds the pinned sandbox image, starts the server, and opens
[http://127.0.0.1:8080](http://127.0.0.1:8080). Later runs reuse the image while
its inputs remain unchanged.

Without Make:

```sh
go run ./cmd/server -runner docker -open
```

Docker mode gives each submission a disposable container with no network,
a read-only root filesystem, dropped capabilities, resource limits, and no host
directory mounts. It is intended for local and classroom use; do not expose the
server as a public code-execution service. See [SECURITY.md](SECURITY.md) for the
full trust model.

## Optional local runner

For trusted code only:

```sh
make run-local
```

The local runner executes submissions with your operating-system permissions.
It keeps the same execution and output limits but is not a sandbox.

## Development

Install the locked frontend dependencies once:

```sh
make frontend-install
```

Run the full verification suite:

```sh
make check
```

Common targets:

| Target | Purpose |
| --- | --- |
| `make run` | Start with the Docker sandbox |
| `make run-local` | Start with the trusted local runner |
| `make frontend-dev` | Start the Vite development server |
| `make frontend-build` | Rebuild the embedded frontend |
| `make test` | Run the Go test suite |
| `make docker-build` | Build the sandbox image directly |
| `make docker-test` | Run the opt-in Docker integration suite |

For live frontend development, run `make run` and `make frontend-dev` in
separate terminals, then open
[http://127.0.0.1:5173](http://127.0.0.1:5173).

## License

See [LICENSE](LICENSE).
