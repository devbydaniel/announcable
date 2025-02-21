#### Widget build stage

FROM node:18-alpine AS widget-builder

WORKDIR /widget

# Copy widget source
COPY widget/ .

# Install dependencies and build widget
RUN npm install
RUN npm run build

#### Backend build stage

FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Add C dependencies needed for go-webp
RUN apk add --no-cache gcc musl-dev libwebp-dev

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ .

# Create static directory and copy widget
RUN mkdir -p static/widget
COPY --from=widget-builder /widget/dist/widget.js static/widget/

# Build the application
RUN GOOS=linux go build -o main .

#### Final stage

FROM alpine:3.17

# Dependency of the go-webp library
RUN apk add --no-cache libwebp-dev

WORKDIR /app

# Copy the binary from builder
COPY --from=backend-builder /app/main .

# Copy env
COPY .env .

# Run the binary
CMD ["./main"]
