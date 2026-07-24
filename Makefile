GO ?= go
GOFMT ?= gofmt
NPM ?= npm
DOCKER ?= docker

RUNNER_IMAGE ?= imperative-go-assessment-runner:manual-check

.DEFAULT_GOAL := help
.PHONY: help run run-local frontend-install frontend-dev frontend-build \
	fmt fmt-check vet test frontend-lint frontend-test check docker-build docker-test
.NOTPARALLEL: check

help:
	@echo "Usage:"
	@echo "  make run               Start with the default Docker sandbox"
	@echo "  make run-local         Start with the trusted local runner"
	@echo "  make check             Run formatting, vet, tests, lint, and frontend build"
	@echo "  make frontend-install  Install locked frontend dependencies"
	@echo "  make frontend-dev      Start the Vite development server"
	@echo "  make docker-test       Run the opt-in Docker integration suite"

run:
	$(GO) run ./cmd/server -runner docker -open

run-local:
	$(GO) run ./cmd/server -runner local -open

frontend-install:
	$(NPM) --prefix web ci

frontend-dev:
	$(NPM) --prefix web run dev

frontend-build:
	$(NPM) --prefix web run build

fmt:
	$(GOFMT) -w cmd internal

fmt-check:
	@files="$$($(GOFMT) -l cmd internal)"; \
	if [ -n "$$files" ]; then \
		echo "$$files"; \
		exit 1; \
	fi

vet:
	$(GO) vet ./...

test:
	$(GO) test -count=1 ./...

frontend-lint:
	$(NPM) --prefix web run lint

frontend-test:
	$(NPM) --prefix web test

check: fmt-check vet test frontend-lint frontend-test frontend-build

docker-build:
	$(DOCKER) build --file docker/runner.Dockerfile --tag $(RUNNER_IMAGE) .

docker-test: export IMPERATIVE_DOCKER_INTEGRATION := 1
docker-test:
	$(GO) test -count=1 -run TestDockerRunnerIntegration -v ./internal/runner
