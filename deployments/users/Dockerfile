# syntax=docker/dockerfile:1

ARG GO_VERSION=1.23.5

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache build-base git ca-certificates

## Pre-cache modules (copy only module files first)
COPY apps/users/go.mod apps/users/go.sum ./apps/users/
WORKDIR /src/apps/users
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

## Copy only users sources to shrink context
WORKDIR /src
COPY apps/users/ ./apps/users/

## Build static binary from apps/users/cmd
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src/apps/users
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/users ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /out/users /app/users
EXPOSE 8080
ENTRYPOINT ["/app/users"]
