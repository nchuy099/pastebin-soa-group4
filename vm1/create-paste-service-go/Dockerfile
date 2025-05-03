# Use a minimal base image
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o main ./main.go

# Final stage: Use a small image to run the app
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/main .

# Expose the application's port
EXPOSE 8081

# Run the application
CMD ["./main"] 