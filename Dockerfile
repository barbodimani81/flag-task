# Dockerfile
FROM golang:1.24.2

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o server .

# Expose the app port
EXPOSE 8080

# Start the app
CMD ["./server"]
