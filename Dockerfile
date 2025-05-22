# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy up dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Install motion apk
RUN apk add --no-cache motion

# Copy motion.conf into image
COPY motion.conf /etc/motion/motion.conf

# Start motion
RUN motion -b

# Expose port 8082
EXPOSE 8082

# Run the application
CMD ["./main"] 