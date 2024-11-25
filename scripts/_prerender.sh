#!/bin/bash

hooks=()

prerender_hook() {
  hooks+=("$1")
}

prerendered=false

prerender() {
  echo "Running prerender hooks"
  for hook in "${hooks[@]}"; do
    echo "Running prerender hook: $hook"
    $hook
  done
  
  echo "Prerendering..."
  for file in $FILES; do
    local output=$(get_output "$file")
    url="$(dirname $output)"
    url="${url##"$OUT"}/"
    _=$(render "$file")
  done
 
  prerendered=true
}
