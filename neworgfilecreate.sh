#!/usr/bin/env bash

ORG_ID="$1"
: ${ORG_ID:="test"}
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------增加新组织证书-------"
./deployFabricTool -r addorgnodecert -orgid $ORG_ID
verifyResult $?

echo "-------创建新组织配置文件-------"
./deployFabricTool -r createneworgconfigtxfile -orgid $ORG_ID
verifyResult $?