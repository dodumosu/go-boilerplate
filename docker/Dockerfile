# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Set Go environment variables for cross-compilation and static linking (good practice)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
# Download dependencies and verify
RUN go mod download && go mod verify

# Copy the rest of the application source code
COPY . .

# Build the application, creating a static binary
# Corrected the path to the main package
RUN go build -ldflags="-w -s" -a -installsuffix cgo -o /app/boilerplate-app ./cmd/main.go

# Stage 2: Create the runtime image
FROM alpine:3.22
# Using a specific stable version of Alpine is good practice, but latest is also fine for dev

RUN apk --no-cache add curl

WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the necessary files from the builder stage to the runtime stage
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations/
COPY --from=builder /app/boilerplate-app .
COPY --from=builder /app/entrypoint.sh .

# Ensure the binary and copied files are owned by the non-root user and the binary/script are executable.
RUN chown -R appuser:appgroup /app && chmod +x /app/boilerplate-app /app/entrypoint.sh

# Switch to the non-root user
USER appuser

# Expose the port the application runs on (default 8080)
# This should be configurable via API_PORT in docker-compose
ARG SERVER_PORT=8080
EXPOSE ${SERVER_PORT}

# Command to run the application
# The actual port binding will be handled by docker-compose
ENTRYPOINT [ "./entrypoint.sh" ]
