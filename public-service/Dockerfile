FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum, then download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Run go mod tidy to update go.sum if needed.
RUN go mod tidy

# Copy the rest of the source code.
COPY . .

# Build the binary.
RUN go build -o public-service

EXPOSE 3001

CMD ["./public-service"]
