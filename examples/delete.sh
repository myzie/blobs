#!/bin/bash

set -e

FILE="$1"

curl -s -X DELETE "http://localhost:8080/blobs${FILE}"
