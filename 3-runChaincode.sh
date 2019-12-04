#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "-------make channel (创建业务通道)-------"
./deployFabricTool -c channel -n mychannel
verifyResult $?

echo "-------update anchor (更新通道机构锚节点)-------"
./deployFabricTool -r updateanchor -n mychannel
verifyResult $?

echo "-------peer join channel(节点加入通道)-------"
./deployFabricTool -r joinchannel -n mychannel
verifyResult $?

echo "-------install chaincode (安装智能合约)-------"
./deployFabricTool -r installchaincode
verifyResult $?

echo "-------instantiate chaincode (实例化智能合约)-------"
./deployFabricTool -r runchaincode -n mychannel
verifyResult $?
