#!/usr/bin/env bash

CHANNEL_NAME="$1"
CHAINCODE_NAME="$2"
CHAINCODE_VERSION="$3"
NODE_NAME="$4"
: ${CHANNEL_NAME:="mychannel"}
: ${CHAINCODE_NAME:="mycc"}
: ${CHAINCODE_VERSION:="1"}
: ${NODE_NAME:="peer3.test.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------新peer节点加入安装智能合约------"
./deployFabricTool -r installcctonewnode -n $CHANNEL_NAME -ccname $CHAINCODE_NAME  -version $CHAINCODE_VERSION -nodename $NODE_NAME
verifyResult $?