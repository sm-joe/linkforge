# syntax=docker/dockerfile:1.7

# ---------- Build stage ----------
FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.work ./
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY packages/url-validation/go.mod ./packages/url-validation/

RUN cd apps/api && go mod download

COPY apps/api ./apps/api
COPY packages/url-validation ./packages/url-validation

WORKDIR /src/apps/api

RUN GOWORK=/src/go.work \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/linkforge-api \
    ./cmd/server

# ---------- Runtime stage ----------
FROM alpine:3.22

RUN addgroup -S linkforge \
    && adduser -S -G linkforge linkforge \
    && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/linkforge-api /app/linkforge-api

RUN mkdir -p /app/data \
    && chown -R linkforge:linkforge /app

USER linkforge

ENV PORT=8080 \
    DATABASE_PATH=/app/data/linkforge.db

EXPOSE 8080

ENTRYPOINT ["/app/linkforge-api"]