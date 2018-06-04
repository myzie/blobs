#!/bin/bash

ssh-keygen -t rsa -b 2048 -f jwt.key

openssl rsa -in jwt.key -pubout -outform PEM -out jwt.key.pub

