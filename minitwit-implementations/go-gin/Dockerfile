# Step 1: Build the Go binary
FROM alpine:edge as BUILDER

# Set the working directory inside the container
WORKDIR /app

# Install Go
RUN apk add --no-cache go

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -v -o /minitwit/app go-gin/src

# Step 2: Run the Go binary
FROM alpine:edge

# Set the working directory inside the container
WORKDIR /minitwit

RUN apk add --no-cache ca-certificates

# Copy the Go binary and required files from the builder stage
COPY --from=BUILDER /minitwit/app /minitwit/app
COPY ./templates/ /minitwit/templates/
COPY ./static/ /minitwit/static/

# Expose the port the web server will run on
EXPOSE 5000

# Run the Go binary
CMD ["/minitwit/app"]