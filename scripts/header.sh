#!/bin/bash

prerender_hook "prerender_header"

prerender_header() {
  echo "{}" > "$tmp/headers.json"
}

header_add() {
  if [ $prerendered == true ]; then
    return
  fi

  local content=""

  if [ -p /dev/stdin ] || [ ! -t 0 ]; then
    while IFS= read -r line; do
      content+=$'\n'"$line"
    done
  fi

  jq -c ".\"$url\" += [$(echo $content | jq -aRs .)]" <<< $(<$tmp/headers.json) > "$tmp/headers.json"
}

header_items() {
  jq -c ".\"$url\" | select(. != null) | join(\"\n\")" <<< $(<$tmp/headers.json) | jq -r
}
