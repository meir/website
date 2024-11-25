#!/bin/bash

set -e

source ./config.sh

# load all the scripts in
for mod in scripts/*.sh; do
  source "$mod"
done

case "$1" in
  test)
    shift
    eval "$@" ;;
  *)
    # run functions to build the site
    clean
    copy_assets
    prerender
    create_cname
    render_all_files
    ;;
esac
