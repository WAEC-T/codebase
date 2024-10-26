#!/bin/sh
echo "Hello from Alpine Linux"

echo "Setting up NTP..."
apk add --no-cache openntpd
date

echo "Installing packages..."
apk update

echo "Setting hostname..."
echo "rpi-b-client" > /etc/hostname

apk add --no-cache vim

echo "Installing Python 3.12 and pip..."
apk add --no-cache python3=3.12.7-r0 python3-dev py3-pip

echo "Creating a Python virtual environment..."
python3 -m venv /home/waect/venv

# shellcheck source=/dev/null
. /home/waect/venv/bin/activate

echo "Installing Python packages in the virtual environment..."
pip install requests flask

echo "Provisioning complete!"