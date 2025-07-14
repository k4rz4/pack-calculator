# Multi-stage build for production optimization
FROM golang:1.22-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS calls and wget for health checks
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy web assets from builder stage
COPY --from=builder /app/web ./web

# Change ownership to non-root user
RUN chown -R appuser:appuser . 
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
