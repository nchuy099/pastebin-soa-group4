FROM golang:1.21-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy all files
COPY . .

# Build the worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /view-update-worker .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /view-update-worker .

# Run the worker
CMD ["/app/view-update-worker"] 