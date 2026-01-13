#!/usr/bin/env bash
set -e

ROOT_DIR="$(pwd)"
BUILD_DIR="$ROOT_DIR/build"
LIB_DIR="$BUILD_DIR/lib"
INC_DIR="$BUILD_DIR/include"

CXX=g++
CXXFLAGS="-fPIC -O2 -std=c++17"
LDFLAGS="-shared -Wl,-undefined,dynamic_lookup"

rm -rf "$BUILD_DIR"
mkdir -p "$LIB_DIR" "$INC_DIR"

# -------- build: one cpp = one so --------
for ver_dir in v*/ ; do
  ver="${ver_dir%/}"
  echo "Processing $ver"

  find "$ver" \
    -type f \
    -name "*.cpp" \
    ! -path "*/example/*" \
    | while read -r cpp; do

        base="$(basename "$cpp" .cpp)"
        so_name="libagai_${ver}_${base}.so"

        echo "  Building $so_name"
        $CXX $CXXFLAGS "$cpp" $LDFLAGS -o "$LIB_DIR/$so_name"
    done

  # -------- copy headers for this version --------
  find "$ver" \
    -type f \
    \( -name "*.h" -o -name "*.hpp" \) \
    -print0 | while IFS= read -r -d '' header; do

      rel="${header#./}"
      dest="$INC_DIR/$(dirname "$rel")"

      mkdir -p "$dest"
      cp "$header" "$dest/"
  done

done
