# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
FROM golang:1.27-alpine AS builder

WORKDIR /src

# Copy only dependency metadata first so dependency downloads stay cached
COPY go.work ./
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY packages/url-validation/go.mod ./packages/url-validation/

RUN --mount=type=cache,target=/go/pkg/mod \
    cd apps/api && go mod download

# Copy source after dependencies
COPY apps/api ./apps/api
COPY packages/url-validation ./packages/url-validation

WORKDIR /src/apps/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOWORK=/src/go.work \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/linkforge-api \
    ./cmd/server

# ---------- Runtime stage ----------
FROM alpine:3.22

RUN addgroup -S linkforge \
    && adduser -S -G linkforge linkforge \
    && apk add --no-cache ca-certificates \
    && mkdir -p /app/data \
    && chown linkforge:linkforge /app/data

WORKDIR /app

COPY --from=builder --chown=linkforge:linkforge /out/linkforge-api /app/linkforge-api

USER linkforge

ENV PORT=8080 \
    DATABASE_PATH=/app/data/linkforge.db

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/healthz || exit 1

ENTRYPOINT ["/app/linkforge-api"]