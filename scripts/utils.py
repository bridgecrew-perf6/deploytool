import sys
from fabric.api import run, settings, local

reload(sys)
sys.setdefaultencoding('utf8')


def get_bin_path(path, type, crypto_type):
    if crypto_type == "GM":
        return ' BCCSP_CRYPTO_TYPE=GM %s/%s ' % (path, type)
    return ' %s/%s ' % (path, type)


def kill_process(name):
    # kill the jmeter processes for unified order project
    with settings(warn_only=True):
        result = run("pgrep %s" % name)
        if result != "":
            run("kill -9 %s" % result)


def check_local_exist(file_path):
    with settings(warn_only=True):
        result = local("[ -e '" + file_path + "' ]")
    if result.return_code == 0:
        return True
    else:
        return False

def check_remote_exist(file_path):
    if int(run(" [ -e "+file_path+" ] && echo 11 || echo 10")) == 11:
        return True
    else:
        return False

def check_container_exist(name):
    containers = run('unset GREP_OPTIONS && docker ps -a |grep "%s" | wc -l' % name)
    if containers != "0":
        return True
    else:
        return False

def check_image_exist(name):
    containers = run('unset GREP_OPTIONS && docker images |grep "%s" | wc -l' % name)
    if containers != "0":
        return True
    else:
        return False

def check_network_exist(name):
    networks = run('unset GREP_OPTIONS && docker network ls |grep "%s" | wc -l' % name)
    if networks != "0":
        return True
    else:
        return False


def set_domain_name(network_name, node_full_name, domain_ip, domain_name):
    set_cmd = "echo %s %s >> /etc/hosts" % (domain_ip, domain_name)
    yaml_file = "~/networklist/%s/%s/%s.yaml" % (network_name, node_full_name, node_full_name)
    run("docker-compose -f %s exec %s bash -c '%s'" % (yaml_file, node_full_name, set_cmd))


def get_domain_name(network_name, node_full_name, domain_name):
    get_cmd = "unset GREP_OPTIONS && cat /etc/hosts | grep -E %s | awk '{print \\\$1}'" % domain_name
    yaml_file = "~/networklist/%s/%s/%s.yaml" % (network_name, node_full_name, node_full_name)
    out = run('docker-compose -f %s exec %s bash -c "%s"' % (yaml_file, node_full_name, get_cmd))


def chmod_all(path, mode):
    local("chmod -R %s %s" % (mode, path))


def rm_local(path):
    local("rm -rf %s" % path)
