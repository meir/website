#!/bin/bash

SRC="./src"
ASSETS="./assets"
OUT="./site"
FILES=$(find "$SRC" -type f ! -name ".htm(l)")
CNAME="yesimhuman.dev"

# add other variables here that should be globally accessible
