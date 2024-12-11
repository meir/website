#!/bin/bash

# register prerender hook
prerender_hook "prerender_tags"

# prerender hook for the tags to keep track of tags inside a temporary json file
prerender_tags() {
  echo "{}" > $tmp/tags.json
}

# add a file to a specific tag
tag() {
  if [ $prerendered == true ]; then
    return
  fi

  local key=$1
  if [[ -z $key || -z "$title" ]]; then
    return
  fi

  local result=$(jq --arg key "$key" --arg url "$url" --arg title "$title" -c '.[$key] += [{ file: $url, title: $title }]' <<< $(<$tmp/tags.json))
  echo "$result" > $tmp/tags.json
}

# get all the pages attached to a specific tag
get_tagged() {
  local key=$1
  local pages=$(jq --arg key "$key" -c 'if .[$key] | length > 0 then .[$key][] else empty end' <<< $(<$tmp/tags.json))
  local output=()
  while IFS= read -r line; do
    output+=("$line")
  done <<< "$pages"
  echo "${output[@]}"
}

# get the url of a page
get_url() {
  local page="$@"
  
  local url=$(echo "$page" | jq -r ".file")
  local url=$(get_path $url)
  echo $url
}

# get the title of a page
get_title() {
  local page="$@"

  local title=$(echo "$page" | jq -r ".title")
  echo $title
}

