#!/usr/bin/env python
# -*- coding: utf-8 -*-

from fabric.api import local, lcd, put, run, cd
import sys
import os
import utils
reload(sys)
sys.setdefaultencoding('utf8')

## create channel

def create_channel(bin_path, yaml_path, out_path, channel_name, org_id,orderer_address,order_tls_path, domain_name, crypto_type):
    if not os.path.exists(yaml_path + "core.yaml"):
        local("cp %s/core.yaml %s"%(bin_path, yaml_path))
    ret = create_channeltx(bin_path, yaml_path, out_path, channel_name,crypto_type)
    print ret
    channeltx_name = channel_name + '.tx'
    msp_path = yaml_path + "crypto-config/peerOrganizations/%s.%s/users/Admin@%s.%s/msp"%(org_id,domain_name,org_id,domain_name)
    channel_dir = out_path + channel_name
    env = 'FABRIC_CFG_PATH=%s '%yaml_path
    env = env + 'CORE_PEER_LOCALMSPID=%s '%org_id
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s  '%msp_path
    bin = utils.get_bin_path(bin_path, "peer",crypto_type)
    param = ' channel create -o %s -t 3000s -c %s -f %s/%s'%(orderer_address, channel_name, channel_dir, channeltx_name)

    tls = ' --tls --cafile %s'%order_tls_path

    command = env + bin + param + tls
    local(command)
    channel_block = channel_name + '.block'
    local('mv %s %s'%(channel_block, channel_dir))
    local("chmod -R 777 %s"%out_path)

def create_channeltx(bin_path, yaml_path, out_path, channel_name, crypto_type):
    bin = utils.get_bin_path(bin_path, "configtxgen",crypto_type)
    channel_dir = out_path + channel_name
    if not os.path.exists(channel_dir):
        local("mkdir -p %s"%channel_dir)
    channeltx_name = channel_name + '.tx'
    env = 'FABRIC_CFG_PATH=%s '%yaml_path
    param = ' -profile OrgsChannel -outputCreateChannelTx %s/%s -channelID %s'%(channel_dir, channeltx_name, channel_name)
    
    command = env + bin + param
    local(command)
    local("chmod -R 777 %s"%out_path)


def update_anchor(bin_path, yaml_path, out_path, channel_name, org_id, orderer_address,order_tls_path, domain_name, crypto_type):

    create_anchor_tx(bin_path, yaml_path, out_path, channel_name, org_id,crypto_type)

    channel_dir = out_path + channel_name

    msp_path = yaml_path + "crypto-config/peerOrganizations/%s.%s/users/Admin@%s.%s/msp"%(org_id,domain_name,org_id,domain_name)
    env = ' FABRIC_CFG_PATH=%s '%yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=%s'%org_id
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s '%msp_path
    bin = utils.get_bin_path(bin_path, "peer",crypto_type)
    param = ' channel update -o %s -c %s -f %s/%sanchors.tx'%(orderer_address, channel_name, channel_dir, org_id)
    tls = ' --tls --cafile %s'%order_tls_path

    command = env + bin + param + tls
    local(command)
    local("chmod -R 777 %s"%out_path)


def create_anchor_tx(bin_path, yaml_path, out_path, channel_name, org_id, crypto_type):
    channel_dir = out_path + channel_name
    env = ' FABRIC_CFG_PATH=%s '%yaml_path
    param = ' -profile OrgsChannel -outputAnchorPeersUpdate %s/%sanchors.tx -channelID %s -asOrg %s'%(channel_dir, org_id, channel_name, org_id)

    bin = utils.get_bin_path(bin_path, "configtxgen",crypto_type)
    command = env + bin + param
    local(command)
    local("chmod -R 777 %s"%out_path)

def join_channel(bin_path, yaml_path, out_path, channel_name, peer_address, peer_id, org_id, domain_name, crypto_type):
    channel_block = channel_name + '.block'
    tls_root_file = yaml_path + "crypto-config/peerOrganizations/%s.%s/peers/peer%s.%s.%s/tls/ca.crt"%(org_id,domain_name,peer_id,org_id,domain_name)
    msp_path = yaml_path + "crypto-config/peerOrganizations/%s.%s/users/Admin@%s.%s/msp"%(org_id,domain_name,org_id,domain_name)
    channel_dir = out_path + channel_name
    env = ' FABRIC_CFG_PATH=%s '%yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=%s'%org_id
    env = env + ' CORE_PEER_TLS_ROOTCERT_FILE=%s'%tls_root_file
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s'%msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true'
    env = env + ' CORE_PEER_ADDRESS=%s '%peer_address
    bin = utils.get_bin_path(bin_path, "peer",crypto_type)
    param = ' channel join -b %s/%s'%(channel_dir, channel_block)

    command = env + bin + param
    local(command)

