#!/usr/bin/env bash

echo "build image"
    docker rmi -f peersafes/deploy-base:ubuntu
    docker build -t peersafes/deploy-base:ubuntu .

