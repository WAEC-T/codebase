#!/bin/sh
echo "Hello from Alpine Linux"

echo "Installing packages..."
apk update

echo "Setting hostname..."
echo "rpi-b-client" > /etc/hostname

echo "Installing Python 3.12 and pip..."
apk add --no-cache python3 py3-pip

echo "Creating a Python virtual environment..."
python3 -m venv ./venv

# shellcheck source=/dev/null
. venv/bin/activate

echo "Installing Python packages in the virtual environment..."
pip install requests flask

echo "Provisioning complete!"