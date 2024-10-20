#!/bin/sh

ENV_FILE="../../.env.production"
REMOTE_USER="username_rbpi_server"
REMOTE_HOST="hostname_rbpi_server"
REMOTE_PATH="/path/of/our/root/project/rbpisv"

check_command() {
    if [ $? -ne 0 ]; then
        echo "Error: $1 failed."
        exit 1
    fi
}

echo "Starting provisioning for Docker on Alpine Linux..."

apk update
check_command "apk update"

echo "Installing Docker..."
apk add --no-cache docker
check_command "Docker installation"

echo "Starting Docker service..."
rc-update add docker default
service docker start
check_command "Starting Docker service"

echo "Installing Docker Compose..."
apk add --no-cache docker-compose
check_command "Docker Compose installation"

chmod +x /usr/local/bin/docker-compose
check_command "Making Docker Compose executable"

DOCKER_IMAGES="waect/rust-actix waect/go-gorilla waect/python-flask"

for IMAGE in $DOCKER_IMAGES; do
    echo "Pulling Docker image: $IMAGE"
    docker pull "$IMAGE"
    check_command "Pulling $IMAGE"
done

echo "Copying .env.production file to remote server..."
scp "$ENV_FILE" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH"
check_command "Copying .env.production file"

echo "Provisioning complete..."