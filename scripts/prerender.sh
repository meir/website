#!/bin/bash

hooks=()

prerender_hook() {
  hooks+=("$1")
}

prerendered=false

prerender() {
  echo "Prerendering..."
  for file in $FILES; do
    _=$(render "$file")
  done

  echo "Running prerender hooks"
  for hook in "${hooks[@]}"; do
    echo "Running prerender hook: $hook"
    $hook
  done
  
  prerendered=true
}
