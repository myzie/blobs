#!/bin/bash

FILE=$1

if [[ $FILE != "/"* ]]; then
  FILE="/$FILE"
fi

curl -s \
    -H "Content-Type: application/json" \
    "http://localhost:8080/blobs$FILE" | jq '.'
