#!/bin/bash



while getopts "p:" OPTION
do # (colon) denotes a flag that requires a value.
  case $OPTION in
    p)  port=$OPTARG;;
    *)  exit;;
  esac
done


lookup=":500$port"
cid=$(docker ps | grep $lookup | awk '{print$1}');
echo "Logging into container $cid"
docker exec -it $cid bash
