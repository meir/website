#!/bin/bash

# use a watcher such as Overseer for nvim in order to have hot reloading
source ./config.sh

./generate.sh
cd "$OUT"
python3 -m http.server 8080

