#!/bin/bash

# copy everything from the assets folder into the root of /assets/ for the website for easy access
copy_assets() {
  files=$(find "$ASSETS" -type f)
  mkdir -p $OUT/assets
  for file in $files; do
    cp $file $OUT/assets
  done
}
