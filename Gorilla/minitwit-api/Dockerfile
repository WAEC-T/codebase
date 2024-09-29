# Step 1: Build the Go binary
FROM golang:1.23-alpine3.19 AS Build

# Set the working directory inside the container
WORKDIR /usr/src/api

# Copy the Go Modules manifests
COPY *.go go.mod go.sum ./

# Copy the Go Modules manifests
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -v -o /minitwit/api .

# Step 2: Run the Go binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /minitwit
RUN apk add --no-cache bash
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh wait-for-it.sh
RUN chmod +x wait-for-it.sh

# Copy the Go binary from the builder
COPY --from=Build /minitwit/api /minitwit/api

# Expose the port the web server will run on
RUN chown -R 1000 /minitwit
USER 1000:1000

# Run the Go binary
CMD ["/minitwit/app"]