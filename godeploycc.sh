#!/usr/bin/env bash

CHANNEL_NAME="$1"
CHAINCODE_NAME="$2"
CHAINCODE_VERSION="$3"
: ${CHANNEL_NAME:="mychannel1"}
: ${CHAINCODE_NAME:="mycc"}
: ${CHAINCODE_VERSION:="1"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "----install chaincode (安装智能合约)----"
./deployFabricTool -r installchaincode -n $CHANNEL_NAME -ccname $CHAINCODE_NAME -version $CHAINCODE_VERSION
verifyResult $?
echo "-------instantiate chaincode (实例化智能合约)-------"
./deployFabricTool -r runchaincode -n $CHANNEL_NAME -ccname $CHAINCODE_NAME -version $CHAINCODE_VERSION
verifyResult $?
#echo "---1.4.x-upgrade chaincode (升级chaincode)----"
#./deployFabricTool -r upgradecc -n $CHANNEL_NAME -ccname $CHAINCODE_NAME
#verifyResult $?



