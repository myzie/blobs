#!/bin/bash

FILE=$1
FNAME=$(basename "$FILE")

if [[ $FILE != "/"* ]]; then
  FILE="/$FILE"
fi

curl -s                                    \
  -H "Authorization: Bearer $BLOBS_TOKEN"  \
  "http://localhost:8080/blobs$FILE" --output "$FNAME"
