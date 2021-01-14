#!/usr/bin/env bash

ORG_ID="$1"
NODE_NAME="$2"
: ${ORG_ID:="test"}
: ${NODE_NAME:="peer3.test.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------增加组织内节点证书-------"
./deployFabricTool -r addorgnodecert -orgid $ORG_ID
verifyResult $?

echo "-------证书传递到远程服务器-------"
./deployFabricTool -r putnodecrypto -nodename $NODE_NAME
verifyResult $?

echo "-------生成peer节点yaml文件-------"
./deployFabricTool -r createpeeryaml -nodename $NODE_NAME
verifyResult $?

echo "-------启动新peer节点-------"
./deployFabricTool -r runaddnode -nodename $NODE_NAME
verifyResult $?