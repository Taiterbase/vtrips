# syntax=docker/dockerfile:1

ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

# Install build deps
RUN apk add --no-cache build-base git ca-certificates

# Pre-cache modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy source
COPY . .

# Build static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /out/backend ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
USER nonroot:nonroot
COPY --from=builder /out/backend /app/backend
EXPOSE 8080
ENTRYPOINT ["/app/backend"]


