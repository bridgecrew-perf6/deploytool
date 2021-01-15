#!/usr/bin/env python
# -*- coding: utf-8 -*-

import json
import os
import sys
from fabric.api import local, task

import utils

reload(sys)
sys.setdefaultencoding('utf8')

fabfile = os.environ.get("FABFILE")


@task
def generate_org_json(bin_path, update_path, yaml_path, org_id, domain_name):
    org_path = yaml_path + "crypto-config/peerOrganizations/%s.%s" % (org_id, domain_name)
    org_json_path = org_path + "/org.json"
    env = "FABRIC_CFG_PATH=%s" % yaml_path
    bin = utils.get_bin_path(bin_path, "configtxgen", "")
    param = " -configPath " + update_path + " -printOrg " + org_id + " > " + org_json_path

    command = env + bin + param

    local(command)


@task
def add_org_new(bin_path, yaml_path, org_id, add_org_list, orderer_address, orderer_tls_path, domain_name,
                channel_name):
    msp_path = yaml_path + "crypto-config/peerOrganizations/%s.%s/users/Admin@%s.%s/msp" % (
        org_id, domain_name, org_id, domain_name)
    update_config_dir = yaml_path + "updateconfig/"

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

    org_array = json.loads(add_org_list)
    if len(org_array) <= 0:
        print  "python scripts add_org error org_list input error"
        os._exit(1)

    # Modify the configuration to append the new org
    last_msp_id = ''
    for value in org_array:
        new_msp_id = value["mspid"]
        new_org_name = value["orgname"]
        anchor_peer_name = value["anchorpeer"]
        port = value["port"]
        new_org_path = yaml_path + "crypto-config/peerOrganizations/%s.%s" % (new_org_name, domain_name)

        # generate new org json used to update config block
        generate_org_json(bin_path, update_config_dir, yaml_path, new_msp_id, domain_name)
        if last_msp_id == '':
            old_config_json = "config"
            new_config_json = new_msp_id
        else:
            old_config_json = last_msp_id
            new_config_json = new_msp_id
        command = 'jq -s \'.[0] * {"channel_group":{"groups":{"Application":{"groups": {"%s":.[1]}}}}}\' %s/%s.json %s/org.json' % (
            new_msp_id, update_config_dir, old_config_json, new_org_path)
        command += ' | jq \'.channel_group.groups.Application.groups.%s.values.AnchorPeers.mod_policy= "Admins"\'' % (
            new_msp_id)
        command += ' | jq \'.channel_group.groups.Application.groups.%s.values.AnchorPeers.value.anchor_peers[0].host= "%s"\'' % (
        new_msp_id, anchor_peer_name)
        command += ' | jq \'.channel_group.groups.Application.groups.%s.values.AnchorPeers.value.anchor_peers[0].port= "%s"\' > %s/%s.json' % (
        new_msp_id, port, update_config_dir, new_config_json)
        local(command)
        last_msp_id = new_msp_id

    command = 'cp %s/%s.json %s/update_config.json' % (update_config_dir, last_msp_id, update_config_dir)
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

    # Set the peerOrg admin of an org and signing the config update
    bin = utils.get_bin_path(bin_path, "peer", "")
    param = ' channel signconfigtx -f %s/config_update_in_envelope.pb' % update_config_dir
    command = env + bin + param
    local(command)

    param = ' channel update -f %s/config_update_in_envelope.pb -c %s -o %s' % (
        update_config_dir, channel_name, orderer_address)
    command = env + bin + param + tls
    local(command)
