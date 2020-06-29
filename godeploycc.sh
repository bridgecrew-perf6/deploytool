#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "----install chaincode (安装智能合约)----"
./deployFabricTool -r installchaincode -n mychannel
verifyResult $?
if [ "$1" == "" ]; then
  echo "-------instantiate chaincode (实例化智能合约)-------"
  ./deployFabricTool -r runchaincode -n mychannel
  verifyResult $?
else
  echo "----upgrade chaincode (升级chaincode)----"
  ./deployFabricTool -r upgradecc -n mychannel
  verifyResult $?
fi


