#!/usr/bin/env bash

CHANNEL_NAME="$1"
CHAINCODE_NAME="$2"
: ${CHANNEL_NAME:="mychannel1"}
: ${CHAINCODE_NAME:="mycc"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "-------make channel (创建业务通道)-------"
./deployFabricTool -c channel -n $CHANNEL_NAME
verifyResult $?

echo "-------update anchor (更新通道机构锚节点)-------"
./deployFabricTool -r updateanchor -n $CHANNEL_NAME
verifyResult $?

echo "-------peer join channel(节点加入通道)-------"
./deployFabricTool -r joinchannel -n $CHANNEL_NAME
verifyResult $?

echo "-------install chaincode (安装智能合约)-------"
./deployFabricTool -r installchaincode -ccname $CHAINCODE_NAME
verifyResult $?

echo "-------instantiate chaincode (实例化智能合约)-------"
./deployFabricTool -r runchaincode -n $CHANNEL_NAME -ccname $CHAINCODE_NAME
verifyResult $?
