#!/bin/bash

TARGET=deployFabricTool
TARGET_TAR=$TARGET.tar.gz
TARGET_PATH=$GOPATH/src/github.com/peersafe/$TARGET

echo "##########################################"
echo "----------build $TARGET image----------"
echo "##########################################"

if [ -f $TARGET_PATH/$TARGET ]; then
    rm -f $TARGET_PATH/$TARGET
fi

if [ -d ./$TARGET ]; then
    echo "remove old file"
    rm -rf ./$TARGET
fi
mkdir ./$TARGET

echo "build $TARGET wait ...."
cd $TARGET_PATH
go build
cd -

if [ -f $TARGET_PATH/$TARGET ]; then
    mv $TARGET_PATH/$TARGET ./$TARGET/
    cp -r $TARGET_PATH/bin ./$TARGET/
    cp -r $TARGET_PATH/data ./$TARGET/
    cp -r $TARGET_PATH/scripts ./$TARGET/
    cp -r $TARGET_PATH/templates ./$TARGET/
    cp $TARGET_PATH/*.sh ./$TARGET/
else
    echo "--------ERROR: $TARGET process has not been build.----------"
    exit
fi

tar -zcvf $TARGET_TAR $TARGET

docker rmi gmhyperledger/deploy-tool:latest
docker build -t gmhyperledger/deploy-tool:latest .

rm -rf $TARGET
rm -rf $TARGET_TAR
rm -rf $TARGET_PATH/config
