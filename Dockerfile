# Use the official Golang image as the base image
FROM golang:1.22.1-bookworm

# Set the working directory to /app
WORKDIR /app

# Install necessary dependencies
RUN apt-get update && \
    apt-get install -y curl wget gnupg2 ca-certificates

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o mortgage-rate-exporter .

# Expose port 8080 for the API server
EXPOSE 8080

CMD ["./mortgage-rate-exporter"]
