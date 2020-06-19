version: '2'

networks:
  outside:
    external:
      name: fabric_network

services:
  apiserver:
    container_name: apiserver
    image: peersafes/fabric-poc-apiserver
    restart: always
    volumes:
        - ./client_sdk.yaml:/opt/apiserver/client_sdk.yaml
        - ../crypto-config/:/opt/apiserver/crypto-config
        - /etc/localtime:/etc/localtime
    working_dir: /opt/apiserver
    logging:
      driver: "json-file"
      options:
        max-size: "50m"
        max-file: "10"
    command: ./apiserver
    ports:
     - {{.apiPort}}:8888
    networks:
      - outside