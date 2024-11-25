#!/bin/bash

echo "creating tmp dir"
tmp=$(mktemp -d)
trap 'echo "removing tmp dir"; rm -rf $tmp' EXIT

