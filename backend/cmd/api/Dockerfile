# ── Stage 1: Build the Go API binary ────────────────────────────────────────
FROM golang:1.24.4-alpine AS builder

# Install git (for go modules) and any C toolchain
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# 1) Copy go.mod & go.sum from ./backend
COPY go.mod go.sum ./
RUN go mod download

# 2) Copy all backend source code
COPY . .

# 3) Build the API binary
WORKDIR /app/cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /quizgenie_api


# ── Stage 2: Create a minimal runtime image ─────────────────────────────────
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /home/appuser

# 1) Copy the API binary
COPY --from=builder /quizgenie_api           .

# 2) Copy the “wait-for-postgres” script from the build context into /usr/local/bin
COPY cmd/api/wait-for-postgres.sh /usr/local/bin/wait-for-postgres.sh

# 3) Make both of them executable
RUN chmod +x ./quizgenie_api /usr/local/bin/wait-for-postgres.sh

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/wait-for-postgres.sh", "./quizgenie_api"]
