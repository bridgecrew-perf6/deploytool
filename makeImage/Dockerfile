FROM peersafes/deploy-base:latest

# tool
ADD ./deployFabricTool.tar.gz  /opt
COPY ./chaincode/  $GOPATH/src/github.com/chaincode
WORKDIR /opt/deployFabricTool