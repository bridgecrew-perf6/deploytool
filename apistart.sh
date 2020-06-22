#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "----only start apiserver--"
./deployFabricTool -f api
verifyResult $?
./deployFabricTool -p api
verifyResult $?
./deployFabricTool -s api
verifyResult $?
