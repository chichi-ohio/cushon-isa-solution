# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY config/config.yaml ./config/
COPY templates/ ./templates/
COPY static/ ./static/

# Create non-root user
RUN adduser -D appuser
USER appuser

# Command to run the application
CMD ["./main"]

EXPOSE 8081 