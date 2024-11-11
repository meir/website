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
  local title=$(echo $title | jq -aRs .)
  local entry=$(cat <<EOF
[{
  "file": "$url",
  "title": "$title"
}]
EOF
)
  jq -c ".\"$key\" += $entry" <<< $(<$tmp/tags.json) > $tmp/tags.json
}

# get all the pages attached to a specific tag
get_tagged() {
  local key=$1
  local pages=$(jq -c ".\"$key\"[]" <<< $(<$tmp/tags.json))
  echo $pages
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
  echo $title | jq -r
}

