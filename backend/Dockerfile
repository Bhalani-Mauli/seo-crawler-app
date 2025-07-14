# --- Build Stage ---
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Install git for go mod
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o seo-crawler ./cmd/server/main.go

# --- Run Stage ---
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/seo-crawler .

# Expose the default port
EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV DB_HOST=127.0.0.1 \
    DB_PORT=3306 \
    DB_USER=seo_user \
    DB_PASSWORD=seo_password \
    DB_NAME=seo_crawler \
    SERVER_PORT=8080 \
    API_KEY=seo-crawler-api-key-2025

# Run the application
ENTRYPOINT ["/app/seo-crawler"] 