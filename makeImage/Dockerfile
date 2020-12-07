FROM peersafes/deploy-base:latest

# tool
ADD ./deployFabricTool.tar.gz  /opt
ADD ./chaincode.tar.gz  $GOPATH/src/github.com/hyperledger/fabric/examples/
ADD ./poc_chaincode.tar.gz  $GOPATH/src/github.com/peersafe/fabric_poc/
ADD ./commoncc.tar.gz  $GOPATH/src/github.com/
ADD ./gosrc.tar.gz $GOPATH/src/
WORKDIR /opt/deployFabricTool