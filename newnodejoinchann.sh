#!/usr/bin/env bash

CHANNEL_NAME="$1"
NODE_NAME="$2"
: ${CHANNEL_NAME:="mychannel"}
: ${NODE_NAME:="peer3.test.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------新peer节点加入channel-------"
./deployFabricTool -r joinchannel -n $CHANNEL_NAME -nodename $NODE_NAME
verifyResult $?