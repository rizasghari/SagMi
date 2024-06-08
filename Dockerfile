# Use the latest version of Go
FROM golang:latest

# Set the working directory in the Docker image
WORKDIR /app

# Copy the Go application code into the Docker image
COPY . .

# Install the necessary dependencies for the Go application
RUN go get -d -v ./...

# Build the Go application
RUN go build -o sagmi

# Expose the necessary ports for the application
EXPOSE 8080

# Set the entry point for the Docker container
CMD ["./sagmi"]