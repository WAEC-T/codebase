#!/bin/sh

echo "Hello bit..es from Alpine Linux"

echo "Setting up NTP..."
apk add --no-cache openntpd
ntpd -c /etc/ntpd.conf

echo "Installing packages..."
apk update

echo "Setting hostname..."
echo "rpi-b-client" > /etc/hostname

apk add --no-cache vim

echo "Setting up BASH as default shell..."
apk add --no-cache bash
chsh -s /bin/bash

echo "Installing Python 3.12 and pip..."
apk add --no-cache python3=3.12.7-r0 python3-dev py3-pip
pip3 install --upgrade pip

pip install requests
pip install flask

echo "Provisioning complete!"
