#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "----install chaincode (安装智能合约)----"
./deployFabricTool -r installchaincode
verifyResult $?
echo "----upgrade chaincode (升级chaincode)----"
./deployFabricTool -r upgradecc -n mychannel
verifyResult $?
