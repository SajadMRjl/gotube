# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /gotube ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /gotube /app/gotube
COPY configs/config.yaml /app/configs/config.yaml

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Expose port (if using webhooks)
EXPOSE 8080

# Set the entry point
CMD ["/app/gotube"]