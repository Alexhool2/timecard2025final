# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with all dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install minimal dependencies
RUN apk add --no-cache ca-certificates tzdata bash  # Add bash here

# Copy binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Create directory for any potential file storage
RUN mkdir -p /app/data

# Set permissions
RUN chmod +x /app/main

# Expose application port
EXPOSE 8081

# Run command
CMD ["./main"]