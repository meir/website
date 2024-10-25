#!/bin/bash

# load and render a component such as a layout or navigation
component() {
  file="$1"
  content=""

  while IFS= read -r line; do
    content=$(cat <<EOF
$content
$line
EOF)
  done

  render "./components/${file}.htm"
}
