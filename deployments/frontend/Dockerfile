# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache build-base git ca-certificates

## Pre-cache modules (copy only module files first)
COPY apps/frontend/go.mod apps/frontend/go.sum ./apps/frontend/
WORKDIR /src/apps/frontend
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

## Copy only frontend sources to shrink context
WORKDIR /src
COPY apps/frontend/ ./apps/frontend/

## Build static binary from apps/frontend/cmd
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src/apps/frontend
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/frontend ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /out/frontend /app/frontend
EXPOSE 8080
ENTRYPOINT ["/app/frontend"]


