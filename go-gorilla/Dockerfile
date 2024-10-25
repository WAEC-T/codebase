# Step 1: Build the Go binary
FROM golang:1.23-alpine3.19 AS Build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY *.go go.mod go.sum ./

# Copy the Go Modules manifests
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -v -o /minitwit/app go-gorilla/src

# Step 2: Run the Go binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /minitwit
RUN apk add --no-cache bash

# Copy the Go binary from the builder
COPY  ./templates/ /minitwit/templates/
COPY  ./static/ /minitwit/static/
COPY --from=Build /minitwit/app /minitwit/app 

# Expose the port the web server will run on
EXPOSE 5000

# Run the Go binary
CMD ["/minitwit/app"]                  
