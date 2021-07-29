#!/usr/bin/env bash

CHANNEL_NAME="$1"
CHAINCODE_NAME="$2"
NODE_NAME="$3"
: ${CHANNEL_NAME:="mychannel"}
: ${CHAINCODE_NAME:="mycc"}
: ${NODE_NAME:="all"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "-------test chaincode (测试chaincode)-------"
./deployFabricTool -r testcc -n $CHANNEL_NAME -ccname $CHAINCODE_NAME -func invoke -nodename $NODE_NAME
verifyResult $?

