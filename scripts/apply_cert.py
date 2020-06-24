#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import utils
from fabric.api import local, lcd, put, run, cd

reload(sys)
sys.setdefaultencoding('utf8')


def generate_genesis_block(model, bin_path, cfg_path, out_path, crypto_type):
    if not os.path.exists(out_path + "crypto-config"):
        with lcd(out_path):
            local("tar -zxvf crypto-config.tar.gz")
    if not os.path.exists(cfg_path + "core.yaml"):
        local("cp %s/core.yaml %s" % (bin_path, cfg_path))

    tool = utils.get_bin_path(bin_path, "configtxgen", crypto_type)
    channel_path = out_path + "channel-artifacts"
    local("rm -rf %s" % channel_path)
    local("mkdir -p %s" % channel_path)
    env = "FABRIC_CFG_PATH=%s" % cfg_path
    local("%s %s -profile %s -channelID byfn-sys-channel -outputBlock %s/genesis.block" % (env, tool, model, channel_path))
    with lcd(out_path):
        local("tar -zcvf channel-artifacts.tar.gz channel-artifacts")
        local("chmod -R 777 channel-artifacts")


## Generates orderer Org certs using cryptogen tool
def generate_certs(bin_path, cfg_path, out_path, crypto_type):
    cryptotool = utils.get_bin_path(bin_path, "cryptogen", crypto_type)
    yamlfile = cfg_path + "crypto-config.yaml"
    mm_path = out_path + "crypto-config"

    local("%s generate --config=%s --output='%s'" % (cryptotool, yamlfile, mm_path))
    local("chmod -R 777 %s" % mm_path)


def put_cryptoconfig(config_path, type, node_name, org_name, cert_peer_name):
    run("mkdir -p ~/fabricNetwork/yaml")
    with lcd(config_path):
        if type == "orderer":
            local('tar -zcvf %s_crypto-config.tar.gz crypto-config/ordererOrganizations/%s/orderers/%s' % (
                node_name, org_name, node_name))
            copy_file(config_path, "%s_crypto-config.tar.gz" % node_name)
            copy_file(config_path, "channel-artifacts.tar.gz")
            # copy_file(config_path,"kafkaTLSclient.tar.gz")
        elif type == "kafka":
            copy_file(config_path, "kafkaTLSserver.tar.gz")
        elif type == "peer":
            local('tar -zcvf %s_crypto-config.tar.gz crypto-config/peerOrganizations/%s/peers/%s' % (
                node_name, org_name, node_name))
            copy_file(config_path, "%s_crypto-config.tar.gz" % node_name)
        elif type == "api":
            local('tar -zcvf %s_crypto-config.tar.gz crypto-config' %node_name)
            copy_file(config_path, "%s_crypto-config.tar.gz" % node_name)
        elif type == "explorer":
            peerTlsFile = "crypto-config/peerOrganizations/%s/peers/%s" % (org_name, cert_peer_name)
            AdminDir = "crypto-config/peerOrganizations/%s/users/Admin@%s" % (org_name, org_name)
            local('tar -zcvf %s_crypto-config.tar.gz %s %s ' % (node_name, peerTlsFile, AdminDir))
            copy_file(config_path, "%s_crypto-config.tar.gz" % node_name)
            local('tar -zcvf %s.tar.gz %s' % (node_name, node_name))
            copy_file(config_path, "%s.tar.gz" % node_name)
            with cd("~/fabricNetwork"):
                run("rm -rf block_fabric_explorer")
                run("mv %s block_fabric_explorer"%node_name)
                run("mv block_fabric_explorer/block_fabric_explorer.yaml ~/fabricNetwork/yaml/")


def copy_file(config_path, file_name):
    remote_file = "~/fabricNetwork/%s" % file_name
    if not utils.check_remote_exist(remote_file):
        put("%s%s" % (config_path, file_name), "~/fabricNetwork/")
        with cd("~/fabricNetwork"):
            run("tar zxfm %s" % file_name)
            run("rm -rf %s"%file_name)


