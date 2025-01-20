# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache build-base

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Runtime
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main", "--port", "8080"]
