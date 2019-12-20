#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}
echo "----clean all node data---清除所有网络信息"
./deployFabricTool -d all
verifyResult $?
