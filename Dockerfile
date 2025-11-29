#### Widget build stage (React)

FROM node:18-alpine AS widget-builder

WORKDIR /widget

# Copy widget source
COPY widget/ .

# Install dependencies and build widget
RUN npm install
RUN npm run build

#### Widget-Lit build stage (Lit)

FROM node:18-alpine AS widget-lit-builder

WORKDIR /widget-lit

# Copy widget-lit source
COPY widget-lit/ .

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

# Create static directories and copy widgets
RUN mkdir -p static/widget static/widget-lit
COPY --from=widget-builder /widget/dist/widget.js static/widget/
COPY --from=widget-lit-builder /widget-lit/dist/widget.js static/widget-lit/

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
