FROM golang:1.22.3-alpine3.20 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files (go.mod and go.sum)
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp .

# Stage 2: Create a smaller runtime image
FROM golang:1.22.3-alpine3.20

# Install necessary dependencies for the Go app (e.g., libc for running Go binaries)
RUN apk --no-cache add ca-certificates

# Set the current working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/myapp .

# Expose port 8080
EXPOSE 8080

# Command to run the Go application
CMD ["./myapp"]