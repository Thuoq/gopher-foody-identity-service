# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/server/main.go

# Development stage with Hot Reload
FROM golang:1.26-alpine AS dev

WORKDIR /app

# Install air and tzdata
RUN go install github.com/air-verse/air@latest && \
    apk add --no-cache tzdata

# Copy go mod and sum files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

CMD ["air", "-c", ".air.toml"]

# Final stage
FROM alpine:latest

# Install tzdata for timezone support
RUN apk add --no-cache tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 8001 9091

CMD ["./main"]
