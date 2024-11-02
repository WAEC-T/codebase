#!/bin/sh
check_command() {
    if [ $? -ne 0 ]; then
        echo "Error: $1 failed."
        exit 1
    fi
}

echo "Starting provisioning for Docker on Alpine Linux..."

apk update
check_command "apk update"

echo "Installing Docker and openrc (manager)..."
apk add --no-cache openrc docker
check_command "Docker installation"

echo "Installing Docker Compose..."
apk add --no-cache docker-compose
check_command "Docker Compose installation"

echo "Enabling and starting Docker service ..."
rc-update add docker default # or boot here...
check_command "Install and set service"

echo "Provisioning complete. Docker and Docker Compose are installed."

echo "Starting Docker..."
service docker start # if this guy don't work then we have a problem

echo "Pulling the images!"

DOCKER_IMAGES="simonharwick97822/rust-actix:latest simonharwick97822/python-flask:latest" #list the images here

for IMAGE in $DOCKER_IMAGES; do
    echo "Pulling Docker image: $IMAGE"
    docker pull "$IMAGE"
    check_command "Pulling $IMAGE"
done

echo "DONE!"