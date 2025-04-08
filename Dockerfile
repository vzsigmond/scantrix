# --- Base stage ---
FROM golang:1.24-alpine AS base
RUN apk add --no-cache bash git inotify-tools

WORKDIR /app

# --- Dev stage ---
FROM base AS dev
CMD ["bash"]

# --- Build stage ---
FROM base AS builder
COPY . .
RUN go build -o scantrix cmd/scantrix/main.go

# --- Binary stage ---
FROM alpine:latest AS prod
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/scantrix .
CMD ["./scantrix"]
