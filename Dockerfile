# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -o main .
RUN ls -la

# Runtime stage
FROM alpine:3.18

# Set working directory
WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Copy binary from build stage
COPY --from=builder /app/main /app/main
RUN chmod +x /app/main

# Copy configuration files
COPY --from=builder /app/template.env .env

# Create directory for SQLite database (if used)
RUN mkdir -p /app/data

# Set default environment variables
ENV GIN_MODE=release \
    DATABASE_TYPE=postgres \
    DATABASE_USER=gin-template \
    DATABASE_PASSWORD=xxxxxxxxxx \
    DATABASE_DBNAME=gin-template

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
