FROM golang:1.11.5-alpine3.9
ENV TZ "Asia/Shanghai"
#RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/main" > /etc/apk/repositories
RUN apk add --update --no-cache  \
    jq \
    g++ \
    gcc \
    busybox-extras \
    bash \
    python \
    python-dev \
    py-pip \
    libxml2-dev \
    libxslt-dev \
    libffi-dev \
    openssl-dev \
    make

RUN mkdir -p /etc/apk/keys
RUN mkdir -p /root/.pip
COPY sgerrand.rsa.pub /etc/apk/keys/sgerrand.rsa.pub
COPY glibc-2.29-r0.apk /etc/glibc-2.29-r0.apk
COPY pip.conf /root/.pip/pip.conf

RUN apk add --no-cache  \
    /etc/glibc-2.29-r0.apk \
  && pip install paramiko==2.4.2 \
  && pip install cryptography==2.4.2 \
  && pip install fabric==1.13.1 \
  && rm -rf /var/cache/apk/*