def generate_certs_to_ca(bin_path, out_path, crypto_type, node_type, full_name, org_name, ca_url, tlsca_url, admin_name,
                         admin_pw):
    ca_tool = utils.get_bin_path(bin_path, "fabric-ca-client", crypto_type)
    cert_path = out_path + "crypto-config"

    node_password = "password"
    ca_admin = "%s/cadata/%s/%s" % (out_path, org_name, admin_name)
    user = "%s/cadata/%s/%s" % (out_path, org_name, full_name)
    tls_ca_admin = "%s/tlscadata/%s/%s" % (out_path, org_name, admin_name)
    tls_user = "%s/tlscadata/%s/%s" % (out_path, org_name, full_name)
    print "----------------------------------------"
    print "------generate  %s cert start-----------" % full_name
    print "----------------------------------------"
    print "-------------%s admin login-------------" % full_name
    if not os.path.exists(ca_admin):
        print "---------%s  do not exist, need admin enroll---------------" % ca_admin
        local("%s enroll -u http://%s:%s@%s -H %s" % (ca_tool, admin_name, admin_pw, ca_url, ca_admin))
        local("%s enroll -u http://%s:%s@%s -H %s" % (ca_tool, admin_name, admin_pw, tlsca_url, tls_ca_admin))
    else:
        print "---------%s  already exist---------------" % ca_admin
    if node_type == "orderer":
        org_path = "%s/ordererOrganizations/%s" % (cert_path, org_name)
        org_user = "%s/orderers/%s" % (org_path, full_name)
    else:
        org_path = "%s/peerOrganizations/%s" % (cert_path, org_name)
        org_user = "%s/peers/%s" % (org_path, full_name)

    print "----------------------------------------------"
    print "------generate  %s tls cert start----------------" % full_name
    print "----------------------------------------------"
    # 注册登记
    local("%s register --id.name %s_tls --id.type %s --id.secret %s -H %s" % (
        ca_tool, full_name, node_type, node_password, tls_ca_admin))
    local("%s enroll -u http://%s_tls:%s@%s --csr.hosts %s,%s -H %s" % (
        ca_tool, full_name, node_password, tlsca_url, full_name, node_type, tls_user))

    # 生成组织 ca tls证书
    local("mkdir -p %s/msp/tlscacerts/" % org_path)
    local("cp %s/msp/cacerts/*.pem %s/msp/tlscacerts/tlsca.%s-cert.pem" % (tls_ca_admin, org_path, org_name))
    # 生成 amdmin tls 证书
    local("mkdir -p %s/users/%s/tls" % (org_path, admin_name))
    local("cp %s/msp/cacerts/*.pem %s/users/%s/tls/ca.crt" % (tls_ca_admin, org_path, admin_name))
    local("cp %s/msp/signcerts/*.pem %s/users/%s/tls/client.crt" % (tls_ca_admin, org_path, admin_name))
    local("cp %s/msp/keystore/*_sk %s/users/%s/tls/client.key" % (tls_ca_admin, org_path, admin_name))
    # 生成 节点 tls 证书
    local("mkdir -p %s/tls" % org_user)
    local("cp %s/msp/cacerts/*.pem %s/tls/ca.crt" % (tls_user, org_user))
    local("cp %s/msp/signcerts/*.pem %s/tls/server.crt" % (tls_user, org_user))
    local("cp %s/msp/keystore/*_sk %s/tls/server.key" % (tls_user, org_user))

    print "---------------------------------------------------------"
    print "-----------------generate %s cert----------------------" % full_name
    print "---------------------------------------------------------"
    # 注册登记
    local("%s register --id.name %s --id.type %s --id.secret %s -H %s" % (
        ca_tool, full_name, node_type, node_password, ca_admin))
    local("%s enroll -u http://%s:%s@%s --csr.hosts %s,%s -H %s" % (
        ca_tool, full_name, node_password, ca_url, full_name, node_type, user))

    # 生成组织ca证书
    local("mkdir -p %s/msp/cacerts/" % org_path)
    local("cp %s/msp/cacerts/*.pem %s/msp/cacerts/ca.%s-cert.pem" % (ca_admin, org_path, org_name))

    # 生成组织admin证书
    local("mkdir -p %s/msp/admincerts/" % org_path)
    local("cp %s/msp/signcerts/cert.pem %s/msp/admincerts/%s-cert.pem" % (ca_admin, org_path, admin_name))
    local("mkdir -p %s/users/%s" % (org_path, admin_name))
    local("cp -r %s/msp %s/users/%s/" % (org_path, org_path, admin_name))
    local("cp -r %s/msp/keystore %s/users/%s/msp/" % (ca_admin, org_path, admin_name))
    local("mkdir -p %s/users/%s/msp/signcerts" % (org_path, admin_name))
    local(
        "cp -r %s/users/%s/msp/admincerts/*  %s/users/%s/msp/signcerts/" % (org_path, admin_name, org_path, admin_name))
    # 生成组织节点证书
    local("mkdir -p %s/msp/signcerts" % org_user)
    local("cp -r %s/msp/ %s/" % (org_path, org_user))
    local("cp %s/msp/signcerts/cert.pem %s/msp/signcerts/%s-cert.pem" % (user, org_user, full_name))
    local("cp -r %s/msp/keystore/ %s/msp/" % (user, org_user))

    config_tpl_path = utils.get_bin_path(bin_path, "config.yamlt", "")

    # 生成config.yaml配置文件
    if node_type == "orderer":
        print "orderer do not generate config.yaml"
    else:
        local('sed "s/ORG_NAME/%s/g" %s > %s/msp/config.yaml' % (org_name, config_tpl_path, org_path))
        local('sed "s/ORG_NAME/%s/g" %s > %s/msp/config.yaml' % (org_name, config_tpl_path, org_user))

    local("chmod -R 777 %s" % org_path)
