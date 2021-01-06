import sys
from fabric.api import run, settings, cd
import utils

reload(sys)
sys.setdefaultencoding('utf8')

def remove_node(type, node_name):
    with settings(warn_only=True):
        file_name = "~/fabricNetwork/yaml/%s.yaml" % node_name
        if utils.check_remote_exist(file_name):
            run("docker-compose -f %s down --volumes" % file_name)
            run("docker rm -f %s" %node_name)
        if type == "peer":
            if utils.check_container_exist("dev\-%s"%node_name):
                run("unset GREP_OPTIONS && docker ps -a |grep 'dev\-%s'|awk '{print $1}'|xargs docker rm -f" % node_name)
            if utils.check_image_exist("dev\-%s"%node_name):
                run("unset GREP_OPTIONS && docker images |grep 'dev\-%s'|awk '{print $3}'|xargs docker rmi -f" % node_name)
                # run("docker network prune -f")
                # run("docker volume prune -f")

def remove_data(image, mount_path, domain_name):
    with settings(warn_only=True):
        del_cmd = "rm -rf /ledgerData/*%s && rm -rf /deployData/fabricNetwork" % domain_name
        run("docker run -it --rm -v %s:/ledgerData -v  ~/:/deployData %s sh -c '%s'" % (mount_path, image, del_cmd))


def remove_client():
    with settings(warn_only=True):
        run("docker ps -a | awk '{print $1}' | xargs docker rm -f")
        run("docker network prune -f")
        run("docker volume prune -f")
        utils.kill_process("eventserver")


def remove_jmeter():
    with settings(warn_only=True):
        utils.kill_process("jmeter")
