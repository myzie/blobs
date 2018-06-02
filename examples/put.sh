#!/bin/bash

set -e

FILE="$1"
DEST="$2"
FILE_NAME=$(basename "${FILE}")
FILE_SIZE=$(wc -c <"${FILE}" | tr -d '[:space:]')

if [ -z "${DEST}" ]; then
  DEST="/${FILE_NAME}"
fi

curl -s -X PUT                  \
  -T "${FILE}"                  \
  -H "Sound-Size: ${FILE_SIZE}" \
  "http://localhost:8080/blobs${DEST}" | jq '.'
