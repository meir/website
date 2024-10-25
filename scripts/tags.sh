#!/bin/bash

prerender_hook "prerender_tags"

prerender_tags() {
  echo "{}" > $tmp/tags.json
}

tag() {
  if [ $prerendered == true ]; then
    return
  fi

  local key=$1
  local value=$file
  local title=$title
  jq -c ".$key += [{\"file\": \"$value\", \"title\": \"$title\"}]" <<< $(<$tmp/tags.json) > $tmp/tags.json
}

get_tagged() {
  local key=$1
  local pages=$(jq -c ".$key[]" <<< $(<$tmp/tags.json))
  echo $pages
}

get_url() {
  local page="$@"
  
  local url=$(echo "$page" | jq -r ".file")
  local url=$(get_path $url)
  echo $url
}

get_title() {
  local page="$@"

  local title=$(echo "$page" | jq -r ".title")
  echo $title
}

