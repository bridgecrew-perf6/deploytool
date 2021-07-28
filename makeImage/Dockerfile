FROM peersafes/deploy-base:ubuntu

# tool
ADD ./deployFabricTool.tar.gz  /opt
COPY ./chaincode/  $GOPATH/src/github.com/chaincode
WORKDIR /opt/deployFabricTool