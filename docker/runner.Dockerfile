FROM golang:1.23.12-alpine3.22@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f AS builder

WORKDIR /src
COPY go.mod ./
COPY internal/sandboxprotocol ./internal/sandboxprotocol
COPY cmd/sandbox-runner ./cmd/sandbox-runner
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/sandbox-runner ./cmd/sandbox-runner

FROM golang:1.23.12-alpine3.22@sha256:383395b794dffa5b53012a212365d40c8e37109a626ca30d6151c8348d380b5f

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
RUN addgroup -S -g 65532 sandbox \
    && adduser -S -D -H -u 65532 -G sandbox sandbox
COPY --from=builder --chown=65532:65532 /out/sandbox-runner /usr/local/bin/sandbox-runner

USER 65532:65532
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/sandbox-runner"]
