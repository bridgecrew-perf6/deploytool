# 1：环境说明

服务器：Centos 7.4 及以上系统， 8C16G 配置， 用root用户执行

# 2：环境准备

**安装 Docker 和 Docker-compose**
卸载可能存在的旧版本

```bash
yum remove docker \
    docker-client \
    docker-client-latest \
    docker-common \
    docker-latest \
    docker-latest-logrotate \
    docker-logrotate \
    docker-engine
```

## 安装 Docker

安装 Docker 仓库源

```
yum install -y yum-utils
yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

安装Docker

```bash
yum install -y docker-ce docker-ce-cli containerd.io
```

更换 Docker 默认仓库源
在 /etc/docker/daemon.json 中写入如下内容、

mkdir -p    /etc/docker

```
{
    "registry-mirrors": [
        "https://registry.aliyuncs.com",
        "https://docker.mirrors.ustc.edu.cn",
        "https://reg-mirror.qiniu.com",
        "https://hub-mirror.c.163.com"
    ]
}
```

启动 Docker 并让其开机自启

```bash
systemctl enable docker.service --now
```

添加 Docker TAB 自动补全

```shell
yum install -y bash-completion
source /usr/share/bash-completion/bash_completion
```

## 安装 Docker-compose

```bash
curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
```

# 3：Fabric网络部署

## 1)准备

选一台服务器作为部署网络的管理服务器

根据部署业务网络需求合理分配服务器资源部署fabric节点

例如：3台服务器： 尽量每台服务器部署一个orderer节点，peer节点分散部署在提供的服务器上。

​           5台服务器：3台单独部署3个orderer,2台单独部署2个peer

记录每台服务器IP地址，下面修改配置文件需要使用。

部署orderer节点服务器开放 "7050" 端口， 如果同时部署多orderer就要额外开放"5050"、"6050"

部署peer节点服务器开放 "7051" 端口， 如果同时部署多peer就要额外开放"8051"、"9051"等

建议开放端口： 5050，6050，7050， 7051，8051， 7054，8054

#### 安装依赖镜像

PS: fabric网络节点以docker容器运行，所以需要提前在所有服务器安装依赖镜像

在线方式

```bash
docker pull gmhyperledger/fabric-orderer:2.2.1-gm     
docker pull gmhyperledger/fabric-peer:2.2.1-gm                 
docker pull gmhyperledger/fabric-ccenv:2.2.1-gm             
docker pull gmhyperledger/fabric-baseos:2.2.1-gm    
```

离线方式

```bash
docker load -i xxx.tar
```

#### 创建部署目录

PS: 在选取得管理服务器执行命令

```bash
mkdir -p ~/deploy && cd  ~/deploy
```

#### 增加资源文件 

PS: 如果要部署自定义智能合约可以将源码目录放到~/deploy目录下

#### 增加启动配置文件

node.json

```json
{
  "fabricVersion":"2.0","domain":"irchain.net",
  "sshUserName":"root","sshPwd":"xxx","sshPort":"22","sshKey":"/etc/login.pem",
  "ccInit":"'{\"Args\":[\"init\"]}'",
  "ccName":"mycc","ccVersion":"1","ccInstallType":"path",
  "testArgs":"'{\"Args\":[\"saveData\"\\,\"key\"\\,\"value\"]}'",
  "testArgs1":"'{\"Args\":[\"getDataByMsgSn\"\\,\"key\"]}'",
  "ccPath":"github.com//chaincode/commoncc",
  "chan_counts":1,"mountPath": "/data", "caType": "cryptogen",
  "consensusType":"raft", "imagePre":"gmhyperledger","imageTag":"2.2.1-gm","log":"info",
  "batchTime":"1s", "batchSize":100, "batchPreferred":"1024 KB",
  "orderers":[
    {"ip":"10.0.2.15","id":"0","orgId":"ord1","ports":["7050:7050"]},
    {"ip":"10.0.2.15","id":"1","orgId":"ord1","ports":["6050:6050"]},
    {"ip":"10.0.2.15","id":"2","orgId":"ord1","ports":["5050:5050"]}
  ],
  "peers": [
    {"ip":"10.0.2.15","id":"0","orgId":"test","ports":["7051:7051"]},
    {"ip":"10.0.2.15","id":"1","orgId":"test","ports":["8051:8051"]}
  ]
}
```

#### 配置文件修改

根据下面配置文件参数解释和实际服务器情况修改

```bash
domain： 生成证书的后缀名，推荐用默认值
sshUserName： ssh默认登陆用户名，"每个机器也可配置自定义用户名"
sshPwd： ssh默认登陆密码，"每个机器也可配置自定义密码"
sshKey： ssh默认登陆私钥(容器内位置),需要宿主机映射到容器中
sshPort: ssh默认登陆端口，"每个机器也可配置自定义密码"
ccInit: 智能合约初始化参数
ccName： 智能合约名称
ccVersion： 智能合约版本， 升级时要修改，必须为整数1，2，3..., 而且第一次必须从1开始
ccPath： 智能合约源码路径或包绝对路径(容器内位置), 
    已内置智能合约：存证类:'github.com/chaincode/commoncc'
