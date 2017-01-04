#!/bin/bash

echo
echo "1. IDENTIFYING CONTAINERS TO CLEAN"
docker ps -a
echo
containers="$(docker ps -a | grep 'fabric-src\|dev-\|/bin/sh -c ' | awk '{print$1}' | tr '\n' ' ')"
if [ "$containers" != "" ]; then
  echo "docker rm -f $containers"
  docker rm -f $containers
else
  echo "No containers to remove."
fi
echo
docker ps -a
echo
echo
echo "2. IDENTIFYING IMAGES TO CLEAN"
docker images
echo
images="$(docker images | grep 'dev-\|fabric-peer\|<none>' | awk '{print$3}' | tr '\n' ' ')"
if [ "$images" != "" ]; then
  echo "docker rmi -f $images"
  docker rmi -f $images
else
  echo "No images to remove."
fi
echo
docker images
echo
echo
echo "3. REBUILDING PEER IMAGE"
(cd $GOPATH/src/github.com/hyperledger/fabric/ && echo making peer: && make peer && echo && echo making peer-image: && make peer-image && echo)
docker ps -a
docker images
