#!/bin/bash

set -e
set -x

FILE="$1"
JSON="$2"

curl -s -X PUT                        \
  -d "$JSON"                          \
  -H "Content-Type: application/json" \
  "http://localhost:8080/blobs${FILE}" | jq '.'
