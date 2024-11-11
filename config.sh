#!/bin/bash

SRC="./src"
ASSETS="./assets"
OUT="./site"
FILES=$(find "$SRC" -type f ! -name ".htm(l)")
CNAME="yesimhuman.dev"

# add other variables here that should be globally accessible

GISCUS_REPO="meir/website"
GISCUS_REPO_ID="R_kgDOIyQOzA"
GISCUS_CATEGORY_ANNOUNCEMENTS="DIC_kwDOIyQOzM4CkKe4"
GISCUS_CATEGORY_GENERAL="DIC_kwDOIyQOzM4CkKe5"
GISCUS_CATEGORY_IDEAS="DIC_kwDOIyQOzM4CkKe7"
GISCUS_CATEGORY_POLLS="DIC_kwDOIyQOzM4CkKe9"
GISCUS_CATEGORY_QNA="DIC_kwDOIyQOzM4CkKe6"
GISCUS_CATEGORY_SHOWNTELL="DIC_kwDOIyQOzM4CkKe8"
