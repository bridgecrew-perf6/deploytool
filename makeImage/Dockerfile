FROM peersafes/deploy-tool:baseimage

# tool
RUN rm -rf /opt/deployFabricTool
ADD ./deployFabricTool.tar.gz  /opt
ADD ./chaincodes.tar.gz  /opt/gopath/src
RUN cp /opt/deployFabricTool/bin/1.4/telnet /usr/bin