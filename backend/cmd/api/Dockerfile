# ── Stage 1: Build the Go API binary ────────────────────────────────────────
FROM golang:1.24.4-alpine AS builder

# Install git (for go modules) and any C toolchain if needed
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod & go.sum from the context root (which is ./backend)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything else (source code)
COPY . .

# Build the API binary
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_api

# ── Stage 2: Create a minimal runtime image ─────────────────────────────────
FROM alpine:latest

# Install certificates
RUN apk --no-cache add ca-certificates

# Create a non‐root user to run the binary
RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser

WORKDIR /home/appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /quizgenie_api .

# Optionally copy in a .env if you prefer (docker‐compose mounts it)
# COPY --from=builder /app/.env ./

# Expose the port that the API listens on
EXPOSE 8080

ENTRYPOINT ["./quizgenie_api"]
