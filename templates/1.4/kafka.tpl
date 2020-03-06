version: '2'

services:
  kafka{{.kfk_id}}:
    image: {{.imagePre}}/fabric-kafka:{{.imageTag}}
    restart: always
    environment:
      - GODEBUG=netdns=go
      - KAFKA_BROKER_ID={{.kfk_id}}
      - KAFKA_ZOOKEEPER_CONNECT=zk0.{{.kfk_domain}}:2181,zk1.{{.kfk_domain}}:2181,zk2.{{.kfk_domain}}:2181
      - KAFKA_MIN_INSYNC_REPLICAS=2
      - KAFKA_DEFAULT_REPLICATION_FACTOR=3
      - KAFKA_MESSAGE_MAX_BYTES=103809024 # 99 * 1024 * 1024 B
      - KAFKA_REPLICA_FETCH_MAX_BYTES=103809024 # 99 * 1024 * 1024 B
      - KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE=false
      - KAFKA_LOG_RETENTION_HOURS=876000
#      - KAFKA_LISTENERS=PLAINTEXT://:8092
#      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://{{.ip}}:8092
     #enable TLS
      - KAFKA_LISTENERS=PLAINTEXT://:8092,SSL://:9092
      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://{{.ip}}:8092,SSL://kafka{{.kfk_id}}:9092
      - KAFKA_SSL_CLIENT_AUTH=required
      - KAFKA_SSL_KEYSTORE_LOCATION=/opt/kafka/ssl/server.keystore.jks
      - KAFKA_SSL_TRUSTSTORE_LOCATION=/opt/kafka/ssl/server.truststore.jks
      - KAFKA_SSL_KEY_PASSWORD=test1234
      - KAFKA_SSL_KEYSTORE_PASSWORD=test1234
      - KAFKA_SSL_TRUSTSTORE_PASSWORD=test1234
      - KAFKA_SSL_KEYSTORE_TYPE=JKS
      - KAFKA_SSL_TRUSTSTORE_TYPE=JKS
      - KAFKA_SSL_ENABLED_PROTOCOLS=TLSv1.2,TLSv1.1,TLSv1
      - KAFKA_SSL_INTER_BROKER_PROTOCOL=SSL
    ports:
      - 9092:9092
      - 8092:8092
    volumes:
      - /etc/localtime:/etc/localtime
      - ../kafkaTlsServer:/opt/kafka/ssl
      - {{.mountPath}}/kafka{{.kfk_id}}:/tmp/kafka-logs
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "10"
    extra_hosts:
      zk0.{{.kfk_domain}}: {{.zk_ip0}}
      zk1.{{.kfk_domain}}: {{.zk_ip1}}
      zk2.{{.kfk_domain}}: {{.zk_ip2}}

