#!/bin/bash

# clean up the output directory and recreate it
clean() {
  rm -rf "$OUT"
  mkdir -p "$OUT"
}
