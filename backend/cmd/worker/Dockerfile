# ── Stage 1: Build the Go worker binary ────────────────────────────────────
FROM golang:1.24.4-alpine AS builder

# Install git and C toolchain for go modules
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go.mod & go.sum (from backend/)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the worker binary
WORKDIR /app/cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_worker

# ── Stage 2: Create a minimal runtime image ─────────────────────────────────
FROM alpine:latest

# Install certificates
RUN apk --no-cache add ca-certificates

# Create non‐root user
RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser

WORKDIR /home/appuser

# Copy the compiled worker binary
COPY --from=builder /quizgenie_worker .

# Optionally copy a .env if you mount it via docker-compose
# COPY --from=builder /app/.env ./

ENTRYPOINT ["./quizgenie_worker"]
