# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the skiff binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /skiff ./cmd/skiff

# Final stage
FROM alpine:latest

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /home/appuser

# Copy the skiff binary from the builder stage
COPY --from=builder /skiff /usr/local/bin/skiff

# Set the user
USER appuser

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/skiff"]
