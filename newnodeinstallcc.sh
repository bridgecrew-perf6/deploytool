#!/usr/bin/env bash

CHANNEL_NAME="$1"
CHAINCODE_NAME="$2"
NODE_NAME="$3"
: ${CHANNEL_NAME:="mychannel"}
: ${CHAINCODE_NAME:="mycc"}
: ${NODE_NAME:="peer3.test.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------新peer节点加入安装智能合约------"
./deployFabricTool -r installcctonewnode -n $CHANNEL_NAME -ccname $CHAINCODE_NAME -nodename $NODE_NAME
verifyResult $?