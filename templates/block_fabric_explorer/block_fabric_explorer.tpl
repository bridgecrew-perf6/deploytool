version: '2'

networks:
  outside:
    external:
      name: fabric_network

services:
  web:
    container_name: fabric_explorer_web
    image: peersafes/fabric_explorer_web:{{.imageTag}}
    restart: always
    ports:
      - {{.webPort}}:3004
    volumes:
      - ../block_fabric_explorer/registerApi.js:/var/www/html/wischain_dist/config/registerApi.js
    depends_on:
      - api
      - mysql
    networks:
      - outside

  api:
    container_name: fabric_explorer_api
    image: peersafes/fabric_explorer_api:{{.imageTag}}
    restart: always
    environment:
      - GODEBUG=netdns=go
    volumes:
      - ../block_fabric_explorer/client_sdk.yaml:/opt/explorer-api/client_sdk.yaml
      - ../crypto-config:/opt/explorer-api/crypto-config
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "10"
    working_dir: /opt/explorer-api
    command: ./explorer-api
    networks:
      - outside
    depends_on:
      - mysql
    ports:
      - {{.apiPort}}:8888
    {{if gt (len .extra_hosts) 0}}extra_hosts:{{range $index,$value:= .extra_hosts}}
      - "{{$value.domain}}:{{$value.ip}}"{{end}}{{end}}
  mysql:
    image: mysql:5.7
    container_name: fabric_explorer_mysql
    restart: always
    ports:
      - "3307:3306"
    environment:
      MYSQL_ROOT_PASSWORD: "root"
      MYSQL_USER: 'test'
      MYSQL_PASS: 'test'
      MYSQL_ROOT_HOST: '%'
    volumes:
      - ../block_fabric_explorer/mysqld.cnf:/etc/mysql/mysql.conf.d/mysqld.cnf
      - ../block_fabric_explorer/mysql_init:/docker-entrypoint-initdb.d/
      - {{.mountPath}}/block_fabric_explorer.{{.domain}}:/var/lib/mysql
    ulimits:
      nproc: 65535
      nofile:
        soft: 65535
        hard: 65535
    networks:
      - outside



