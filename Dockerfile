#### Widget build stage
FROM node:20-alpine AS widget-builder

ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL

WORKDIR /widget
COPY widget/ .
RUN npm install
RUN npm run build

#### Backend assets build stage (Vite CSS/JS)
FROM node:20-alpine AS backend-assets-builder

WORKDIR /backend
COPY backend/package.json backend/package-lock.json* ./
RUN npm install
COPY backend/vite.config.js ./
COPY backend/assets/ ./assets/
RUN npm run build

#### Backend build stage
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Add C dependencies needed for go-webp
RUN apk add --no-cache gcc musl-dev libwebp-dev

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ .

# Copy built assets from asset builder
COPY --from=backend-assets-builder /backend/static/dist/ ./static/dist/

# Copy widget
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

# Run the binary
CMD ["./main"]
