version: '2'

services:
  orderer{{.id}}.ord{{.orgId}}.{{.domain}}:
    container_name: orderer{{.id}}.ord{{.orgId}}.{{.domain}}
    image: {{.imagePre}}/fabric-orderer:{{.imageTag}}
    restart: always
    environment:
      - GODEBUG=netdns=go
      - FABRIC_LOGGING_SPEC={{.log}}
      {{if eq .log "info"}}
      - FABRIC_LOGGING_SPEC=orderer.common.cluster.step=info:orderer.consensus.etcdraft=info:{{.log}}
      {{end}}
      - ORDERER_GENERAL_LISTENPORT={{.externalPort}}
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block
      - ORDERER_GENERAL_LOCALMSPID=Orderer{{.orgId}}MSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_OPERATIONS_LISTENADDRESS=0.0.0.0:9443
      # enabled TLS
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
      - ../crypto-config/ordererOrganizations/ord{{.orgId}}.{{.domain}}/orderers/orderer{{.id}}.ord{{.orgId}}.{{.domain}}/msp:/var/hyperledger/orderer/msp
      - ../crypto-config/ordererOrganizations/ord{{.orgId}}.{{.domain}}/orderers/orderer{{.id}}.ord{{.orgId}}.{{.domain}}/tls:/var/hyperledger/orderer/tls
      - {{.mountPath}}/orderer{{.id}}.ord{{.orgId}}.{{.domain}}:/var/hyperledger/production
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "5"
    networks:
      - outside
    ports:{{range $index,$value:= .ports}}
      - {{$value}}{{end}}
    {{if gt (len .extra_hosts) 0}}extra_hosts:{{range $index,$value:= .extra_hosts}}
      - "{{$value.domain}}:{{$value.ip}}"{{end}}{{end}}

networks:
  outside:
    external:
      name: fabric_network




