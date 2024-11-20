#!/bin/sh

if [ -d "/media/mmcblk0p2/venv" ]; then
    . "/media/mmcblk0p2/venv/bin/activate"
else
    echo "Virtual environment activate not found at $VENV_PATH"
    exit 1
fi

if [ -f "/media/mmcblk0p2/scenario.py" ]; then
    python "/media/mmcblk0p2/scenario.py" >> /media/mmcblk0p2/scenario.log 2>&1 &
else
    echo "Script not found at $SCRIPT_PATH"
    exit 1
fi

echo "Scenario script is running in the background. Check logs with tail -f /media/mmcblk0p2/scenario.log"
