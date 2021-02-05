#!/usr/bin/env bash

ORG_ID="$1"
NODE_NAME="$2"
: ${ORG_ID:="ord"}
: ${NODE_NAME:="orderer3.ord.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------证书传递到远程服务器-------"
./deployFabricTool -r putnodecrypto -nodename $NODE_NAME
verifyResult $?

echo "-------生成peer/orderer节点yaml文件-------"
./deployFabricTool -r createnodeyaml -nodename $NODE_NAME
verifyResult $?

echo "-------启动新peer/orderer节点-------"
./deployFabricTool -r runaddnode -nodename $NODE_NAME
verifyResult $?