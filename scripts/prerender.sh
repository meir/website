#!/bin/bash

hooks=()

prerender_hook() {
  hooks+=("$1")
}

prerender() {
  echo "Running prerender hooks"
  for hook in "${hooks[@]}"; do
    echo "Running prerender hook: $hook"
    $hook
  done
}
