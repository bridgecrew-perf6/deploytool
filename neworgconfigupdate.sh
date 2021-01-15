#!/usr/bin/env bash

ORG_ID="$1"
CHANNEL_NAME="$2"
: ${ORG_ID:="test"}
: ${CHANNEL_NAME:="mychannel"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------更新新组织到指定通道配置块-------"
./deployFabricTool -r addorgtoconfigblock -orgid $ORG_ID  -n $CHANNEL_NAME
verifyResult $?