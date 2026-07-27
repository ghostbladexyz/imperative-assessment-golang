FROM golang:1.26.5-alpine3.24@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS builder

WORKDIR /src
COPY go.mod ./
COPY third_party/z01 ./third_party/z01
COPY internal/sandboxprotocol ./internal/sandboxprotocol
COPY cmd/sandbox-runner ./cmd/sandbox-runner
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sandbox-runner ./cmd/sandbox-runner

FROM golang:1.26.5-alpine3.24@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS toolchain

WORKDIR /src
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
    && chmod -R a+rX /opt/go-build \
    && rm -rf \
        /usr/local/go/api \
        /usr/local/go/doc \
        /usr/local/go/misc \
        /usr/local/go/test \
        /usr/local/go/src/cmd \
        /usr/local/go/bin/gofmt \
    && find /usr/local/go/src -type f -name '*_test.go' -delete \
    && find /usr/local/go/src -type d -name testdata -prune -exec rm -rf '{}' +

FROM scratch

ENV PATH=/usr/local/go/bin \
    GOPROXY=off
COPY --from=toolchain --chown=65532:65532 /usr/local/go /usr/local/go
COPY --from=toolchain --chown=65532:65532 /opt/go-build /opt/go-build
COPY --chown=65532:65532 third_party/z01 /opt/z01
COPY --from=builder --chown=65532:65532 /out/sandbox-runner /sandbox-runner

USER 65532:65532
WORKDIR /workspace
ENTRYPOINT ["/sandbox-runner"]
