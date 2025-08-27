# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# -ldflags="-w -s" reduces the size of the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o cepm-backend ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/cepm-backend .

# Copy the configuration file template
# In a real setup, you'd mount a config file, but we include the template for reference
COPY --from=builder /app/config/config.yaml.template ./config/config.yaml.template

# Expose port 8080
EXPOSE 8080

# Command to run the application
# The application will look for a config file at ./config/config.yaml
CMD ["./cepm-backend"]
