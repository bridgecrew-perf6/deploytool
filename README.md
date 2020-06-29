# 北京众享比特科技有限公司

## 区块链层部署文档

## 检查部署所需资源文件

chaincode目录： 包含智能合约源码

node.json 文件：节点服务器配置文件

## 根据需求修改配置文件node.json

下面为模板文件： 采用fabric1.4版本， etcdraft 共识类型，3个orderer，2个peer 节点

*具体参数需根据实际数据进行调整*，主要关注“xxx” 这些参数修改，ip写内网ip

```json
{
  "fabricVersion":"1.4","domain":"example.com","cryptoType":"FGM",
  "sshUserName":"peersafe","sshPwd":"dev1@peersafe","sshPort":"22","sshKey":"/etc/login.pem",
  "ccInit":"'{\"Args\":[\"init\"\\,\"a\"\\,\"100\"\\,\"b\"\\,\"200\"]}'",
  "ccPolicy":"\"OR  ('Org1MSP.member'\\,'Org2MSP.member')\"",
  "ccName":"mycc","ccVersion":"1","ccInstallType":"path",
  "testArgs":"'{\"Args\":[\"invoke\"\\,\"a\"\\,\"b\"\\,\"1\"]}'",
  "ccPath":"github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd",
  "chan_counts":1,"mountPath": "/data", "caType": "cryptogen",
  "consensusType":"raft", "imagePre":"hyperledger","imageTag":"1.4.3","log":"INFO",
  "batchTime":"2s", "batchSize":100, "batchPreferred":"512 KB", "useCouchdb":"false",
  "orderers":[
    {"ip":"XXX","id":"0","orgId":"1","ports":["7050:7050","5443:9443"]},
    {"ip":"XXX","id":"1","orgId":"1","ports":["8050:7050","6443:9443"]},
    {"ip":"XXX","id":"2","orgId":"1","ports":["9050:7050","7443:9443"]}
  ],
  "peers": [
    {"ip":"XXX","id":"0","orgId":"1","ports":["7051:7051","8443:9443"]},
    {"ip":"XXX","id":"1","orgId":"1","ports":["8051:7051","9443:9443"]},
    {"ip":"XXX","id":"0","orgId":"2","ports":["9051:7051","10443:9443"]},
    {"ip":"XXX","id":"1","orgId":"2","ports":["10051:7051","11443:9443"]}
  ],
  "cas": [
    {"ip":"XXX","certType":"orderer","orgId":"1","ports":["7054:7054","9543:9443"]},
    {"ip":"XXX","certType":"peer","orgId":"1","ports":["8054:7054","9643:9443"]},
    {"ip":"XXX","certType":"peer","orgId":"2","ports":["9054:7054","9743:9443"]}
  ],
  "apis": [
    {"ip":"XXX","orgId":"1","imageTag":"latest","apiPort":"5984"}
  ],
  "explorers": [
    {"ip":"XXX","explorerIp":"XXX","webPort":"7054","apiPort":"8080",
      "peerId":"0","orgId":"1","imageTag":"onlyblock","chList": "[\"mychannel\"]"}
  ],
  "zookeepers": [
  ],
  "kafkas": [
  ]
}
```

## 部署

### 创建部署目录

```bash
mkdir -p ~/deploy && cd  ~/deploy
```

### 复制资源文件 

复制所需的资源文件（chaincode源码目录和 node.json 配置文件) 到deploy目录下

### 修改文件权限

修改宿主机node.json 权限为777，否则修改宿主机上的文件会引起内容不同步的问题

```bash
sudo chmod 777 ~/deploy/node.json   
```

### 配置文件修改

node.json 参数解释

```bash
fabricVersion： fabric的版本，容器内bin目录下可执行文件和tpl模板文件对应的版本,支持"1.4"和"2.0"
domain： 生成证书的后缀名，推荐用默认值
cryptoType： 算法类型, "GM" 国密、"FGM" 非国密
sshUserName： ssh默认登陆用户名，"每个机器也可配置自定义用户名"
sshPwd： ssh默认登陆密码，"每个机器也可配置自定义密码"
sshKey： ssh默认登陆私钥(容器内位置),需要宿主机映射到容器中
sshPort: ssh默认登陆端口，"每个机器也可配置自定义密码"
ccInit: 智能合约初始化参数
ccPolicy： 智能合约背书策略
ccName： 智能合约名称
ccVersion： 智能合约版本， 升级时要修改，必须为整数1，2，3..., 而且第一次必须从1开始
ccPath： 智能合约源码路径或包绝对路径(容器内位置), 
 内置2个智能合约：
 	转账cc 'github.com/hyperledger/fabric/examples/chaincode/go/example02/cmd'
    poc写/读数据cc：'peersafe/fabric_poc/chaincode'
ccInstallType：智能合约安装方式， "path" 源码路径方式 "pkg" 包安装方式
testArgs： 执行调用智能合约的参数
caType: 证书生成方式， "cryptogen","fabric-ca" 方式，默认为 cryptogen
chan_counts： 创建的业务通道个数，默认为1对应通道"mychannel" 修改后为"mychannel2" ...
mountPath: orderer和peer节点账本数据挂载的宿主机位置，默认"/data"eg:/data/peer0.org1.example.com
consensusType: 共识方式，"raft"、"solo"、"kafka"  目前实现只raft
imagePre： 镜像前缀,   eg：  "peersafes"、"hyperledger"
imageTag: 镜像标签， eg: "1.4"、"1.4.3"、"1.4.3-gm"、"2.1.0" 
log: orderer和peer日志级别， eg: "INFO"、"DEBUG"
batchTime、batchSize、batchPreferred: 切块的条件
orderers： 对应orderer节点数组
ip: 服务器ip
id： 当前节点序列号
orgId: 当前节点归属组织id
ports： 当前节点端口映射数组列表， eg: ["8050:7050","10443:9443"] , 前面为外部访问端口
peers: 对应peers节点数组， 和orderer解释一样
certType: 表示ca对应的组织类型,eg: "orderer"、"peer"
explorers: 相关介绍
	ip: 要部署的机器ip
	peerId,orgId：表示event连接的peer
	imageTag： 为镜像标签
	chList: 为channel数组,多个用逗号隔开
	apiPort: 为后台服务端口，默认8888
	webPort: 为前端页面端口，默认3004
apis: 
	ip： 要部署的机器ip
	orgId: 表示客户端所属组织id
	imageTag: 镜像标签,镜像名：fabric-poc-apiserver
    apiPort: api对外端口
```

