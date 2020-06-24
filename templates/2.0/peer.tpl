version: '2'

services:
  peer{{.id}}.org{{.orgId}}.{{.domain}}:
    image: {{.imagePre}}/fabric-peer:{{.imageTag}}
    restart: always
    container_name: peer{{.id}}.org{{.orgId}}.{{.domain}}
    environment:
      # base env
      - GODEBUG=netdns=go
      - BCCSP_CRYPTO_TYPE={{.cryptoType}}
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_network
      - FABRIC_LOGGING_SPEC=gossip=warning:msp=warning:grpc=warning:leveldbhelper=warning:comm.grpc.server=warning:{{.log}}
      - CORE_CHAINCODE_LOGGING_LEVEL={{.log}}
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
      - CORE_CHAINCODE_EXECUTETIMEOUT=300s
      {{if eq .imagePre "peersafes"}}
      - CORE_CHAINCODE_BUILDER=peersafes/fabric-ccenv:{{.imageTag}}
      - CORE_CHAINCODE_GOLANG_RUNTIME=peersafes/fabric-baseos:{{.imageTag}}{{end}}
      # improve env
      - CORE_PEER_ID=peer{{.id}}.org{{.orgId}}.{{.domain}}
      #peer listen service
      - CORE_PEER_LISTENADDRESS=0.0.0.0:7051
      #for same org peer connect
      - CORE_PEER_ADDRESS=peer{{.id}}.org{{.orgId}}.{{.domain}}:{{.externalPort}}
      #listen chaincode connect
      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
      #chaincode connect to peer
      - CORE_PEER_CHAINCODEADDRESS=peer{{.id}}.org{{.orgId}}.{{.domain}}:7052
      #connect peer when peer init{{if ne .bootStrapAddress ""}}
      - CORE_PEER_GOSSIP_BOOTSTRAP={{.bootStrapAddress}}{{else}}
      - CORE_PEER_GOSSIP_BOOTSTRAP=127.0.0.1:7051{{end}}
      #for other org peer connect
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer{{.id}}.org{{.orgId}}.{{.domain}}:{{.externalPort}}
      - CORE_OPERATIONS_LISTENADDRESS=0.0.0.0:9443
      - CORE_PEER_LOCALMSPID=Org{{.orgId}}MSP
      {{if eq .useCouchdb "true"}}
      - CORE_LEDGER_STATE_STATEDATABASE=CouchDB
      - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb:5984{{end}}
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
        - /var/run/:/host/var/run/
        - ../crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/peers/peer{{.id}}.org{{.orgId}}.{{.domain}}/msp:/etc/hyperledger/fabric/msp
        - ../crypto-config/peerOrganizations/org{{.orgId}}.{{.domain}}/peers/peer{{.id}}.org{{.orgId}}.{{.domain}}/tls:/etc/hyperledger/fabric/tls
        - {{.mountPath}}/peer{{.id}}.org{{.orgId}}.{{.domain}}:/var/hyperledger/production
    networks:
      - outside
    logging:
      driver: "json-file"
      options:
        max-size: "200m"
        max-file: "50"
    ports:{{range $index,$value:= .ports}}
      - {{$value}}{{end}}
    {{if gt (len .extra_hosts) 0}}extra_hosts:{{range $index,$value:= .extra_hosts}}
      - "{{$value.domain}}:{{$value.ip}}"{{end}}{{end}}
  {{if eq .useCouchdb "true"}}
    depends_on:
      - couchdb

  couchdb:
    container_name: couchdb
    image: couchdb:2.3
    ports:
       - "5984:5984"
    volumes:
       - {{.mountPath}}/couchdb:/opt/couchdb/data
    networks:
      - outside
   {{end}}

networks:
  outside:
    external:
      name: fabric_network


