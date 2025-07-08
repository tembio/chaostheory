#!/bin/bash

set -e

# Get the directory of this script (should be mockeventgenerator)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

echo "[INFO] Building mockeventgenerator"

# Ensure cleanup of the common directory on exit, even if the script fails
trap 'rm -rf "$SCRIPT_DIR/common"' EXIT

# Copy the common directory from parent if it doesn't already exist
if [ -d "$SCRIPT_DIR/common" ]; then
  echo "[INFO] 'common' directory already exists in mockeventgenerator. Removing it first."
  rm -rf "$SCRIPT_DIR/common"
fi

cp -r "$PARENT_DIR/common" "$SCRIPT_DIR/common"
echo "[INFO] Copied 'common' directory into mockeventgenerator."

docker build -t mockeventgenerator "$SCRIPT_DIR"

echo "[INFO] mockeventgenerator built"



echo "[INFO] Building rabbitMQ"

docker build -f Dockerfile.rabbitmq -t rabbitleaderboard .

echo "[INFO] rabbitMQ built"