testArgs： 执行测试智能合约的参数
chan_counts： 创建的业务通道个数，默认为1对应通道"mychannel" 修改后为"mychannel2" ...
mountPath: orderer和peer节点账本数据挂载的宿主机位置，默认"/data"eg:/data/peer0.org1.irchain.net
imagePre： 镜像前缀,   eg：  "gmhyperledger"
imageTag: 镜像标签， eg: "2.0.0-gm" 
log: orderer和peer日志级别， eg: "INFO"
batchTime、batchSize、batchPreferred: 切块的条件
orderers： 对应orderer节点数组
ip: 服务器ip
id： 当前节点序列号
orgId: 当前节点归属组织组织名
ports： 当前节点端口映射数组列表， eg: ["7050:7050"] , 前面为外部访问端口
peers: 对应peers节点数组， 和orderer解释一样
```

#### 修改配置文件权限

修改宿主机node.json 权限为777，否则修改宿主机上的文件会引起内容不同步的问题

```bash
sudo chmod 777 ~/deploy/node.json   
```

### 启动部署工具

先安装部署工具镜像

```bash
docker pull gmhyperledger/deploy-tool:latest
```

chaincode映射路径要与node.json配置的ccpath相同

```bash
docker run -it -d --name manager -e GODEBUG=netdns=go -v $PWD/config:/opt/deployFabricTool/config -v $PWD/node.json:/opt/deployFabricTool/data/node.json -v $PWD/chaincode:/go/src/github.com/user/chaincode gmhyperledger/deploy-tool:latest
```

如果报错，需要根据具体错误调整

### --下面一定要按顺序执行--

### 检测配置文件登录参数是否配置正确

```bash
docker exec manager bash -c ./0-checknode.sh
```

如果连接失败或卡住，就要检测node.json里面的ip和ssh配置，或者用ssh命令先测试一下

### 生成fabric需要资源文件

需要生成新的crypto-config 证书目录

```bash
docker exec manager bash -c ./1-makeConfig.sh
```

或

用已经存在的crypto-config证书目录，需要在启动容器前将证书放在config/crypto-config目录下

```bash
docker exec manager bash -c './1-makeConfig.sh 1'
```

### 启动fabric节点

```bash
docker exec manager bash -c ./2-startNode.sh
```

### 创建通道并部署

部署基础链-----通道名（basechannel),  智能合约名（basecc)

```bash
docker exec manager bash -c './multiChannel.sh basechannel basecc'
```

PS: 可能要等待一段时间

部署业务链-----通道名（tranchannel),  智能合约名（trancc) 

```bash
docker exec manager bash -c './multiChannel.sh tranchannel trancc'
```

PS: 可能要等待一段时间

### 检查所有节点和智能合约是否启动成功

```bash
docker exec manager bash -c ./0-checknode.sh
```

如果上面部署命令执行有错误，先根据错误日志判定是否node.json文件参数错误，然后要清理环境再重新部署。

### 发交易测试智能合约是否部署成功

修改node.json文件中的testArgs对应的参数

```bash
docker exec manager bash -c './invokecc.sh basechannel basecc'
```

### 查询指定节点已加入通道

参数1： 操作方法名， 参数2：节点名称

```bash
docker exec manager bash -c './deployFabricTool -r chanlist -nodename peer0.test.irchain.net'
```

### 新增已存在组织peer节点

#### 1. 添加正确节点信息

修改node.json文件在peers里面添加正确节点信息

```json
{"ip":"10.0.2.15","id":"0","orgId":"test","ports":["7051:7051"]},
//新增 peer1.test.irchain.net 节点
{"ip":"10.0.2.15","id":"1","orgId":"test","ports":["6051:7051"]} 
```

#### 2. 生成节点证书文件并启动

脚本后参数1： 组织名， 参数2：节点名

```bash
docker exec manager bash -c './newnodeadd.sh test peer1.test.irchain.net'
```

#### 3. 新节点加入通道

脚本后参数1： 通道名， 参数2：节点名

```bash
docker exec manager bash -c './newnodejoinchann.sh basechannel peer1.test.irchain.net'
```

#### 4. 新节点部署智能合约

脚本后参数1： 通道名， 参数2：合约名  参数3：当前智能合约版本 参数3：节点名

```bash
docker exec manager bash -c './newnodeinstallcc.sh basechannel basecc 1 peer1.test.irchain.net'
```

### 添加Peer组织

#### 1. 添加正确组织节点信息

修改node.json文件在peers里面添加正确节点信息

```json
//新增 peer0.test3.irchain.net 节点
{"ip":"10.0.2.15","id":"0","orgId":"test3","ports":["7051:7051"]},
```

#### 2. 创建新组织配置文件

脚本参数1： 组织名

```bash
docker exec manager bash -c './neworgfilecreate.sh test3'
```

#### 3. 更新组织配置到指定通道

脚本参数1： 组织名 参数2：通道名

```bash
docker exec manager bash -c './neworgconfigupdate.sh test3 basechannel'
```

#### 4. 添加新组织节点

参照 上一节： 《新增已存在组织peer节点》

### 删除原有Peer组织

执行如下命令：参数 -orgid （要删除组织id)     -n （通道名）

```bash
docker exec manager bash -c './deployFabricTool -r rmorgfromconfigblock -orgid test2 -n mychannel'
```

### 删除原有Peer节点

执行如下命令：参数 -nodename （要删除节点名)   注意该命令不会删除挂载账本目录

```bash
docker exec manager bash -c './deployFabricTool -r rmnode -nodename peer0.test2.irchain.net'
```

### 添加新Orderer节点

#### 1.修改node.json文件在orderers里面添加正确节点信息

```json
//新增 orderer3.ord.irchain.net 节点，注意按照json格式上一行结尾加逗号
{"ip":"10.0.2.15","id":"3","orgId":"ord","ports":["9050:7050"]}
```

#### 2.更新原有orderer节点的域名列表

PS:  添加新orderer节点前需要在原有节点的域名映射列表里面增加新orderer节点连接方式

脚本后参数1： 操作方法名， 参数2：节点名    

```bash
docker exec manager bash -c './deployFabricTool -r updatenodedomain -nodename orderer0.ord.irchain.net'
```

```bash
docker exec manager bash -c './deployFabricTool -r updatenodedomain -nodename orderer1.ord.irchain.net'
```

```bash
docker exec manager bash -c './deployFabricTool -r updatenodedomain -nodename orderer2.ord.irchain.net'
```

#### 3.生成新orderer节点证书

脚本后参数1： 操作方法名， 参数2：节点所在组织名（orderer节点的都是ord)

```bash
docker exec manager bash -c './deployFabricTool -r addorgnodecert -orgid ord'
```

#### 3.添加新节点到系统通道

PS: 更新业务通道前，必须先更新系统通道

脚本后参数1： 操作方法名， 参数2：节点名    ， 参数3： 通道名

```bash
docker exec manager bash -c './deployFabricTool -r addordertoconfigblock -nodename orderer3.ord.irchain.net -n byfn-sys-channel'
```

#### 4.添加新节点到业务通道

PS: 更新业务通道前，必须先更新系统通道

脚本后参数1： 操作命， 参数2：节点名    ， 参数3： 通道名

```bash
docker exec manager bash -c './deployFabricTool -r addordertoconfigblock -nodename orderer3.ord.irchain.net -n basechannel'
```

#### 5.更新新orderer节点启动依赖的创世区块

PS: 必须用最新的系统通道配置块，作为新orderer节点的创世块

```bash
docker exec manager bash -c './deployFabricTool -r updategenesisblock'
```

#### 6.生成orderer节点yaml文件并启动

脚本后参数1： 组织名， 参数2：节点名

```bash
docker exec manager bash -c './newordereradd.sh ord orderer3.ord.irchain.net'
```

#### 7. 确认新orderer节点加入通道成功

PS: 等待一段时间，新orderer节点需要同步之前的区块

```bash
docker logs -f orderer3.ord.irchain.net --tail 1000 2>&1
```

PS： 执行上面命令如果新orderer节点最后写的区块号为当前网络最新区块号，则说明新orderer加入集群成功。

### 删除Orderer节点

#### 1.删除orderer节点从系统通道

PS: 更新业务通道前，必须先更新系统通道

脚本后参数1： 操作方法名， 参数2：节点名    ， 参数3： 通道名

```bash
docker exec manager bash -c './deployFabricTool -r rmorderfromconfigblock -nodename orderer3.ord.irchain.net -n byfn-sys-channel'
```

#### 2.删除orderer节点从业务通道

PS: 更新业务通道前，必须先更新系统通道

脚本后参数1： 操作命， 参数2：节点名    ， 参数3： 通道名

```bash
docker exec manager bash -c './deployFabricTool -r rmorderfromconfigblock -nodename orderer3.ord.irchain.net -n basechannel'
```

#### 3.删除orderer节点容器

执行如下命令：参数 -nodename （要删除节点名)   注意该命令不会删除挂载账本目录

```bash
docker exec manager bash -c './deployFabricTool -r rmnode -nodename orderer3.ord.irchain.net'
```

## 后台客户端所需证书目录

```bash
~/depoly/config/crypto-config
```

## 重新部署-清理环境

清理fabric

```bash
docker exec manager bash -c ./removenode.sh
```

然后再重头开始执行

## 升级或单独安装fabric智能合约

替换新版本chaincode源码到~/deploy/chaincode目录下

参数解释： 通道名（basechannel),  智能合约名（basecc) , 智能合约版本（2）

```bash
docker exec manager bash -c './godeploycc.sh basechannel basecc 2'
```

## 其他功能

### 查看所有命令和参数

```bash
docker exec manager bash -c ./deployFabricTool
```

