FROM peersafes/deploy-base:latest

# tool
ADD ./deployFabricTool.tar.gz  /opt
#ADD ./gosrc.tar.gz $GOPATH/src/
COPY ./chaincode/  $GOPATH/src/github.com/chaincode
WORKDIR /opt/deployFabricTool