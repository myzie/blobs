#!/bin/bash

curl -s                                \
  -H "Content-Type: application/json"  \
  http://localhost:8080/blobs | jq '.'
