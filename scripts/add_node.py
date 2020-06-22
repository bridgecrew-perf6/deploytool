#!/bin/python

from fabric.api import cd, put, lcd, local, run, settings, sudo
import sys
import os
import utils

reload(sys)
sys.setdefaultencoding('utf8')


def load_images(type, images_path):
    filter = type
    if type == "ca":
        filter = "fabric-ca"
    result = run('unset GREP_OPTIONS && docker images | grep -e "%s" | wc -l' % filter)
    if result == "0":
        with settings(warn_only=True):
            run("mkdir -p ~/images")
        print "check local image package is exist"
        local("ls %s/%s.tar.gz" % (images_path, type))
        put("%s/%s.tar.gz" % (images_path, type), "~/images/")
        with cd("~/images/"):
            # load image
            run("tar zxvfm %s.tar.gz" % type)
            run("rm %s.tar.gz" % type)
            run("docker load -i %s.tar" % type)
    else:
        sys.stdout.write("%s image is exsit" % type)


def replace_images(type, config_path):
    if type == "order":
        type = "orderer"
    local("docker save hyperledger/fabric-%s:latest -o %s%s.tar" % (type, config_path, type))
    run("rm -rf ~/%s.tar" % type)
    put("%s%s.tar" % (config_path, type), "~")
    sudo("systemctl restart docker")
    run("docker rmi hyperledger/fabric-%s:latest" % type)
    run("docker load -i ~/%s.tar" % type)
    sudo("systemctl restart docker")


def up_node(yaml_name):
    # start container
    with cd("~/fabricNetwork/yaml/"):
        run("docker-compose -f %s.yaml up -d" % yaml_name)


def start_node(node_name, config_dir):
    if not utils.check_network_exist("fabric_network"):
        run("docker network create fabric_network")

    if node_name == "block_fabric_explorer":
        run("docker-compose -f ~/fabricNetwork/yaml/block_fabric_explorer.yaml up -d")
        return
    if node_name == "apiserver":
        with lcd(config_dir):
            local("tar -zcvf %s.tar.gz %s.yaml" % ("client_sdk", "client_sdk"))
            run("mkdir -p ~/fabricNetwork/yaml/")
            put("client_sdk.tar.gz", "~/fabricNetwork/yaml/")
            local("rm client_sdk.tar.gz")

    with lcd(config_dir):
        local("tar -zcvf %s.tar.gz %s.yaml" % (node_name, node_name))
        # remote yaml
        run("mkdir -p ~/fabricNetwork/yaml/")
        put("%s.tar.gz" % node_name, "~/fabricNetwork/yaml/")
        local("rm %s.tar.gz" % node_name)
    # start container
    with cd("~/fabricNetwork/yaml/"):
        if node_name == "apiserver":
            run("tar zxvfm client_sdk.tar.gz")
            run("rm client_sdk.tar.gz")

        run("tar zxvfm %s.tar.gz" % node_name)
        run("rm %s.tar.gz" % node_name)
        run("docker-compose -f %s.yaml up -d" % node_name)


def start_docker():
    # start docker service
    sudo("systemctl restart docker")


def start_api(peer_id, org_id, config_dir, api_id):
    name = "peer" + peer_id + "org" + org_id + "api" + api_id
    apiclientname = name + "apiclient"
    apidockername = name + "apidocker"
    parent_path = os.path.dirname(config_dir)
    # apiserver
    with lcd(parent_path):
        # remote yaml
        run("mkdir -p ~/fabricNetwork/yaml//api_server/%s" % name)
        run("rm -rf ~/fabricNetwork/yaml//api_server/%s/*" % name)
        put("api_server.tar.gz", "~/fabricNetwork/yaml//api_server/%s" % name)
    with cd("~/fabricNetwork/yaml//api_server/%s" % name):
        run("tar zxvfm api_server.tar.gz --strip-components=1")
        run("rm -rf api_server.tar.gz")
    with lcd(config_dir):
        put("%s.yaml" % apiclientname, "~/fabricNetwork/yaml/api_server/%s/client_sdk.yaml" % name)
        put("%s.yaml" % apidockername, "~/fabricNetwork/yaml/api_server/%s/docker-compose.yaml" % name)
    with cd("~/fabricNetwork/yaml/api_server/%s" % name):
        run("docker-compose -f docker-compose.yaml down")
        run("docker-compose -f docker-compose.yaml up -d")


def start_event(peer_id, org_id, config_dir, clitype, api_id):
    name = "peer" + peer_id + "org" + org_id + "api" + api_id
    yamlname = name + "%sclient" % clitype
    parent_path = os.path.dirname(config_dir)
    # apiserver or eventserver
    with lcd(parent_path):
        put("/etc/hosts", "~")
        # remote yaml
        run("mkdir -p ~/fabricNetwork/yaml/%s_server/%s" % (clitype, name))
        run("rm -rf ~/fabricNetwork/yaml/%s_server/%s/*" % (clitype, name))
        put("%s_server.tar.gz" % clitype, "~/fabricNetwork/yaml/%s_server/%s" % (clitype, name))
    # utils.kill_process("%sserver"%clitype)
    with cd("~/fabricNetwork/yaml/%s_server/%s" % (clitype, name)):
        run("tar zxvfm %s_server.tar.gz --strip-components=1" % clitype)
        run("rm %s_server.tar.gz" % clitype)
    with lcd(config_dir):
        put("%s.yaml" % yamlname, "~/fabricNetwork/yaml/%s_server/%s/client_sdk.yaml" % (clitype, name))
    with cd("~/fabricNetwork/yaml/%s_server/%s" % (clitype, name)):
        sudo("cp ~/hosts /etc/hosts")
        run("chmod +x %sserver" % clitype)
        run("rm -rf %sserver.log" % clitype)
        run("$(nohup ./%sserver >> %sserver.log 2>&1 &) && sleep 1" % (clitype, clitype))
        run("cat /dev/null > %sserver.log" % clitype)


def check_node():
    # check container
    run("docker ps")
