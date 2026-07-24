# Security Policy

## Supported use

Imperative Go Practice Assessment is intended for local development and supervised classroom use. The default Docker runner materially improves isolation over direct local execution, but the project is not a hardened anonymous public code-execution service.

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

In Docker mode, each submission receives a new non-root container with no network, a read-only root filesystem, dropped capabilities, `no-new-privileges`, resource limits, bounded tmpfs mounts, and no host bind mounts. The submitted source and server-generated harness travel through a bounded JSON request on standard input. The exact container is force-removed after every outcome.

The following remain trusted:

- The Docker daemon and Docker Desktop or Engine installation
- Docker's VM, host kernel, container runtime, and configuration
- The assessment server process and its server-owned harness
- The local receipt signing key

The runner image deliberately contains a pinned Go toolchain, fixed entrypoint, and precompiled standard-library cache. It must not contain assessment solutions, host credentials, receipt keys, the Docker CLI, or the Docker socket.

## Known limitations

- A container escape or Docker/runtime vulnerability may compromise the host boundary.
- Resource limits reduce denial-of-service risk but do not eliminate all host pressure or runtime bugs.
- The Docker daemon is privileged infrastructure; anyone who can control it is outside this threat model.
- Local runner mode executes submitted code with the current user's permissions and provides no OS isolation.
- Binding beyond localhost or forwarding the port expands exposure and is unsupported.
- This design does not provide multi-tenant identity, authentication, quotas, abuse detection, or workload separation.

A public service would require stronger isolation such as a fresh microVM per submission, dedicated execution infrastructure, authentication, rate limiting, monitoring, incident response, and an independent security review.

## Operational checks

Keep Docker Desktop or Engine and the host OS patched. Review base-image digest updates before merging them. After integration testing, confirm no assessment containers remain:

```sh
docker ps -a --filter name=imperative-go-assessment- --format "{{.Names}} {{.Status}}"
```

Remove only exact project image tags when cleaning the cache. Do not use broad prune commands as part of this application's normal operation.
