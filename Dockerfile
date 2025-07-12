# Build stage for Go backend
FROM golang:1.24-alpine AS go-builder

# Install build dependencies
# RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o playlists .

# Build stage for frontend
FROM oven/bun:1 AS frontend-builder

# Set working directory
WORKDIR /app

# Copy frontend package files
COPY frontend/package.json frontend/bun.lock ./

# Install dependencies
RUN bun install

# Copy frontend source code
COPY frontend/ ./

# Build the frontend
RUN bun run build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the built Go binary from go-builder stage
COPY --from=go-builder /app/playlists .

# Copy the built frontend from frontend-builder stage
COPY --from=frontend-builder /app/dist ./frontend/dist

# Copy migrations directory
COPY --from=go-builder /app/migrations ./migrations

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./playlists"] 