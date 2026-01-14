#!/usr/bin/env bash
set -e

ROOT_DIR="$(pwd)"
BUILD_DIR="$ROOT_DIR/build"
BIN_DIR="$BUILD_DIR/bin"
LIB_OUT="$BUILD_DIR/lib"
# New directory for generated source code
GEN_DIR=".pre-build"

CXX=g++
CXXFLAGS="-O2 -std=c++17"

rm -rf "$BUILD_DIR"
mkdir -p "$BIN_DIR" "$LIB_OUT" "$GEN_DIR"

# --- NEW: GENERATE BINARY VIEWS ---
GENERATED_VIEW_CPP="$GEN_DIR/views.cpp"
echo "=== Embedding Views as Binary ==="
echo "#include <map>" > "$GENERATED_VIEW_CPP"
echo "#include <string>" >> "$GENERATED_VIEW_CPP"
echo "#include <vector>" >> "$GENERATED_VIEW_CPP"
echo "#include \"../vendor/agai/include/v1/agai.h\"" >> "$GENERATED_VIEW_CPP"

echo "namespace Agai {" >> "$GENERATED_VIEW_CPP"


if [ -d "src/views" ]; then
  # 1. Generate hex arrays for each file
  for view in src/views/*; do
    if [ -f "$view" ]; then
      filename=$(basename "$view")
      echo "  [BINARY] $filename"
      # xxd -i creates: unsigned char name[] = { ... }; and unsigned int name_len = ...;
      xxd -i "$view" >> "$GENERATED_VIEW_CPP"
    fi
  done

  # 2. Generate the accessor function
  echo "std::map<std::string, std::vector<unsigned char>> register_embedded_views() {" >> "$GENERATED_VIEW_CPP"
  echo "  return {" >> "$GENERATED_VIEW_CPP"
  for view in src/views/*; do
    if [ -f "$view" ]; then
      filename=$(basename "$view")
      # Sanitize filename for variable name (matching xxd's behavior)
      var_name=$(echo "$filename" | sed 's/[^a-zA-Z0-9]/_/g')
      echo "    {\"$filename\", std::vector<unsigned char>(src_views_${var_name}, src_views_${var_name} + src_views_${var_name}_len)}," >> "$GENERATED_VIEW_CPP"
    fi
  done
  echo "  };" >> "$GENERATED_VIEW_CPP"
  echo "}" >> "$GENERATED_VIEW_CPP"
  echo "}" >> "$GENERATED_VIEW_CPP"
else
  echo "  [SKIP] No src/views directory found."
  # Create empty function if no views exist to prevent linker errors
  echo "std::map<std::string, std::vector<unsigned char>> get_embedded_views() { return {}; }" >> "$GENERATED_VIEW_CPP"
fi

echo "=== Scanning C++ source files ==="
# Notice we now include the $GENERATED_VIEW_CPP in the list
# CPP_FILES="$GENERATED_VIEW_CPP"
find . \
  -type f \
  -name "*.cpp" \
  ! -path "./build/*" \
  | while read -r cpp; do
      echo "  [CPP] $cpp"
      echo "$CPP_FILES $cpp" > /tmp/.cpp_files
      CPP_FILES=$(cat /tmp/.cpp_files)
    done
CPP_FILES=$(cat /tmp/.cpp_files 2>/dev/null || echo "$GENERATED_VIEW_CPP")
rm -f /tmp/.cpp_files

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

# Note: We no longer need to copy src/views to BIN_DIR because they are inside the binary!
# rm -r "$GENERATED_VIEW_CPP"
echo "Build complete: build/bin/app"
