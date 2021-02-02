#!/usr/bin/env python
# -*- coding: utf-8 -*-

from fabric.api import local, task
import os
import utils
import sys
import json
reload(sys)
sys.setdefaultencoding('utf8')

fabfile = os.environ.get( "FABFILE" )

@task
def delete_org_new(bin_path, yaml_path, org_id, rm_org_list, orderer_address, orderer_tls_path, domain_name,
                   channel_name):
    msp_path = yaml_path + "crypto-config/peerOrganizations/%s.%s/users/Admin@%s.%s/msp" % (
        org_id, domain_name, org_id, domain_name)
    update_config_dir = yaml_path + "updateconfig/"
    local('mkdir -p %s'%update_config_dir)
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

    # remove Org3
    command = 'jq . %s/config.json '%update_config_dir
    msp_array = json.loads(rm_org_list)
    
    for value in msp_array:
        peer_org_mspid = value['orgid']
        command += ' | jq \'del(.channel_group.groups.Application.groups.%s)\''%peer_org_mspid

    command += ' > %s/updated_config.json'%update_config_dir
    local(command)

    # Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to org3_update_in_envelope.pb
    bin = utils.get_bin_path(bin_path, "configtxlator", "")
    param = ' proto_encode --input %s/%s.json --type common.Config > %s/original_config.pb'%(update_config_dir, 'config', update_config_dir)
    command = bin + param
    local(command)

    param = ' proto_encode --input %s/%s.json --type common.Config > %s/modified_config.pb'%(update_config_dir, 'updated_config', update_config_dir)
    command = bin + param
    local(command)

    param = ' compute_update --channel_id %s --original %s/original_config.pb --updated %s/modified_config.pb > %s/config_update.pb'%(channel_name, update_config_dir, update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    param = ' proto_decode --input %s/config_update.pb  --type common.ConfigUpdate > %s/config_update.json'%(update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    command = 'echo \'{"payload":{"header":{"channel_header":{"channel_id":"%s", "type":2}},"data":{"config_update":\'$(cat %s/config_update.json)\'}}}\' | jq . > %s/config_update_in_envelope.json'%(channel_name, update_config_dir, update_config_dir)
    local(command)

    param = ' proto_encode --input %s/config_update_in_envelope.json --type common.Envelope > %s/config_update_in_envelope.pb'%(update_config_dir, update_config_dir)
    command = bin + param
    local(command)

    # Set the peerOrg admin of an org and signing the config update
    bin = utils.get_bin_path(bin_path, "peer", "")
    param = ' channel signconfigtx -f %s/config_update_in_envelope.pb'%(update_config_dir)
    command = env + bin + param
    local(command)

    param = ' channel update -f %s/config_update_in_envelope.pb -c %s -o %s'%(update_config_dir, channel_name, orderer_address)
    command = env + bin + param + tls
    local(command)
