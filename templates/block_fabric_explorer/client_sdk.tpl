crypto:
  {{if eq .cryptoType "GM"}}
  family: gm
  algorithm: P256SM2
  hash: SM3
  {{else}}
  family: ecdsa
  algorithm: P256-SHA256
  hash: SHA2-256
  {{end}}
eventPeer:
  host: peer{{.peerId}}.org{{.orgId}}.{{.domain}}:{{.peerPort}}
  useTLS: true
  tlsPath: ./crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/peers/peer{{.peerId}}.org{{.orgId}}.{{.domain}}/tls/server.crt
  clientCert: ./crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/peers/peer{{.peerId}}.org{{.orgId}}.{{.domain}}/msp/signcerts
  clientKey: ./crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/peers/peer{{.peerId}}.org{{.orgId}}.{{.domain}}/msp/keystore
channels:
  mspConfigPath:    ./crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/users/Admin@org{{.orgId}}.{{.domain}}/msp
  localMspId:       Org{{.orgId}}MSP
  channelIds:       {{.chList}}
  chaincodeName:    {{.ccName}}
  chaincodeVersion: 1.0
service:
    listenPort: 8888
    isHttps: false
db:
  address: fabric_explorer_mysql
  port: 3306
  type: mysql
  name: root
  user1:
  password1:
  user2:
  password2:
  user3: root
  password3: root
network:
  zookeeper:
    zk1:
      name: zk1
      host: zk1.blockchain.cn
      ip: 192.168.0.169
  kafka:
    kafka1:
      name: kafka1
      host: kafka1.blockchain.cn
      ip: 192.168.0.169
  orderers:
    orderer0-org1:
      orgName: ord1
      name: orderer0
      host: orderer0.ord1.example.com
      ip: 192.168.0.169
  peers:
    peer0-org1:
      orgName: org1
      name: peer0
      host: peer0.org1.example.com
      ip: 192.168.0.169



