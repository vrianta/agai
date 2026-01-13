#!/usr/bin/env bash
set -e

ROOT_DIR="$(pwd)"
BUILD_DIR="$ROOT_DIR/build"
LIB_DIR="$BUILD_DIR/lib"
INC_DIR="$BUILD_DIR/include"

CXX=g++
CXXFLAGS="-fPIC -O2 -std=c++17"
LDFLAGS="-shared"

rm -rf "$BUILD_DIR"
mkdir -p "$LIB_DIR" "$INC_DIR"

# Find cpp files excluding example/
find . \
  -type f \
  -name "*.cpp" \
  ! -path "./example/*" \
  | while read -r cpp; do

    name=$(basename "$cpp" .cpp)
    echo "Building $name.so"

    $CXX $CXXFLAGS "$cpp" $LDFLAGS -o "$LIB_DIR/lib$name.so" -Wl,-undefined,dynamic_lookup
done

echo "Copying headers (preserving structure)..."

find . \
  -type f \
  \( -name "*.h" -o -name "*.hpp" \) \
  ! -path "./example/*" \
  -print0 | while IFS= read -r -d '' header; do

    rel_path="${header#./}"
    dest="$INC_DIR/$(dirname "$rel_path")"

    mkdir -p "$dest"
    cp "$header" "$dest/"
done

