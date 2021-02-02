package cmd

import (
	"fmt"
	"sync"
)

func RemoveData(stringType string) error {
	var wg sync.WaitGroup
	for _, obj := range HostMapList {
		wg.Add(1)
		go func(n NodeObj, w *sync.WaitGroup) {
			defer w.Done()
			cmd := NewFabCmd("removenode.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
			err := cmd.RunShow("remove_data", n.ImageName, GlobalConfig.MountPath, GlobalConfig.Domain)
			if err != nil {
				fmt.Println("remove_node err or")
			}
		}(obj, &wg)
	}
	wg.Wait()
	return nil
}

func ForeachNode(nodeList []NodeObj, isStart bool) {
	var wg sync.WaitGroup
	for _, node := range nodeList {
		wg.Add(1)
		go func(n NodeObj, w *sync.WaitGroup) {
			defer w.Done()
			var err error
			if isStart {
				obj := NewFabCmd("add_node.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("start_node", n.NodeName, ConfigDir())
			} else {
				obj := NewFabCmd("removenode.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("remove_node", n.NodeType, n.NodeName)
			}
			if err != nil {
				fmt.Printf("ForeachNode error: %s\n", err.Error())
			}
		}(node, &wg)
	}
	wg.Wait()
}

func operationExplorer(isStart bool) {
	var wg sync.WaitGroup
	for _, e := range GlobalConfig.Explorers {
		wg.Add(1)
		go func(n ExplorerObj, w *sync.WaitGroup) {
			defer w.Done()
			var err error
			if isStart {
				obj := NewFabCmd("add_node.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("start_node", "block_fabric_explorer", ConfigDir())
			} else {
				obj := NewFabCmd("removenode.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("remove_node", n.NodeType, "block_fabric_explorer")
			}
			if err != nil {
				fmt.Printf("handleExplorer error: %s\n", err.Error())
			}
		}(e, &wg)
	}
	wg.Wait()
}
func operationApiserver(isStart bool) {
	var wg sync.WaitGroup
	for _, e := range GlobalConfig.Apiservers {
		wg.Add(1)
		go func(n ApiserverObj, w *sync.WaitGroup) {
			defer w.Done()
			var err error
			if isStart {
				obj := NewFabCmd("add_node.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("start_node", "apiserver", ConfigDir())
			} else {
				obj := NewFabCmd("removenode.py", n.Ip, n.SshUserName, n.SshPwd, n.SshPort, n.SshKey)
				err = obj.RunShow("remove_node", n.NodeType, "apiserver")
			}
			if err != nil {
				fmt.Printf("handleAPiserver error: %s\n", err.Error())
			}
		}(e, &wg)
	}
	wg.Wait()
}

func RunRmNode(nodename string) error {
	if nodename == "" {
		return fmt.Errorf("nodename is empty")
	}
	for _, peer := range GlobalConfig.Peers {
		if peer.NodeName == nodename {
			//删除节点
			obj := NewFabCmd("removenode.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			if err := obj.RunShow("remove_node", peer.NodeType, peer.NodeName); err != nil {
				return err
			}
		}
	}
	for _, orderer := range GlobalConfig.Orderers {
		if orderer.NodeName == nodename {
			//删除节点
			obj := NewFabCmd("removenode.py", orderer.Ip, orderer.SshUserName, orderer.SshPwd, orderer.SshPort, orderer.SshKey)
			if err := obj.RunShow("remove_node", orderer.NodeType, orderer.NodeName); err != nil {
				return err
			}
		}
	}
	return nil
}

func RunAddNode(nodename string) error {
	if nodename == "" {
		return fmt.Errorf("nodename is empty")
	}
	for _, peer := range GlobalConfig.Peers {
		if peer.NodeName == nodename {
			//启动节点
			obj := NewFabCmd("add_node.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			if err := obj.RunShow("start_node", peer.NodeName, ConfigDir()); err != nil {
				return err
			}
		}
	}
	return nil
}

func HandleNode(stringType string, isStart bool) error {
	if stringType == TypeExplorer {
		operationExplorer(isStart)
	}
	if stringType == TypeApi {
		operationApiserver(isStart)
	}
	if isStart {
		if err := WriteHost(); err != nil {
			return err
		}
	}

	if stringType == "all" || stringType == TypeKafka {
		ForeachNode(GlobalConfig.Kafkas, isStart)
	}
	if stringType == "all" || stringType == TypeZookeeper {
		ForeachNode(GlobalConfig.Zookeepers, isStart)
	}
	if stringType == "all" || stringType == TypeOrder {
		ForeachNode(GlobalConfig.Orderers, isStart)
	}
	if stringType == "all" || stringType == TypePeer {
		ForeachNode(GlobalConfig.Peers, isStart)
	}
	if stringType == "all" || stringType == TypeCa {
		if GlobalConfig.CaType == "fabric-ca" {
			ForeachNode(GlobalConfig.Cas, isStart)
		} else {
			fmt.Println("Type isn't fabric-ca don't handle ca node")
		}
	}

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

	for _, obj := range HostMapList {
		obj := NewFabCmd("add_node.py", obj.Ip, obj.SshUserName, obj.SshPwd, obj.SshPort, obj.SshKey)
		if err := obj.RunShow("check_node"); err != nil {
			return err
		}
	}

	return nil
}
