#!/usr/bin/env bash
set -e

SRC_DIR="build"
DEST_DIR="example/vendor/agai"

if [ ! -d "$SRC_DIR" ]; then
  echo "Error: build directory not found"
  exit 1
fi

mkdir -p "$DEST_DIR"

echo "Copying build artifacts to example/vendor/agai..."

cp -R "$SRC_DIR/"* "$DEST_DIR/"

echo "Done."
