#!/bin/bash

docker rmi peersafes/deploy-base:latest
docker build -t peersafes/deploy-base:latest .
