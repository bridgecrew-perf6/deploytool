#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import base64
from fabric.api import local, task

import utils

reload(sys)
sys.setdefaultencoding('utf8')

fabfile = os.environ.get("FABFILE")


@task
def handle_orderer(bin_path, yaml_path, new_node_name, new_node_port, org_id, orderer_address, orderer_tls_path,
                   domain_name,
                   channel_name, isAdd):
    msp_path = yaml_path + "crypto-config/ordererOrganizations/%s.%s/users/Admin@%s.%s/msp" % (
        org_id, domain_name, org_id, domain_name)
    update_config_dir = yaml_path + "updateconfig/"
    local('mkdir -p %s' % update_config_dir)

    env = 'FABRIC_CFG_PATH=%s ' % yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=%s ' % org_id
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s ' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true '

    # fetch block
    bin = utils.get_bin_path(bin_path, "peer", "")
    param = ' channel fetch config %s/config_block.pb -o %s -c %s ' % (update_config_dir, orderer_address, channel_name)
    tls = ' --tls --cafile %s' % orderer_tls_path
    command = env + bin + param + tls
    local(command)

    # decode block
    bin = utils.get_bin_path(bin_path, "configtxlator", "")
    param = ' proto_decode --input %s/config_block.pb --type common.Block | jq .data.data[0].payload.data.config > %s/config.json' % (
        update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    if isAdd == "add":
        OPERATION = "+="
    else:
        OPERATION = "-="

    new_orderer_tls_file = yaml_path + "crypto-config/ordererOrganizations/%s.%s/orderers/%s/tls/server.crt" % (
        org_id, domain_name, new_node_name)
    # if baseimage=ubuntu
    # command = 'base64 %s -w 0' % new_orderer_tls_file
    #new_orderer_tls_str = utils.safe_local(command)
    # elseif baseimage=alpine
    cert_str = ""
    with open(new_orderer_tls_file, 'rb') as file:   # 将文件路径和文件名改成自己需要的
        for line in file.readlines():  #  去除每一行之后的换行符
            cert_str += line.strip()
    new_orderer_tls_str = base64.b64encode(cert_str)
   #end
    new_cert_struct = '{"client_tls_cert":"%s","host":"%s","port":%s,"server_tls_cert":"%s"}' % (
        new_orderer_tls_str, new_node_name, new_node_port, new_orderer_tls_str)

    command = ' jq \'.channel_group.groups.Orderer.values.ConsensusType.value.metadata.consenters %s [%s]\' %s/config.json > %s/temp.json' % (
        OPERATION, new_cert_struct, update_config_dir, update_config_dir)
    local(command)
    command = ' jq \'.channel_group.values.OrdererAddresses.value.addresses %s ["%s:%s"]\' %s/temp.json > %s/update_config.json' % (
        OPERATION, new_node_name, new_node_port, update_config_dir, update_config_dir)
    local(command)

    # Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to org3_update_in_envelope.pb
    bin = utils.get_bin_path(bin_path, "configtxlator", "")
    param = ' proto_encode --input %s/%s.json --type common.Config > %s/original_config.pb' % (
        update_config_dir, 'config', update_config_dir)
    command = bin + param
    local(command)

    param = ' proto_encode --input %s/update_config.json --type common.Config > %s/modified_config.pb' % (
        update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    param = ' compute_update --channel_id %s --original %s/original_config.pb --updated %s/modified_config.pb > %s/config_update.pb' % (
        channel_name, update_config_dir, update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    param = ' proto_decode --input %s/config_update.pb  --type common.ConfigUpdate > %s/config_update.json' % (
        update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    command = 'echo \'{"payload":{"header":{"channel_header":{"channel_id":"%s", "type":2}},"data":{"config_update":\'$(cat %s/config_update.json)\'}}}\' | jq . > %s/config_update_in_envelope.json' % (
        channel_name, update_config_dir, update_config_dir)
    local(command)

    param = ' proto_encode --input %s/config_update_in_envelope.json --type common.Envelope > %s/config_update_in_envelope.pb' % (
        update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    # Set the orderer admin of an org and signing the config update
    bin = utils.get_bin_path(bin_path, "peer", "")
    param = ' channel signconfigtx -f %s/config_update_in_envelope.pb' % update_config_dir
    command = env + bin + param
    local(command)

    param = ' channel update -f %s/config_update_in_envelope.pb -c %s -o %s' % (
        update_config_dir, channel_name, orderer_address)
    command = env + bin + param + tls
    local(command)


@task
def update_genesis_block(bin_path, yaml_path, org_id, orderer_address, orderer_tls_path, domain_name):
    msp_path = yaml_path + "crypto-config/ordererOrganizations/%s.%s/users/Admin@%s.%s/msp" % (
        org_id, domain_name, org_id, domain_name)
    update_config_dir = yaml_path + "channel-artifacts/"

    env = 'FABRIC_CFG_PATH=%s ' % yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=%s ' % org_id
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s ' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true '
    # fetch block
    bin = utils.get_bin_path(bin_path, "peer", "")
    param = ' channel fetch config %s/genesis.block -o %s -c byfn-sys-channel ' % (update_config_dir, orderer_address)
    tls = ' --tls --cafile %s' % orderer_tls_path
    command = env + bin + param + tls
    local(command)
