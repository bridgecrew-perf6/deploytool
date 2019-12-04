#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "-------put crypto and config file (向节点传输证书、配置文件)-------"
./deployFabricTool -p all
verifyResult $?

echo "-------start node (启动节点)-------"
./deployFabricTool -s orderer
verifyResult $?

./deployFabricTool -s peer
verifyResult $?

