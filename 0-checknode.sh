#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

if [[ "$1" == "" ]]; then
    echo "-------writehost (写域名映射)-------"
    ./deployFabricTool -r writehost
    verifyResult $?
else
    echo "-------check node (验证所有节点)-------"
    ./deployFabricTool -r checknode
    verifyResult $?
fi


