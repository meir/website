#!/bin/bash

echo "Escaping single quotes in src files..."
for file in $FILES; do
  content=$(sed -E "s|([^\\])'|\1\\\'|g" "$file")
  echo "$content" > "$file"
done
echo "Done escaping single quotes in src files."
