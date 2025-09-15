# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.5

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache build-base git ca-certificates

## Pre-cache modules (copy only module files first)
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
WORKDIR /src/apps/backend
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

## Copy only backend sources to shrink context
WORKDIR /src
COPY apps/backend/ ./apps/backend/

## Build static binary from apps/backend/cmd
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src/apps/backend
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/backend ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /out/backend /app/backend
EXPOSE 8080
ENTRYPOINT ["/app/backend"]


