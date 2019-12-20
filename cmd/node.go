package cmd

import (
	"fmt"
	"sync"
)

func StartNode(stringType string) error {
	if err := WriteHost(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	StartN := func(Ip, Sshuser, Sshpwd, Sshport, Sshkey, NodeName string, w1 *sync.WaitGroup) {
		defer w1.Done()
		obj := NewFabCmd("add_node.py", Ip, Sshuser, Sshpwd, Sshport, Sshkey)
		err := obj.RunShow("start_node", NodeName, ConfigDir())
		if err != nil {
			fmt.Println("start node err or")
		}
	}
	if stringType == "all" || stringType == TypeKafka {
		for _, kafka := range GlobalConfig.Kafkas {
			wg.Add(1)
			go StartN(kafka.Ip, kafka.SshUserName, kafka.SshPwd, kafka.SshPort, kafka.SshKey, kafka.NodeName, &wg)
		}
	}
	if stringType == "all" || stringType == TypeZookeeper {
		for _, zk := range GlobalConfig.Zookeepers {
			wg.Add(1)
			go StartN(zk.Ip, zk.SshUserName, zk.SshPwd, zk.SshPort, zk.SshKey, zk.NodeName, &wg)
		}
	}
	if stringType == "all" || stringType == TypeOrder {
		for _, ord := range GlobalConfig.Orderers {
			wg.Add(1)
			go StartN(ord.Ip, ord.SshUserName, ord.SshPwd, ord.SshPort, ord.SshKey, ord.NodeName, &wg)
		}
	}
	if stringType == "all" || stringType == TypePeer {
		for _, peer := range GlobalConfig.Peers {
			wg.Add(1)
			go StartN(peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey, peer.NodeName, &wg)
		}
	}
	if GlobalConfig.CaType == "fabric-ca" {
		if stringType == "all" || stringType == TypeCa {
			for _, ca := range GlobalConfig.Cas {
				wg.Add(1)
				go StartN(ca.Ip, ca.SshUserName, ca.SshPwd, ca.SshPort, ca.SshKey, ca.NodeName, &wg)
			}
		}
	}
	wg.Wait()
	return nil
}

func WriteHost() error {
	for _, ord := range GlobalConfig.Orderers {
		if err := LocalHostsSet(ord.Ip, ord.NodeName); err != nil {
			return err
		}
	}
	for _, peer := range GlobalConfig.Peers {
		if err := LocalHostsSet(peer.Ip, peer.NodeName); err != nil {
			return err
		}
	}
	//for _, kafka := range GlobalConfig.Kafkas {
	//	if err := LocalHostsSet(kafka.Ip, kafka.NodeName); err != nil {
	//		return err
	//	}
	//}
	//for _, zk := range GlobalConfig.Zookeepers {
	//	if err := LocalHostsSet(zk.Ip, zk.NodeName); err != nil {
	//		return err
	//	}
	//}

	return nil
}

func DeleteObj(stringType string) error {
	var wg sync.WaitGroup
	StopN := func(Ip, Sshuser, Sshpwd, Sshport, Sshkey, Ty, ImageName, MountName string, w1 *sync.WaitGroup) {
		defer w1.Done()
		obj := NewFabCmd("removenode.py", Ip, Sshuser, Sshpwd, Sshport, Sshkey)
		err := obj.RunShow("remove_node", Ty, ImageName, GlobalConfig.MountPath, MountName)
		if err != nil {
			fmt.Println("stopnode err or")
		}
	}
	if stringType == "all" || stringType == TypeKafka {
		for _, kafka := range GlobalConfig.Kafkas {
			wg.Add(1)
			imageName := fmt.Sprintf("%s/fabric-kafka:%s", GlobalConfig.ImagePre, GlobalConfig.ImageTag)
			mountName := fmt.Sprintf("kafka%s", kafka.Id)
			go StopN(kafka.Ip, kafka.SshUserName, kafka.SshPwd, kafka.SshPort, kafka.SshKey, TypeKafka, imageName, mountName, &wg)
		}
	}
	if stringType == "all" || stringType == TypeZookeeper {
		for _, zk := range GlobalConfig.Zookeepers {
			wg.Add(1)
			imageName := fmt.Sprintf("%s/fabric-zookeeper:%s", GlobalConfig.ImagePre, GlobalConfig.ImageTag)
			mountName := fmt.Sprintf("zk%s", zk.Id)
			go StopN(zk.Ip, zk.SshUserName, zk.SshPwd, zk.SshPort, zk.SshKey, TypeZookeeper, imageName, mountName, &wg)
		}
	}
	if stringType == "all" || stringType == TypeOrder {
		for _, ord := range GlobalConfig.Orderers {
			wg.Add(1)
			imageName := fmt.Sprintf("%s/fabric-orderer:%s", GlobalConfig.ImagePre, GlobalConfig.ImageTag)
			mountName := fmt.Sprintf("orderer%s.ord%s.%s", ord.Id, ord.OrgId, GlobalConfig.Domain)
			go StopN(ord.Ip, ord.SshUserName, ord.SshPwd, ord.SshPort, ord.SshKey, TypeOrder, imageName, mountName, &wg)
		}
	}
	if stringType == "all" || stringType == TypePeer {
		for _, peer := range GlobalConfig.Peers {
			wg.Add(1)
			imageName := fmt.Sprintf("%s/fabric-peer:%s", GlobalConfig.ImagePre, GlobalConfig.ImageTag)
			mountName := fmt.Sprintf("peer%s.org%s.%s", peer.Id, peer.OrgId, GlobalConfig.Domain)
			go StopN(peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey, TypePeer, imageName, mountName, &wg)
		}
	}
	if stringType == "all" || stringType == TypeCa {
		for _, ca := range GlobalConfig.Cas {
			wg.Add(1)
			imageName := fmt.Sprintf("%s/fabric-ca:%s", GlobalConfig.ImagePre, GlobalConfig.ImageTag)
			go StopN(ca.Ip, ca.SshUserName, ca.SshPwd, ca.SshPort, ca.SshKey, TypeCa, imageName, ca.NodeName, &wg)
		}
	}
	wg.Wait()
	return nil
}

func LocalHostsSet(ip, domain string) error {
	if ip == domain {
		return nil
	}
	if err := ModifyHosts("/etc/hosts", ip, domain); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func CheckNode(stringType string) error {
	if err := WriteHost(); err != nil {
		return err
	}
	if stringType == "all" || stringType == TypeKafka {
		for _, kafka := range GlobalConfig.Kafkas {
			obj := NewFabCmd("add_node.py", kafka.Ip, kafka.SshUserName, kafka.SshPwd, kafka.SshPort, kafka.SshKey)
			if err := obj.RunShow("check_node"); err != nil {
				return err
			}
		}
	}
	if stringType == "all" || stringType == TypeZookeeper {
		for _, zk := range GlobalConfig.Zookeepers {
			obj := NewFabCmd("add_node.py", zk.Ip, zk.SshUserName, zk.SshPwd, zk.SshPort, zk.SshKey)
			if err := obj.RunShow("check_node"); err != nil {
				return err
			}
		}
	}
	if stringType == "all" || stringType == TypeOrder {
		for _, ord := range GlobalConfig.Orderers {
			obj := NewFabCmd("add_node.py", ord.Ip, ord.SshUserName, ord.SshPwd, ord.SshPort, ord.SshKey)
			if err := obj.RunShow("check_node"); err != nil {
				return err
			}
		}
	}
	if stringType == "all" || stringType == TypePeer {
		for _, peer := range GlobalConfig.Peers {
			obj := NewFabCmd("add_node.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			if err := obj.RunShow("check_node"); err != nil {
				return err
			}
		}
	}
	if stringType == "all" || stringType == TypeCa {
		for _, ca := range GlobalConfig.Cas {
			obj := NewFabCmd("add_node.py", ca.Ip, ca.SshUserName, ca.SshPwd, ca.SshPort, ca.SshKey)
			if err := obj.RunShow("check_node"); err != nil {
				return err
			}
		}
	}
	return nil
}
