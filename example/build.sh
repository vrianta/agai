#!/usr/bin/env bash
set -e

ROOT_DIR="$(pwd)"
BUILD_DIR="$ROOT_DIR/build"
BIN_DIR="$BUILD_DIR/bin"
LIB_OUT="$BUILD_DIR/lib"

CXX=g++
CXXFLAGS="-O2 -std=c++17"

rm -rf "$BUILD_DIR"
mkdir -p "$BIN_DIR" "$LIB_OUT"

CPP_FILES=""
INC_FLAGS=""
LIB_FLAGS=""
LIB_DIR_FLAGS=""

echo "=== Scanning C++ source files ==="
find . \
  -type f \
  -name "*.cpp" \
  ! -path "./build/*" \
  | while read -r cpp; do
      echo "  [CPP] $cpp"
      CPP_FILES="$CPP_FILES $cpp"
      echo "$CPP_FILES" > /tmp/.cpp_files
    done
CPP_FILES=$(cat /tmp/.cpp_files 2>/dev/null || true)
rm -f /tmp/.cpp_files

[ -z "$CPP_FILES" ] && echo "ERROR: no .cpp files found" && exit 1

echo "=== Scanning include directories ==="
find . \
  -type f \
  \( -name "*.h" -o -name "*.hpp" \) \
  ! -path "./build/*" \
  | while read -r hdr; do
      dir=$(dirname "$hdr")
      echo "  [INC] $dir"
      INC_FLAGS="$INC_FLAGS -I$dir"
      echo "$INC_FLAGS" > /tmp/.inc_flags
    done
INC_FLAGS=$(cat /tmp/.inc_flags 2>/dev/null | tr ' ' '\n' | sort -u | tr '\n' ' ')
rm -f /tmp/.inc_flags

echo "=== Scanning shared libraries ==="
find . \
  -type f \
  -name "*.so" \
  ! -path "./build/*" \
  | while read -r so; do
      dir=$(dirname "$so")
      base=$(basename "$so")
      name=$(echo "$base" | sed 's/^lib//;s/\.so$//')

      echo "  [SO ] $so -> -l$name"

      LIB_FLAGS="$LIB_FLAGS -l$name"
      LIB_DIR_FLAGS="$LIB_DIR_FLAGS -L$dir"

      echo "$LIB_FLAGS" > /tmp/.lib_flags
      echo "$LIB_DIR_FLAGS" > /tmp/.lib_dir_flags
    done

LIB_FLAGS=$(cat /tmp/.lib_flags 2>/dev/null || true)
LIB_DIR_FLAGS=$(cat /tmp/.lib_dir_flags 2>/dev/null | tr ' ' '\n' | sort -u | tr '\n' ' ')
rm -f /tmp/.lib_flags /tmp/.lib_dir_flags

echo "=== Compiling & linking ==="
echo "$CXX $CXXFLAGS <cpp> <includes> <libs>"

$CXX $CXXFLAGS \
  $CPP_FILES \
  $INC_FLAGS \
  $LIB_DIR_FLAGS \
  $LIB_FLAGS \
  -Wl,-rpath,'$ORIGIN/../lib' \
  -o "$BIN_DIR/app"

echo "=== Copying shared libraries ==="
find . \
  -type f \
  -name "*.so" \
  ! -path "./build/*" \
  -exec cp -v {} "$LIB_OUT/" \;

if [ -d "Views" ]; then
  echo "=== Copying Views directory ==="
  cp -Rv "Views" "$BIN_DIR/"
fi

echo "Build complete: build/bin/app"
