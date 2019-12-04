#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import utils
from fabric.api import local

reload(sys)
sys.setdefaultencoding('utf8')


def pkg_chaincode(bin_path, config_path, org_id, domain_name, ccname, ccversion, ccpath, ccinstalltype, crypto_type):
    msp_path = config_path + "crypto-config/peerOrganizations/org%s.%s/users/Admin@org%s.%s/msp" % (
    org_id, domain_name, org_id, domain_name)
    env = ' FABRIC_CFG_PATH=%s ' % config_path
    env = env + ' CORE_PEER_LOCALMSPID=Org%sMSP' % org_id
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s' % msp_path
    bin = utils.get_bin_path(bin_path, "peer", crypto_type)
    if ccinstalltype == "path":
        param = ' chaincode package -n %s -p %s -v %s %s/%s_%s.pkg' % (
        ccname, ccpath, ccversion, config_path, ccname, ccversion)
        command = env + bin + param
        local(command)


def install_chaincode(bin_path, config_path, peer_address, peer_id, org_id, domain_name, ccname, ccversion, ccpath,
                      ccinstalltype, crypto_type):
    tls_root_file = config_path + "crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/ca.crt" % (
    org_id, domain_name, peer_id, org_id, domain_name)
    msp_path = config_path + "crypto-config/peerOrganizations/org%s.%s/users/Admin@org%s.%s/msp" % (
    org_id, domain_name, org_id, domain_name)
    env = ' FABRIC_CFG_PATH=%s ' % config_path
    env = env + ' CORE_PEER_LOCALMSPID=Org%sMSP' % org_id
    env = env + ' CORE_PEER_TLS_ROOTCERT_FILE=%s' % tls_root_file
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true'
    env = env + ' CORE_PEER_ADDRESS=%s ' % peer_address
    bin = utils.get_bin_path(bin_path, "peer", crypto_type)
    if ccinstalltype == "path":
        param = ' chaincode install %s/%s_%s.pkg' % (config_path, ccname, ccversion)
    else:
        param = ' chaincode install  %s' % ccpath

    command = env + bin + param
    local(command)


def instantiate_chaincode(bin_path, operation, yaml_path, peer_address, order_address, peer_id, org_id, domain_name,
                          channel_name, ccname, ccversion, init_param, policy, crypto_type):
    tls_root_file = yaml_path + "crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/ca.crt" % (
    org_id, domain_name, peer_id, org_id, domain_name)
    msp_path = yaml_path + "crypto-config/peerOrganizations/org%s.%s/users/Admin@org%s.%s/msp" % (
    org_id, domain_name, org_id, domain_name)
    order_tls_path = yaml_path + "crypto-config/ordererOrganizations/ord1.%s/orderers/orderer0.ord1.%s/msp/tlscacerts/tlsca.ord1.%s-cert.pem" % (
    domain_name, domain_name, domain_name)
    env = ' FABRIC_CFG_PATH=%s ' % yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=Org%sMSP' % org_id
    env = env + ' CORE_PEER_TLS_ROOTCERT_FILE=%s' % tls_root_file
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true'
    env = env + ' CORE_PEER_ADDRESS=%s ' % peer_address
    bin = utils.get_bin_path(bin_path, "peer", crypto_type)
    param = ' chaincode %s -o %s -C %s -n %s -v %s -c %s -P %s ' % (
    operation, order_address, channel_name, ccname, ccversion, init_param, policy)
    tls = ' --tls --cafile %s' % order_tls_path

    command = env + bin + param + tls
    local(command)


def test_query_tx(bin_path, yaml_path, peer_address, peer_id, org_id, domain_name, channel_name, ccname, tx_args,
                  crypto_type):
    tls_root_file = yaml_path + "crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/ca.crt" % (
    org_id, domain_name, peer_id, org_id, domain_name)
    msp_path = yaml_path + "crypto-config/peerOrganizations/org%s.%s/users/Admin@org%s.%s/msp" % (
    org_id, domain_name, org_id, domain_name)
    env = ' FABRIC_CFG_PATH=%s ' % yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=Org%sMSP' % org_id
    env = env + ' CORE_PEER_TLS_ROOTCERT_FILE=%s' % tls_root_file
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true'
    env = env + ' CORE_PEER_ADDRESS=%s ' % peer_address
    bin = utils.get_bin_path(bin_path, "peer", crypto_type)
    param = '  chaincode query -C %s -n %s -c %s ' % (channel_name, ccname, tx_args)
    command = env + bin + param
    local(command)


def test_chaincode(func, bin_path, yaml_path, peer_address, order_address, peer_id, org_id, domain_name, channel_name,
                   ccname, args, crypto_type):
    tls_root_file = yaml_path + "crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/ca.crt" % (
    org_id, domain_name, peer_id, org_id, domain_name)
    msp_path = yaml_path + "crypto-config/peerOrganizations/org%s.%s/users/Admin@org%s.%s/msp" % (
    org_id, domain_name, org_id, domain_name)
    order_tls_path = yaml_path + "crypto-config/ordererOrganizations/ord1.%s/orderers/orderer0.ord1.%s/msp/tlscacerts/tlsca.ord1.%s-cert.pem" % (
    domain_name, domain_name, domain_name)
    env = ' FABRIC_CFG_PATH=%s ' % yaml_path
    env = env + ' CORE_PEER_LOCALMSPID=Org%sMSP' % org_id
    env = env + ' CORE_PEER_TLS_ROOTCERT_FILE=%s' % tls_root_file
    env = env + ' CORE_PEER_MSPCONFIGPATH=%s' % msp_path
    env = env + ' CORE_PEER_TLS_ENABLED=true'
    env = env + ' CORE_PEER_ADDRESS=%s ' % peer_address
    bin = utils.get_bin_path(bin_path, "peer", crypto_type)
    param = ' chaincode %s -o %s -C %s -n %s -c %s ' % (func, order_address, channel_name, ccname, args)
    tls = ' --tls --cafile %s' % order_tls_path

    command = env + bin + param + tls
    local(command)
