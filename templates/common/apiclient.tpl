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

orderers:{{range $value:= .ordList}}
  orderer{{$value.serverId}}:
    host:       {{$value.serverHost}}
    domainName: {{$value.serverDomain}}
    useTLS:     true
    tlsPath:    {{$value.serverTlsPath}}{{end}}
peers:{{range $index,$value:= .peerList}}
  peer{{$index}}:
    host:         {{$value.serverHost}}
    domainName:   {{$value.serverDomain}}
    useTLS:       true
    tlsPath:      {{$value.serverTlsPath}}{{end}}

eventPeers:{{range $index,$value:= .eventPeerList}}
  peer{{$index}}:
    host:         {{$value.serverHost}}
    domainName:   {{$value.serverDomain}}
    useTLS:       true
    tlsPath:      {{$value.serverTlsPath}}{{end}}

channel:
    mspConfigPath: ./crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/users/Admin@org{{.orgId}}.{{.domain}}/msp
    localMspId:          Org{{.orgId}}MSP
    channelId:           mychannel
    chaincodeName:       {{.ccName}}

log:
    logLevel: {{.log}}





