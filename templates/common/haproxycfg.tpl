frontend web #自定义前端服务器名，这里为web
    bind *:3333 #设置监听的端口与IP
    default_backend     websrvs #配置默认调用的名为websrvs的后端服务器组

backend websrvs #定义后端服务器名
    balance roundrobin #配置算法类型
    server srv1 10.10.134.12:5555 check #配置主机参数，check为健康状态检测
    server srv2 10.10.134.12:6666 check
