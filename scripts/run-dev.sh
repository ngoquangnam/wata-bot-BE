#!/bin/bash
# Bash script to run in development environment
# Usage: ./scripts/run-dev.sh

echo "Starting Wata Bot Backend in DEV mode..."

# Load dev environment variables
if [ -f .env.dev ]; then
    export $(cat .env.dev | grep -v '^#' | xargs)
    echo "Loaded .env.dev file"
else
    echo "Warning: .env.dev file not found, using config file only"
fi

# Run the application with dev config
go run wata-bot.go -f etc/wata-bot-api.dev.yaml

