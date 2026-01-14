#!/usr/bin/env bash
set -e

# ===============================
# Build mode: dev | release
# ===============================
BUILD_MODE=${1:-dev}

ROOT_DIR="$(pwd)"
BUILD_DIR="$ROOT_DIR/build/$BUILD_MODE"
BIN_DIR="$BUILD_DIR/bin"
LIB_OUT="$BUILD_DIR/lib"
GEN_DIR=".pre-build"

CXX=g++

COMMON_FLAGS="-std=c++17 -Wall -Wextra"

if [ "$BUILD_MODE" = "dev" ]; then
  CXXFLAGS="$COMMON_FLAGS -g -O0 -DDEBUG"
  echo "=== DEV BUILD (debug) ==="
else
  CXXFLAGS="$COMMON_FLAGS -O2 -DNDEBUG"
  echo "=== RELEASE BUILD ==="
fi

mkdir -p "$BIN_DIR" "$LIB_OUT" "$GEN_DIR"

# ===============================
# Generate embedded views
# ===============================
GENERATED_VIEW_CPP="$GEN_DIR/views.cpp"
echo "=== Embedding Views as Binary ==="

cat > "$GENERATED_VIEW_CPP" <<EOF
#include <map>
#include <string>
#include <vector>
#include "../vendor/agai/include/v1/agai.h"

namespace Agai {
EOF

if [ -d "src/views" ]; then

  # --- 1. Embed all view files (recursive) ---
  while IFS= read -r view; do
    [ -f "$view" ] || continue

    rel="${view#src/views/}"
    echo "  [BINARY] $rel"
    xxd -i "$view" >> "$GENERATED_VIEW_CPP"
  done < <(find src/views -type f)

  # --- 2. Generate map accessor ---
  echo "std::map<std::string, std::vector<unsigned char>> register_embedded_views() {" >> "$GENERATED_VIEW_CPP"
  echo "  return {" >> "$GENERATED_VIEW_CPP"

  while IFS= read -r view; do
    [ -f "$view" ] || continue

    # logical id: path → dot, drop extension
    rel="${view#src/views/}"
    no_ext="${rel%.*}"
    id="${no_ext//\//.}"

    # symbol name: must match xxd exactly
    var_name=$(echo "$rel" | sed 's/[^a-zA-Z0-9]/_/g')

    echo "    {\"$id\", std::vector<unsigned char>(src_views_${var_name}, src_views_${var_name} + src_views_${var_name}_len)}," \
      >> "$GENERATED_VIEW_CPP"
  done < <(find src/views -type f)

  echo "  };" >> "$GENERATED_VIEW_CPP"
  echo "}" >> "$GENERATED_VIEW_CPP"

else
  echo "std::map<std::string, std::vector<unsigned char>> register_embedded_views() { return {}; }" \
    >> "$GENERATED_VIEW_CPP"
fi


echo "}" >> "$GENERATED_VIEW_CPP"

# ===============================
# Collect source files
# ===============================
echo "=== Scanning C++ source files ==="
# CPP_FILES="$GENERATED_VIEW_CPP"

while IFS= read -r cpp; do
  echo "  [CPP] $cpp"
  CPP_FILES="$CPP_FILES $cpp"
done < <(
  find . -type f -name "*.cpp" ! -path "./build/*"
)

# ===============================
# Include directories
# ===============================
echo "=== Scanning include directories ==="
INC_FLAGS=$(
  find . -type f \( -name "*.h" -o -name "*.hpp" \) ! -path "./build/*" \
  | xargs -n1 dirname \
  | sort -u \
  | sed 's/^/-I/'
)

# ===============================
# Shared libraries
# ===============================
echo "=== Scanning shared libraries ==="
LIB_FLAGS=""
LIB_DIR_FLAGS=""

while IFS= read -r so; do
  dir=$(dirname "$so")
  base=$(basename "$so")
  name=$(echo "$base" | sed 's/^lib//;s/\.so$//')
  echo "  [SO ] $so"
  LIB_FLAGS="$LIB_FLAGS -l$name"
  LIB_DIR_FLAGS="$LIB_DIR_FLAGS -L$dir"
done < <(
  find . -type f -name "*.so" ! -path "./build/*"
)

LIB_DIR_FLAGS=$(echo "$LIB_DIR_FLAGS" | tr ' ' '\n' | sort -u | tr '\n' ' ')

# ===============================
# Compile & link
# ===============================
echo "=== Compiling & linking ==="
$CXX $CXXFLAGS \
  $CPP_FILES \
  $INC_FLAGS \
  $LIB_DIR_FLAGS \
  $LIB_FLAGS \
  -Wl,-rpath,'$ORIGIN/../lib' \
  -o "$BIN_DIR/app"

# ===============================
# Copy shared libraries
# ===============================
echo "=== Copying shared libraries ==="
find . -type f -name "*.so" ! -path "./build/*" -exec cp -v {} "$LIB_OUT/" \;

rm -r "$GEN_DIR"

echo "=== Build complete: $BIN_DIR/app ==="
