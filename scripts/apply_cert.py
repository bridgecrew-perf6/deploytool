#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import utils
from fabric.api import local, lcd, put, run, cd

reload(sys)
sys.setdefaultencoding('utf8')

def generate_genesis_block(model, bin_path, cfg_path ,out_path, crypto_type):
    if not os.path.exists(out_path + "crypto-config"):
        with lcd(out_path):
            local("tar -zxvf crypto-config.tar.gz")
    if not os.path.exists(cfg_path + "core.yaml"):
        local("cp %s/core.yaml %s"%(bin_path, cfg_path))

    tool = utils.get_bin_path(bin_path, "configtxgen",crypto_type)
    channel_path = out_path + "channel-artifacts"
    local("rm -rf %s"%channel_path)
    local("mkdir -p %s"%channel_path)
    env = "FABRIC_CFG_PATH=%s"%cfg_path
    local("%s %s -profile %s -outputBlock %s/genesis.block"%(env,tool,model,channel_path))
    with lcd(out_path):
        local("tar -zcvf channel-artifacts.tar.gz channel-artifacts")

## Generates orderer Org certs using cryptogen tool
def generate_certs(bin_path, cfg_path ,out_path, crypto_type):
    cryptotool = utils.get_bin_path(bin_path, "cryptogen",crypto_type)
    yamlfile =  cfg_path + "crypto-config.yaml"
    mm_path = out_path + "crypto-config"

    with lcd(out_path):
        local("rm -rf crypto-config.tar.gz crypto-config")
    local("%s generate --config=%s --output='%s'"%(cryptotool,yamlfile,mm_path))
    with lcd(out_path):
        local("tar -zcf crypto-config.tar.gz crypto-config")


def put_cryptoconfig(config_path, type, node_name, org_name):
    run("mkdir -p ~/deployFabricTool")
    with lcd(config_path):
        if type == "orderer":
            local('tar -zcvf %s_crypto-config.tar.gz crypto-config/ordererOrganizations/%s/orderers/%s'%(node_name,org_name,node_name))
            copy_file(config_path,"%s_crypto-config.tar.gz"%node_name)
            copy_file(config_path,"channel-artifacts.tar.gz")
            # copy_file(config_path,"kafkaTLSclient.tar.gz")
        elif type == "kafka":
            copy_file(config_path,"kafkaTLSserver.tar.gz")
        elif type == "peer":
            local('tar -zcvf %s_crypto-config.tar.gz crypto-config/peerOrganizations/%s/peers/%s'%(node_name,org_name,node_name))
            copy_file(config_path,"%s_crypto-config.tar.gz"%node_name)
        # elif type == "api":
        #     copy_file(config_path,"crypto-config.tar.gz")

def copy_file(config_path, file_name):
    remote_file = "~/deployFabricTool/%s"%file_name
    if utils.check_remote_file_exist(remote_file) == "false":
        put("%s%s"%(config_path,file_name), "~/deployFabricTool/")
        local("rm -rf %s%s"%(config_path,file_name))
        with cd("~/deployFabricTool"):
            run("tar zxfm %s"%file_name)
            run("rm %s"%file_name)

