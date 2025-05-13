# Start from the official Go image as a build stage
FROM golang:1.24.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod ./
COPY go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Use a minimal base image for the final stage
FROM alpine:3.21.3

# Set the working directory inside the container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /app/server .

# Copy the static files to serve
COPY --from=builder /app/cmd/server/static /static

# Copy env files
COPY --from=builder /app/cmd/server/.env ./.env
COPY --from=builder /app/cmd/server/.env.development ./.env.development
COPY --from=builder /app/cmd/server/.env.development ./.env.production

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./server"]