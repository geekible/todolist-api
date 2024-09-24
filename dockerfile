# Use the official Golang image as the base image
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and configuration file
COPY . .

# Build the Go application
RUN go build -o todo-api .

# Use a minimal image for the final build
FROM gcr.io/distroless/base

# Copy the binary and configuration file from the builder stage
COPY --from=builder /app/todo-api /todo-api
COPY --from=builder /app/config.yaml /config.yaml

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/todo-api"]