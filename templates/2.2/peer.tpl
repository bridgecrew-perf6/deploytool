version: '2'

services:
  peer{{.id}}.{{.orgId}}.{{.domain}}:
    image: {{.imagePre}}/fabric-peer:{{.imageTag}}
    restart: always
    container_name: peer{{.id}}.{{.orgId}}.{{.domain}}
    environment:
      # base env
      - GODEBUG=netdns=go
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_network
      - FABRIC_LOGGING_SPEC={{.log}}
      {{if or (eq .log "info") (eq .log "INFO")}}
      - FABRIC_LOGGING_SPEC=gossip=warning:msp=warning:grpc=warning:leveldbhelper=warning:comm.grpc.server=warning:{{.log}}
      {{end}}
      - CORE_CHAINCODE_LOGGING_LEVEL={{.log}}
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_GOSSIP_USELEADERELECTION=false
      - CORE_PEER_GOSSIP_ORGLEADER=true
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_CHAINCODE_EXECUTETIMEOUT=300s
      - CORE_CHAINCODE_BUILDER={{.imagePre}}/fabric-ccenv:{{.imageTag}}
      - CORE_CHAINCODE_GOLANG_RUNTIME={{.imagePre}}/fabric-baseos:{{.imageTag}}
      # improve env
      - CORE_PEER_ID=peer{{.id}}.{{.orgId}}.{{.domain}}
      #peer listen service
      - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
      #for same org peer connect
      - CORE_PEER_ADDRESS=peer{{.id}}.{{.orgId}}.{{.domain}}:{{.externalPort}}
      #listen chaincode connect
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
      #chaincode connect to peer
      - CORE_PEER_CHAINCODEADDRESS=peer{{.id}}.{{.orgId}}.{{.domain}}:7052
      #for other org peer connect
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer{{.id}}.{{.orgId}}.{{.domain}}:{{.externalPort}}
      - CORE_OPERATIONS_LISTENADDRESS=0.0.0.0:9443
      - CORE_PEER_LOCALMSPID={{.orgId}}
      #connect peer when peer init
      {{if ne .bootStrapAddress ""}}
      - CORE_PEER_GOSSIP_BOOTSTRAP={{.bootStrapAddress}}{{else}}
      - CORE_PEER_GOSSIP_BOOTSTRAP=127.0.0.1:{{.externalPort}}{{end}}
     # - CORE_PEER_GOSSIP_BOOTSTRAP=127.0.0.1:{{.externalPort}}
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
        - /etc/timezone:/etc/timezone
        - /etc/localtime:/etc/localtime
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/{{.orgId}}.{{.domain}}/peers/peer{{.id}}.{{.orgId}}.{{.domain}}/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/{{.orgId}}.{{.domain}}/peers/peer{{.id}}.{{.orgId}}.{{.domain}}/tls:/etc/hyperledger/fabric/tls
        - {{.mountPath}}/peer{{.id}}.{{.orgId}}.{{.domain}}:/var/hyperledger/production
    networks:
      - outside
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "5"
    ports:{{range $index,$value:= .ports}}
      - {{$value}}{{end}}
    {{if gt (len .extra_hosts) 0}}extra_hosts:{{range $index,$value:= .extra_hosts}}
      - "{{$value.domain}}:{{$value.ip}}"{{end}}{{end}}

networks:
  outside:
    external:
      name: fabric_network



