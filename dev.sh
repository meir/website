#!/bin/bash

# Load configuration
source ./config.sh

IP=127.0.0.1
PORT=8080

list=""
pid=""
running=true

# Ensure the Python server is killed on exit
cleanup() {
  if [ -n "$pid" ]; then
    echo "Killing server with PID: $pid"
    kill $pid 2>/dev/null
    wait $pid 2>/dev/null
  fi
  running=false
}
trap cleanup EXIT

# netcat is my nightmare
# 
# get_request_path() {
#   local request="$1"
#   # only get first line
#   local path=$(echo "$request" | head -n 1 | cut -d ' ' -f 2)
#   echo "$path"
# }
#
# get_requested_file() {
#   local request="$1"
#   local path=$(get_request_path "$request")
#   local file="$OUT$path"
#
#   if [ ! -f "$file" ]; then
#     file="$file/index.htm"
#   fi
#
#   if [ -f "$file" ]; then
#     echo "$file" | tr -s /
#   else
#     echo ""
#   fi
# }
#
# handle_request() {
#   content=""
#   while read line; do
#     if [ -z "$line" ]; then
#       break
#     fi
#     content="$content$line\n"
#   done
#
#   local file=$(get_requested_file "$request")
#   if [ -f "$file" ]; then
#     cat <<EOF
# HTTP/1.1 200 OK
# Content-Type: $(file -b --mime-type "$file")
#
# $(<"$file")
# ${request}
# $(date)
# EOF
#   else
#     cat <<EOF
# HTTP/1.1 404 Not Found
#
# ${request}
# 404 Not Found
# EOF
#   fi
# }
#
# netcat_server() {
#   # check if on osx
#   if [ "$(uname)" == "Darwin" ]; then
#     osx_netcat_server
#     return
#   fi
#
#   echo "Starting server on port $IP:$PORT"
#   while true; do
#     response=$(nc -l -p $PORT | handle_request)
#   done
# }
#
# osx_netcat_server() {
#   echo "Starting OSX server on port $IP:$PORT"
#   while $running; do
#     nc -l $IP $PORT | handle_request | nc -l $PORT
#   done
# }

while true; do
  # Generate a new list of file modification times
  new_list=$(find "$COMPONENTS" "$ASSETS" "$SRC" -type f -exec stat -c "%Y %n" {} + | sort -n | md5sum)
  
  # Check if the list has changed
  if [ "$list" != "$new_list" ]; then
    echo "Change detected from $list to $new_list"
    list="$new_list"
    ./generate.sh

    # If a server is running, kill it
    if [ -n "$pid" ]; then
      cleanup
    fi

    # Start a new server in the background and store its PID
    if [ -d "$OUT" ]; then
      # netcat_server &
      python3 -m http.server $PORT --bind $IP --directory "$OUT" &
      pid=$!
      echo "Server PID: $pid"
    else
      echo "Directory '$OUT' does not exist. Please check your config.sh."
      exit 1
    fi
  fi

  # Sleep for a second before rechecking
  sleep 1
done
