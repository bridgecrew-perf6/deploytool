FROM golang:1.11.5-alpine3.9
ENV TZ "Asia/Shanghai"
RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.9/main" > /etc/apk/repositories

RUN apk add --update --no-cache  \
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
    make \
  && pip install fabric==1.13.1 \
  && rm -rf /var/cache/apk/*
