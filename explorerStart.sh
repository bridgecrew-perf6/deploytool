#!/usr/bin/env bash

verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "----only start explorer--"
./deployFabricTool -f explorer
verifyResult $?
./deployFabricTool -p explorer
verifyResult $?
./deployFabricTool -s explorer
verifyResult $?
