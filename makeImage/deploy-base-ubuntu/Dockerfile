FROM peersafes/golang-ubuntu-16.04:1.11.5
MAINTAINER <peersafe>

# update apt from aliyun
RUN mv /etc/apt/sources.list /etc/apt/sources.list.bak
ADD ./sources.list /etc/apt/
RUN apt-get update

## install telnet
RUN apt-get install -y telnet

#install ping
RUN apt-get install -y inetutils-ping

# install wget
RUN apt-get install -y wget \
    curl \
    python \
    python-dev \
    jq \
    vim \
    gcc \
    libtool \
    libffi-dev

RUN curl https://bootstrap.pypa.io/pip/2.7/get-pip.py -o get-pip.py
RUN python get-pip.py

# install pip
RUN pip install --upgrade pip

# install fabric
RUN pip install fabric==1.14.0
RUN pip install paramiko==2.4.2
RUN pip install cryptography==2.4.2

RUN pip install cython