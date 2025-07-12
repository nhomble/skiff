# Builder stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /skiff ./cmd/skiff

FROM gcr.io/distroless/static:nonroot
WORKDIR /home/skiff
COPY --from=builder /skiff /usr/local/bin/skiff

ENTRYPOINT ["/usr/local/bin/skiff"]
