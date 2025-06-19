#!/usr/bin/env bash

if command -v watchexec ] >/dev/null 2>&1; then
  watchexec --restart --watch src/ --watch components/ --watch assets/ -- "daisy build && cd site && python3.12 -m http.server 8000"
elif command -v , >/dev/null 2>&1; then
  , watchexec --restart --watch src/ --watch components/ --watch assets/ -- "daisy build && cd site && python3.12 -m http.server 8000"
else
  echo "watchexec is not installed. Please install it to use this script."
fi
