#!/bin/bash

set -e

source ./config.sh

# load all the scripts in
for mod in scripts/*.sh; do
  source "$mod"
done

# run functions to build the site
clean
copy_assets
prerender
render_all_files
