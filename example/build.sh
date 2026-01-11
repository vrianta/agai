#!/usr/bin/env bash
set -e

CXX=g++
CXXFLAGS="-std=c++17 -O2"
VENDOR_DIR="vendor/agai"
INC_DIR="$VENDOR_DIR/include"
LIB_DIR="$VENDOR_DIR/lib"

BUILD_DIR="build"
BIN_DIR="$BUILD_DIR/bin"
LIB_OUT="$BUILD_DIR/lib"

mkdir -p "$BIN_DIR" "$LIB_OUT"

# Find all cpp files recursively
SRC_FILES=$(find src -type f -name "*.cpp")

# Auto-detect all vendor .so libraries
LIB_FLAGS=""
for so in "$LIB_DIR"/*.so; do
  name=$(basename "$so" | sed 's/^lib//;s/\.so$//')
  LIB_FLAGS="$LIB_FLAGS -l$name"
done

echo "Compiling sources:"
echo "$SRC_FILES"

$CXX $CXXFLAGS \
  $SRC_FILES \
  -I"$INC_DIR" \
  -L"$LIB_DIR" \
  $LIB_FLAGS \
  -Wl,-rpath,'$ORIGIN/../lib' \
  -o "$BIN_DIR/app"

# Copy runtime libraries
cp "$LIB_DIR"/*.so "$LIB_OUT/"

echo "Build complete: build/bin/app"
