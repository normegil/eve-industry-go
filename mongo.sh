#!/bin/sh
docker run \
  -d \
  -p 27017:27017 \
  --name eve-vulcan-mongo \
  mongo