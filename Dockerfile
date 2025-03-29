# Build stage
FROM golang:1.24-alpine AS builder

# Install make and build dependencies
RUN apk add --no-cache make git

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with explicit GOOS and GOARCH
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=0.1.0-alpha" -o bin/kubectl-k8smed ./cmd/kubectl-k8smed

# Final stage
FROM alpine:3.19

# Install kubectl and CA certificates
RUN apk add --no-cache curl ca-certificates && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Copy binary from builder stage
COPY --from=builder /app/bin/kubectl-k8smed /usr/local/bin/kubectl-k8smed

# Set PATH so that kubectl can find the plugin
ENV PATH="/usr/local/bin:${PATH}"

# Create non-root user
RUN addgroup -S k8smed && adduser -S k8smed -G k8smed

# Use non-root user
USER k8smed

# Set default command
ENTRYPOINT ["kubectl-k8smed"]
CMD ["--help"] 