# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies (go.sum will be created if needed)
RUN go mod download

# Copy source code
COPY main.go main_test.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o palindromic-fuel main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/palindromic-fuel .

# Expose port for web server
EXPOSE 8080

# Set default port environment variable
ENV PORT=8080

# Run the application
CMD ["./palindromic-fuel", "-web"]