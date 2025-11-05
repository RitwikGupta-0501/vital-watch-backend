# ---- Build Stage ----
# Use the official Go image as the builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy all your project files
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
# This creates a static binary at /app/main
# It uses the path from your main.go file
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/main.go

# ---- Final Stage ----
# Use a tiny, clean Alpine image for the final container
FROM alpine:latest

WORKDIR /app

# Copy *only* the built binary from the builder stage
COPY --from=builder /app/main .

# Copy your migrations so the app can find them
# Your main.go file looks for migrations at "file://./migrations"
COPY ./migrations ./migrations

# EXPOSE the port your Gin app runs on (default is 8080)
EXPOSE 8080

# This is the command that will run when the container starts
ENTRYPOINT ["/app/main"]