#!/bin/bash

curl -s                                   \
  -H "Authorization: Bearer $BLOBS_TOKEN" \
  -H "Content-Type: application/json"     \
  http://localhost:8080/blobs | jq '.'
