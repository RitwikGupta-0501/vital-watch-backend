# ---- Build Stage ----
# Use the official Go image as the builder
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# 1. Copy *only* the module files
COPY go.mod go.sum ./

# 2. Download dependencies
# This layer is cached as long as go.mod/go.sum don't change
RUN go mod download

# 3. Copy *only* the Go source code
# This is the key optimization. We only copy folders
# that contain Go code needed for the build.
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY utils/ ./utils/

# 4. Build the Go application
# This layer is now *only* invalidated if Go files (above) change.
# It will no longer be invalidated by changes to .sql, .env, etc.
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main/main.go

# 5. Copy other assets *after* the build.
# We copy them to the builder so the final stage can get them.
COPY ./migrations ./migrations

# ---- Final Stage ----
# Use a tiny, clean Alpine image for the final container
FROM alpine:latest

WORKDIR /app

# Copy *only* the built binary from the builder stage
COPY --from=builder /app/main .

# Copy your migrations from the builder stage
COPY --from=builder /app/migrations ./migrations

# EXPOSE the port your Gin app runs on (default is 8080)
EXPOSE 8080

# This is the command that will run when the container starts
ENTRYPOINT ["/app/main"]