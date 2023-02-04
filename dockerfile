# Use a lightweight Alpine Linux image as the base image
FROM golang:alpine as builder

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the source code to the container
COPY . .

# Build the Go application
RUN go build main.go 

# Use a new image based on the Alpine Linux image
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /app/

EXPOSE 8080

# Set the working directory to /app
WORKDIR /app

# Specify the command to run when a container is started from the image
CMD ["./main"]


