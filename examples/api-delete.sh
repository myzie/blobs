#!/bin/bash

set -e

FILE="$1"

curl -s -H "Authorization: Bearer $BLOBS_TOKEN" -X DELETE "http://localhost:8080/blobs${FILE}"
