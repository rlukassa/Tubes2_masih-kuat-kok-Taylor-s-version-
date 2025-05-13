# Backend Dockerfile
FROM golang:1.19-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Set environment variables
ENV GIN_MODE=release

# Expose port 8081
EXPOSE 8081

# Command to run the executable
CMD ["./main"]