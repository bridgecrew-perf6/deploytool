//配置文件
/* eslint-disable */
const prefix = "{{.explorerIp}}"; //该地址需能被web页面所在的ip访问，即配置外部ip
const SERVER_URL = "http://"+prefix + ":{{.apiPort}}/v3"; //项目部署后台地址
const WEBSOCKET_URL = "ws://"+prefix+":{{.apiPort}}/v3/fabric/ws"; //websocket请求地址
