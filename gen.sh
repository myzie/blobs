#!/bin/bash

mkdir -p gen/server
mkdir -p gen/client

swagger generate server -f ./swagger.yml -t gen/server
swagger generate client -f ./swagger.yml -t gen/client
