#!/bin/bash

# return the path without the src prefix
get_path() {
  local path=$(echo $1 | sed "s|$SRC||")
  echo $path
}

# return the name of the file without htm(l) extension
get_name() {
  local name=$(basename $1)
  local name=$(echo $name | sed 's|\.htm||' | sed 's|\.html||')
  echo $name
}

# return the ouput path of the file
# example: ./src/topics/keyboards.htm > ./site/topics/keyboards/index.htm
# this makes it so that the user can go to /topics/keyboards instead of /topics/keyboards.htm
get_output() {
  local path=$(get_path $1)
  # if the file is already named index we dont have to make an index dir
  if [ "$(get_name $1)" == "index" ]; then
    path=$(dirname $path)
  else
    path="$(dirname $path)/$(get_name $1)"
  fi
  mkdir -p "${OUT}${path}"
  touch "${OUT}${path}/index.htm"
  echo "${OUT}${path}/index.htm"
}

# render the content of the file
render() {
  export file="$1"
  local content=$(eval "cat <<EOD
$(<$file)
EOD
")
  echo "$content"
}

# render all the files in the src dir and output them in the output dir
render_all_files() {
  for file in $FILES; do
    local output=$(get_output "$file")
    echo "Rendering $file to $output"
    url="$(dirname $output)"
    url="${url##"$OUT"}/"
    root_file="$file"
    render "$file" | ./prettify.sh > "$output" 
  done
  
  echo "Done rendering all files"
}
