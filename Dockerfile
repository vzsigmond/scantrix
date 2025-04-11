# --- Base stage (Go)
FROM golang:1.24-bookworm AS base
# Install common utilities using apt.
RUN apt-get update && apt-get install -y \
    bash \
    git \
    inotify-tools \
    wget && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app

# --- Dev stage
FROM base AS dev
CMD ["bash"]

# --- Build stage
FROM base AS builder
COPY . .
RUN go build -o scantrix cmd/scantrix/main.go

# --- Final binary stage
FROM debian:bookworm-slim AS prod
RUN apt-get update && apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /root/
COPY --from=builder /app/scantrix .
CMD ["./scantrix"]