### 启动部署工具

*这个过程需要拉取部署工具镜像文件,第一次执行需要等待几分钟，以下命令为一行*

chaincode映射路径要与node.json配置的ccpath相同

```bash
docker run -it -d --name manager -v $PWD/config:/opt/deployFabricTool/config -v $PWD/node.json:/opt/deployFabricTool/data/node.json -v $PWD/chaincode:/opt/gopath/src/github.com/peersafe/xxx/chaincode peersafes/deploy-tool:latest
```

如果报错，需要根据具体错误调整

### -下面一定要按顺序执行-

### 检测配置文件登录参数是否配置正确

```bash
docker exec manager bash -c ./0-checknode.sh
```

如果连接失败或卡住，就要检测node.json里面的ip和ssh配置，或者用ssh命令先测试一下

### 生成fabric需要配置文件

需要生成新的crypto-config 证书目录

```bash
docker exec manager bash -c ./1-makeConfig.sh
```

或

用已经存在的crypto-config证书目录，需要在启动容器签将证书放在config/crypto-config目录下

```bash
docker exec manager bash -c './1-makeConfig.sh 1'
```

### 启动fabric节点

```bash
docker exec manager bash -c ./2-startNode.sh
```

### 启动fabric智能合约

```bash
docker exec manager bash -c ./3-runChaincode.sh
```

PS: 可能要等待一段时间

### 检查所有节点是否启动成功

```bash
docker exec manager bash -c ./0-checknode.sh
```

如果上面部署命令执行有错误，先根据错误日志判定是否node.json文件参数错误，然后要清理环境再重新部署。

### 发交易测试CC

修改node.json文件中的testArgs对应得参数

```bash
docker exec manager bash -c ./invokecc.sh
```

## 部署-apiserver

```bash
docker exec manager bash -c ./apistart.sh
```

## 关闭-apiserver

```bash
docker exec manager bash -c ./apidown.sh
```

## 部署-浏览器

```bash
docker exec manager bash -c ./explorerstart.sh
```

## 关闭-浏览器

```bash
docker exec manager bash -c ./explorerdown.sh
```

## 后台APi服务客户端所需证书目录

```bash
~/depoly/config/crypto-config
```

## 重新部署-清理环境

如果部署了apiserver和浏览器一定要先关闭apiserver和浏览器

```bash
docker exec manager bash -c ./apidown.sh
docker exec manager bash -c ./explorerdown.sh
```

然后再清理fabric

```bash
docker exec manager bash -c ./removenode.sh
```

然后再重头开始执行

## 业务变动升级fabric智能合约

替换新版本chaincode源码或包文件

修改node.json文件里面的ccVersion为新版本号：  "ccVersion":"1.1"

```bash
docker exec manager bash -c './godeploycc.sh upgrade'
```

## 安装自定义智能合约

替换新版本chaincode源码或包文件

修改node.json文件里面的ccName、ccVersion、ccInit、ccPolicy、ccInstallType、ccPath

执行安装和实例化智能合约命令

```bash
docker exec manager bash -c ./godeploycc.sh
```

### 隐藏功能

#### 1. 自定义bin、tpl， 比如需要1.3版本的bin和tpl, 先自己准备好bin和tpl

```bash
docker run -it -d --name manager -v $PWD/config:/opt/deployFabricTool/config -v $PWD/node.json:/opt/deployFabricTool/data/node.json -v $PWD/chaincode:/opt/gopath/src/github.com/peersafe/xxx/chaincode -v
$PWD/bin/1.3/:/opt/deployFabricTool/bin/1.3/ -v
$PWD/templates/1.3/:/opt/deployFabricTool/templates/1.3/ -v
peersafes/deploy-tool:latest
```

#### 2. 单独执行拆分命令

先看原有包含哪些命令

```bash
docker exec manager bash -c 'cat 1-makeConfig.sh'
docker exec manager bash -c 'cat 2-startNode.sh'
```

只执行其中包含的某条命令

```bash
docker exec manager bash -c './deployFabricTool -r runchaincode -n mychannel'
```

#### 3.查看所有命令和参数

```bash
docker exec manager bash -c ./deployFabricTool
```

## 特别说明：

以上部署需要依赖外网环境

如果想要部署在内网环境机器，可以先用外网机器拉取所需镜像，在将镜像导入到内网服务器

```bash
docker pull peersafes/deploy-tool:latest   		#部署工具镜像
docker pull peersafes/fabric-orderer:XXX		#orderer镜像
docker pull peersafes/fabric-peer:XXX		#peer镜像
docker pull peersafes/fabric-ccenv:XXX		#编译智能合约依赖镜像
docker pull peersafes/fabric-baseos:XXX		#智能合约运行镜像
```



