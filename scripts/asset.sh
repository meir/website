#!/bin/bash

# copy everything from the assets folder into the root of /assets/ for the website for easy access
copy_assets() {
  files=$(find "$ASSETS" -type f)
  mkdir -p $OUT/assets
  for file in $files; do
    cp $file $OUT/assets
  done
}

find_assets() {
  local type="$1"
  local name="$2"

  local files="$(find "$OUT/assets/" -type f ! -name "$type" | grep "$name")"
  echo $files
}
