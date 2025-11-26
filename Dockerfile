# syntax=docker/dockerfile:1

### 1. Build Stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of the project (including .env, config/, app/, main.go, etc.)
COPY . .

# Build static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main ./main.go


### 2. Runtime Stage
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy compiled binary
COPY --from=builder /app/main .

# Copy .env from host → container
COPY .env .env

# Copy docs directory for Swagger
COPY --from=builder /app/docs ./docs

# Create a non-root user
RUN adduser -D appuser
USER appuser

EXPOSE 8080

CMD ["./main"]
