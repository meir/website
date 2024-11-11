#!/bin/bash

# load and render a component such as a layout or navigation
component() {
  file="$1"
  content=""

  if [ -p /dev/stdin ] || [ ! -t 0 ]; then
    while IFS= read -r line; do
      content+=$'\n'"$line"
    done
  fi

  render "./components/${file}.htm"
}
