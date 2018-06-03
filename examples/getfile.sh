#!/bin/bash

FILE=$1
FNAME=$(basename "$FILE")

if [[ $FILE != "/"* ]]; then
  FILE="/$FILE"
fi

curl -s "http://localhost:8080/blobs$FILE" --output "$FNAME"
