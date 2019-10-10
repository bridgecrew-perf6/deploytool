#!/bin/python

import sys
import os
from fabric.api import cd, put, lcd, local, run, get
import utils
reload(sys)
sys.setdefaultencoding('utf8')

#cp jmeter config to remote
def start_jmeter(config_dir):
    dir_name = "jmeter_config"
    with lcd(config_dir):
        local("tar -zcvf jmeterjmx.tar.gz jmeter.jmx")
        #remote yaml
        run("rm -rf ~/deployFabricTool/%s"%dir_name)
        run("mkdir -p ~/deployFabricTool/%s"%dir_name)
        put("jmeterjmx.tar.gz", "~/deployFabricTool/%s"%dir_name)
        local("rm jmeterjmx.tar.gz")
    with cd("~/deployFabricTool/%s"%dir_name):
        run("tar zxvfm jmeterjmx.tar.gz")
        run("rm jmeterjmx.tar.gz")
        run("screen -d -m ~/jmeter/apache-jmeter-3.2/bin/jmeter -n -t jmeter.jmx -l jmeter.jtl", pty=False)

#get jmeter log from remote
def get_jmeter_log(config_dir, log_dir):
    dir = "%sevent_logs/%s"%(config_dir,log_dir)
    local("mkdir -p %s"%dir)
    file = '%s/jmeter.jtl'%dir
    log =  '%s/jmeter_send.txt'%dir
    if os.path.exists(file):
        local("rm -rf %s"%file)
    get('~/deployFabricTool/jmeter_config/jmeter.jtl',file)
    get('~/deployFabricTool/jmeter_config/jmeter.log',log)
    with lcd(dir):
        local("~/jmeter/apache-jmeter-3.2/bin/jmeter -g %s -e -o ./jmeterReport"%file)
        local("rm -rf jmeter.log")

#get eventserver log from remote
def get_eventserver_log(yaml_name, config_dir, log_dir):
    dir = "%sevent_logs/%s"%(config_dir,log_dir)
    local("mkdir -p %s"%dir)
    file = '%s/%s_eventserver.log'%(dir,yaml_name)
    if os.path.exists(file):
        local("rm -rf %s"%file)
    get('~/deployFabricTool/event_server/%s/eventserver.log'%yaml_name,'%s'%file)
    #echo  empty log
    run("cat /dev/null > ~/deployFabricTool/event_server/%s/eventserver.log"%yaml_name)

# remote
def start_haproxy(config_dir):
    dir_name = "haproxy_config"
    with lcd(config_dir):
        local("tar -zcvf haproxyconfig.tar.gz haproxy_config")
        #remote yaml
        run("rm -rf ~/deployFabricTool/%s"%dir_name)
        run("mkdir -p ~/deployFabricTool/")
        put("haproxyconfig.tar.gz", "~/deployFabricTool/")
        local("rm haproxyconfig.tar.gz")
    with cd("~/deployFabricTool"):
        run("tar zxvfm haproxyconfig.tar.gz")
        run("rm -rf haproxyconfig.tar.gz")
    with cd("~/deployFabricTool/haproxy_config"):
        run("docker-compose -f docker-compose.yaml down")
        run("docker-compose -f docker-compose.yaml up -d")


def start_nmon(rate,times_number,out_file_name):
    run("rm -rf ~/nmon_log")
    run("mkdir -p ~/nmon_log")
    with cd("~/nmon_log"):
        utils.kill_process("nmon")
        run("nmon -s%s -c%s -F %s.nmon"%(rate,times_number,out_file_name))

#get nmon log from remote
def get_nmon_log(rate,times_number,out_file_name,config_dir,log_dir):
    dir = "%sevent_logs/%s"%(config_dir,log_dir)
    local("mkdir -p %s"%dir)
    file = '%s/%s_nmon'%(dir,out_file_name)
    if os.path.exists(file):
        local("rm -rf %s"%file)
    with cd("~/nmon_log"):
        utils.kill_process("nmon")
        get("~/nmon_log/%s.nmon"%out_file_name,"%s"%file)
