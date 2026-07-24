FROM golang:1.26.5-alpine3.24@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

WORKDIR /src
COPY go.mod ./
COPY internal/sandboxprotocol ./internal/sandboxprotocol
COPY cmd/sandbox-runner ./cmd/sandbox-runner
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sandbox-runner ./cmd/sandbox-runner

FROM golang:1.26.5-alpine3.24@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS toolchain

RUN mkdir -p /opt/go-build \
    && GOCACHE=/opt/go-build CGO_ENABLED=0 GOMAXPROCS=1 go install -p=1 -trimpath \
    bufio \
    context \
    encoding/json \
    errors \
    fmt \
    io \
    net/http \
    net/http/httptest \
    sort \
    strconv \
    strings \
    sync \
    sync/atomic \
    time \
    && chmod -R a+rX /opt/go-build

FROM scratch

ENV PATH=/usr/local/go/bin
COPY --from=toolchain --chown=65532:65532 /usr/local/go /usr/local/go
COPY --from=toolchain --chown=65532:65532 /opt/go-build /opt/go-build
COPY --from=builder --chown=65532:65532 /out/sandbox-runner /sandbox-runner

USER 65532:65532
WORKDIR /workspace
ENTRYPOINT ["/sandbox-runner"]
