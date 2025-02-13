# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.17

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy the env file
COPY ../.env .

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]
