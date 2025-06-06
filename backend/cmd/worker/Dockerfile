# ── Stage 1: Build the Go worker binary ────────────────────────────────────
FROM golang:1.24.4-alpine AS builder

# Install git and C toolchain
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# 1) Copy go.mod & go.sum from ./backend
COPY go.mod go.sum ./
RUN go mod download

# 2) Copy all backend source code
COPY . .

# 3) Build the worker binary
WORKDIR /app/cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_worker


# ── Stage 2: Create a minimal runtime image ─────────────────────────────────
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /home/appuser

# 1) Copy the worker binary
COPY --from=builder /quizgenie_worker           .

# 2) Copy the “wait-for-postgres” script
COPY cmd/worker/wait-for-postgres.sh /usr/local/bin/wait-for-postgres.sh

# 3) Make them executable
RUN chmod +x ./quizgenie_worker /usr/local/bin/wait-for-postgres.sh

# 4) Create a non-root user and switch to it
RUN addgroup -S appuser && adduser -S -G appuser appuser
USER appuser

ENTRYPOINT ["/usr/local/bin/wait-for-postgres.sh", "./quizgenie_worker"]
