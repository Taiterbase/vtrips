# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache build-base git ca-certificates

# Pre-cache modules
COPY ../../apps/frontend/go.mod ../../apps/frontend/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy full repo
COPY ../../ ./

# Build static binary from apps/frontend/cmd
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/frontend ./apps/frontend/cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /out/frontend /app/frontend
EXPOSE 8080
ENTRYPOINT ["/app/frontend"]


