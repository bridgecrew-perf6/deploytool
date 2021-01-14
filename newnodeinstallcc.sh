#!/usr/bin/env bash

NODE_NAME="$1"
: ${NODE_NAME:="peer3.test.example.com"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------新peer节点加入安装智能合约------"
./deployFabricTool -r installcctonewnode -nodename $NODE_NAME
verifyResult $?