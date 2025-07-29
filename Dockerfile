# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates tzdata

# Create app directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls and clean up cache to reduce image size
RUN apk --no-cache add ca-certificates curl && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Create app directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy blog data file
COPY --from=builder /app/blog_data.json .

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV LOG_LEVEL=info

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"] 