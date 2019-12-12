#!/usr/bin/env bash

#chmod +x ./bin/*
#sudo chown ubuntu:ubuntu /etc/hosts
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! FAIL !!!!!!!!!!!!!!!!"
    exit 1
  fi
}

echo "-------make config.yaml (生成configtx配置文件)-------"
./deployFabricTool -f configtx
verifyResult $?

echo "-------make docker-compose.yaml (生成节点docker启动文件)-------"
./deployFabricTool -f node
verifyResult $?

if [ "$1" = "" ]; then
    echo "-------make crypto-config.yaml (生成证书配置文件)-------"
    ./deployFabricTool -f crypto-config
    verifyResult $?
    echo "-------if ca exist, will start fabric-ca (通过fabric-ca生成证书)-------"
    ./deployFabricTool -s ca
    verifyResult $?
    echo "-------make crypto-config dir (生成证书目录)-------"
    ./deployFabricTool -c crypto-config
    verifyResult $?
else
    echo "$1"
    echo "-------use exist crypto-config (使用config目录里的证书)-------"
fi

echo "-------make genesisblock (生成创世区块)-------"
./deployFabricTool -c genesisblock
verifyResult $?
