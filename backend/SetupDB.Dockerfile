# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ./cmd/setup-db/main.go .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./setupdb .

# Final stage
FROM gcr.io/distroless/base-debian12 AS final

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/setupdb .

# Command to run the executable
ENTRYPOINT ["./setupdb"]