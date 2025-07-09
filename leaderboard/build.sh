#!/bin/bash

set -e

# Get the directory of this script (should be leaderboard)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

# Ensure cleanup of the common directory on exit, even if the script fails
trap 'rm -rf "$SCRIPT_DIR/common"' EXIT

# Copy the common directory from parent if it doesn't already exist
if [ -d "$SCRIPT_DIR/common" ]; then
  echo "[INFO] 'common' directory already exists in leaderboard. Removing it first."
  rm -rf "$SCRIPT_DIR/common"
fi

cp -r "$PARENT_DIR/common" "$SCRIPT_DIR/common"
echo "[INFO] Copied 'common' directory into leaderboard."

docker build -t leaderboard "$SCRIPT_DIR"

echo "[INFO] Build complete."