# Security Policy

## Supported use

Imperative Checkpoint Practice Assessment is intended for local development and supervised classroom use. Its Docker runner materially improves isolation, but the project is not a hardened anonymous public code-execution service.

Security fixes are applied to the latest version on the default branch.

## Reporting a vulnerability

Please report vulnerabilities through GitHub's private security advisory feature for this repository. Include:

- The affected commit or version
- The runner mode and host platform
- Reproduction steps
- The expected and observed isolation boundary
- Any evidence of host, credential, network, or cross-submission access

Do not include real credentials, private keys, or unrelated personal data in a report. Avoid publishing a working escape before a fix is available.

## Trust boundaries

The server binds to `127.0.0.1` by default and selects the runner at startup. Browser requests cannot change that selection.

Each submission receives a new container with no network or shared IPC, a read-only root filesystem, resource limits, a bounded tmpfs, and no Docker log persistence. The temporary submitted `main.go` and declared exercise resources are bind-mounted read-only. The pinned official entrypoint compiles them and then uses `setpriv` to execute learner code as UID/GID 65534 with cleared supplementary groups and `no-new-privileges`. The container is force-removed after every outcome, and startup removes labeled non-running containers left by an abrupt host-process exit.

The outer container starts as root because the official entrypoint needs to compile and then perform that privilege transition. Applying Docker's blanket capability drop or outer `no-new-privileges` prevents the pinned grader from switching users. This is a deliberate compatibility trade-off, not a claim of hardened multi-tenant isolation.

The following remain trusted:

- The Docker daemon and Docker Desktop or Engine installation
- Docker's VM, host kernel, container runtime, and configuration
- The digest-pinned official grader image, including its tests and compiled oracles
- The assessment server process and result parser
- The local receipt signing key

The official image reference includes an immutable SHA-256 digest. It contains the grader implementation and compiled oracles by design, but receives no host credentials, receipt keys, Docker CLI, or Docker socket.

## Known limitations

- A container escape or Docker/runtime vulnerability may compromise the host boundary.
- Resource limits reduce denial-of-service risk but do not eliminate all host pressure or runtime bugs.
- The Docker daemon is privileged infrastructure; anyone who can control it is outside this threat model.
- Binding beyond localhost or forwarding the port expands exposure and is unsupported.
- This design does not provide multi-tenant identity, authentication, quotas, abuse detection, or workload separation.

A public service would require stronger isolation such as a fresh microVM per submission, dedicated execution infrastructure, authentication, rate limiting, monitoring, incident response, and an independent security review.

## Operational checks

Keep Docker Desktop or Engine and the host OS patched. Review base-image digest updates before merging them. After integration testing, confirm no assessment containers remain:

```sh
docker ps -a --filter name=imperative-go-assessment- --format "{{.Names}} {{.Status}}"
```

Remove only the exact digest-pinned grader reference when cleaning the cache. Do not use broad prune commands as part of this application's normal operation.
