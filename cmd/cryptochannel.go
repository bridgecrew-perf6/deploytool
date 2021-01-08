package cmd

import (
	"fmt"
	"github.com/peersafe/deployFabricTool/tpl"
	"os"
	"sync"
)

func CreateCert() error {
	obj := NewLocalFabCmd("apply_cert.py")
	if err := os.RemoveAll(ConfigDir() + "/crypto-config"); err != nil {
		return err
	}
	if err := os.RemoveAll(ConfigDir() + "/cadata"); err != nil {
		return err
	}
	if err := os.RemoveAll(ConfigDir() + "/tlscadata"); err != nil {
		return err
	}
	if GlobalConfig.CaType == "fabric-ca" {
		fmt.Println("-----Apply cert by Fabric-ca---")
		for _, v := range GlobalConfig.Peers {
			orgName := fmt.Sprintf("%s.%s", v.OrgId, GlobalConfig.Domain)
			err := obj.RunShow("generate_certs_to_ca", BinPath(), ConfigDir(), GlobalConfig.CryptoType, v.NodeType, v.NodeName, orgName, v.CaUrl, v.CaUrl, v.AdminName, v.AdminPw)
			if err != nil {
				return err
			}
		}
		for _, v := range GlobalConfig.Orderers {
			orgName := fmt.Sprintf("%s.%s", v.OrgId, GlobalConfig.Domain)
			err := obj.RunShow("generate_certs_to_ca", BinPath(), ConfigDir(), GlobalConfig.CryptoType, v.NodeType, v.NodeName, orgName, v.CaUrl, v.CaUrl, v.AdminName, v.AdminPw)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("-----Apply cert by Cryptogen---")
		for name, _ := range GlobalConfig.OrdList {
			configFile := ConfigDir() + name + "-crypto-config-orderer.yaml"
			if err := obj.RunShow("generate_certs", BinPath(), configFile, ConfigDir(), GlobalConfig.CryptoType); err != nil {
				return err
			}
		}
		for name, _ := range GlobalConfig.OrgList {
			configFile := ConfigDir() + name + "-crypto-config-peer.yaml"
			if err := obj.RunShow("generate_certs", BinPath(), configFile, ConfigDir(), GlobalConfig.CryptoType); err != nil {
				return err
			}
		}
	}
	return nil
}

func makeExploreYaml() error {
	for _, e := range GlobalConfig.Explorers {
		e.Domain = GlobalConfig.Domain
		e.Log = GlobalConfig.Log
		e.MountPath = GlobalConfig.MountPath
		e.CryptoType = GlobalConfig.CryptoType
		e.CCName = GlobalConfig.CCName
		ParentPath := ConfigDir() + e.NodeName + "/"
		outFile := "block_fabric_explorer.yaml"
		err := tpl.Handler(e, ExplorerTplPath(TplExplorer), ParentPath+outFile)
		if err != nil {
			return err
		}
		outFile = "client_sdk.yaml"
		err = tpl.Handler(e, ExplorerTplPath(TplClient), ParentPath+outFile)
		if err != nil {
			return err
		}
		outFile = "registerApi.js"
		err = tpl.Handler(e, ExplorerTplPath(TplRegister), ParentPath+outFile)
		if err != nil {
			return err
		}
		outFile = "mysql.sql"
		err = tpl.Handler(e, ExplorerTplPath(outFile), ParentPath+"mysql_init/"+outFile)
		if err != nil {
			return err
		}
		outFile = "mysqld.cnf"
		err = tpl.Handler(e, ExplorerTplPath(outFile), ParentPath+outFile)
		if err != nil {
			return err
		}

	}
	return nil
}

func CreateYamlByJson(strType string) error {
	if strType == "configtx" {
		return tpl.Handler(GlobalConfig, TplPath(TplConfigtx), ConfigDir()+"configtx.yaml")
	} else if strType == "crypto-config" {
		for name, counts := range GlobalConfig.OrdList {
			outFile := ConfigDir() + name + "-crypto-config-orderer.yaml"
			orgObj := OrgObj{name, counts, GlobalConfig.Domain}
			if err := tpl.Handler(orgObj, TplPath(TplOrdererCryptoConfig), outFile); err != nil {
				return err
			}
		}
		for name, counts := range GlobalConfig.OrgList {
			outFile := ConfigDir() + name + "-crypto-config-peer.yaml"
			orgObj := OrgObj{name, counts, GlobalConfig.Domain}
			if err := tpl.Handler(orgObj, TplPath(TplPeerCryptoConfig), outFile); err != nil {
				return err
			}
		}
	} else if strType == TypeExplorer {
		return makeExploreYaml()
	} else if strType == TypeApi {
		for _, api := range GlobalConfig.Apiservers {
			api.Domain = GlobalConfig.Domain
			api.CryptoType = GlobalConfig.CryptoType
			api.Log = GlobalConfig.Log
			outfile := ConfigDir() + "client_sdk"
			if err := tpl.Handler(api, TplCommonPath(TplApiClient), outfile+".yaml"); err != nil {
				return err
			}
			outfile = ConfigDir() + api.NodeName
			if err := tpl.Handler(api, TplCommonPath(TplApiDocker), outfile+".yaml"); err != nil {
				return err
			}
		}
	} else if strType == "node" || strType == "client" {
		for _, ca := range GlobalConfig.Cas {
			CopyConfig(&ca)
			outfile := ConfigDir() + ca.NodeName
			if err := tpl.Handler(ca, TplPath(TplCa), outfile+".yaml"); err != nil {
				return err
			}
		}
		for _, ord := range GlobalConfig.Orderers {
			CopyConfig(&ord)
			outfile := ConfigDir() + ord.NodeName
			if err := tpl.Handler(ord, TplPath(TplOrderer), outfile+".yaml"); err != nil {
				return err
			}
		}
		for _, peer := range GlobalConfig.Peers {
			CopyConfig(&peer)
			outfile := ConfigDir() + peer.NodeName
			if err := tpl.Handler(peer, TplPath(TplPeer), outfile+".yaml"); err != nil {
				return err
			}
		}
		for _, kafka := range GlobalConfig.Kafkas {
			outfile := ConfigDir() + kafka.NodeName
			if err := tpl.Handler(kafka, TplPath(TplPeer), outfile+".yaml"); err != nil {
				return err
			}
		}
		for _, zk := range GlobalConfig.Zookeepers {
			outfile := ConfigDir() + zk.NodeName
			if err := tpl.Handler(zk, TplPath(TplZookeeper), outfile+".yaml"); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("%s not exist", strType)
	}
	return nil
}

func CreateGenesisBlock() error {
	model := ""
	if GlobalConfig.ConsensusType == "solo" {
		model = "OrgsOrdererGenesis"
	} else if GlobalConfig.ConsensusType == "kafka" {
		model = "SampleDevModeKafka"
	} else if GlobalConfig.ConsensusType == "raft" {
		model = "SampleMultiNodeEtcdRaft"
	} else {
		return fmt.Errorf("ConsensusType %s unknow", GlobalConfig.ConsensusType)
	}
	obj := NewLocalFabCmd("apply_cert.py")
	err := obj.RunShow("generate_genesis_block", model, BinPath(), ConfigDir(), ConfigDir(), GlobalConfig.CryptoType)
	if err != nil {
		return err
	}
	return nil
}

func CreateChannel(channelName string) error {
	if channelName == "" {
		return fmt.Errorf("channel name is nil")
	}
	obj := NewLocalFabCmd("create_channel.py")
	ordererAddress := ""
	order_tls_path := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/orderer0.%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, dirPath, dirPath)
		break
	}
	OrgId := ""
	for _, peer := range GlobalConfig.Peers {
		OrgId = peer.OrgId
		break
	}
	err := obj.RunShow("create_channel", BinPath(), ConfigDir(), ChannelPath(), channelName, OrgId, ordererAddress, order_tls_path, GlobalConfig.Domain, GlobalConfig.CryptoType)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAnchor(channelName string) error {
	if channelName == "" {
		return fmt.Errorf("channel name is nil")
	}
	ordererAddress := ""
	order_tls_path := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/orderer0.%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, dirPath, dirPath)
	}
	for _, peer := range GlobalConfig.Peers {
		if peer.Id == "0" {
			obj := NewFabCmd("create_channel.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			mspid := peer.OrgId
			err := obj.RunShow("update_anchor", BinPath(), ConfigDir(), ChannelPath(), channelName, mspid, ordererAddress, order_tls_path,GlobalConfig.Domain, GlobalConfig.CryptoType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyConfig(obj *NodeObj) {
	obj.Domain = GlobalConfig.Domain
	obj.Log = GlobalConfig.Log
	obj.UseCouchdb = GlobalConfig.UseCouchdb
	obj.ImageTag = GlobalConfig.ImageTag
	obj.ImagePre = GlobalConfig.ImagePre
	obj.MountPath = GlobalConfig.MountPath
	obj.CryptoType = GlobalConfig.CryptoType
}

func JoinChannel(channelName string) error {
	if channelName == "" {
		return fmt.Errorf("channel name is nil")
	}
	for _, peer := range GlobalConfig.Peers {
		peerAddress := fmt.Sprintf("%s:%s", peer.NodeName, peer.ExternalPort)
		obj := NewLocalFabCmd("create_channel.py")
		err := obj.RunShow("join_channel", BinPath(), ConfigDir(), ChannelPath(), channelName, peerAddress, peer.Id, peer.OrgId, GlobalConfig.Domain, GlobalConfig.CryptoType)
		if err != nil {
			return err
		}
	}
	return nil
}

func PutCryptoConfig(stringType string) error {
	var wg sync.WaitGroup
	putCrypto := func(ip, sshuser, sshpwd, sshport, sshkey, cfg, nodeTy, nodeName, orgName, certPeerName string, w1 *sync.WaitGroup) {
		obj := NewFabCmd("apply_cert.py", ip, sshuser, sshpwd, sshport, sshkey)
		err := obj.RunShow("put_cryptoconfig", cfg, nodeTy, nodeName, orgName, certPeerName)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer w1.Done()
	}
	if stringType == TypeExplorer {
		for _, exp := range GlobalConfig.Explorers {
			wg.Add(1)
			orgName := fmt.Sprintf("%s.%s", exp.OrgId, GlobalConfig.Domain)
			certPeerName := fmt.Sprintf("peer%s.%s.%s", exp.PeerId, exp.OrgId, GlobalConfig.Domain)
			putCrypto(exp.Ip, exp.SshUserName, exp.SshPwd, exp.SshPort, exp.SshKey, ConfigDir(), TypeExplorer, exp.NodeName, orgName, certPeerName, &wg)
		}
	} else {
		for _, api := range GlobalConfig.Apiservers {
			wg.Add(1)
			go putCrypto(api.Ip, api.SshUserName, api.SshPwd, api.SshPort, api.SshKey, ConfigDir(), TypeApi, api.NodeName, "", "", &wg)
		}
		for _, kafka := range GlobalConfig.Kafkas {
			wg.Add(1)
			go putCrypto(kafka.Ip, kafka.SshUserName, kafka.SshPwd, kafka.SshPort, kafka.SshKey, ConfigDir(), TypeKafka, kafka.NodeName, "", "", &wg)
		}
		for _, zk := range GlobalConfig.Zookeepers {
			wg.Add(1)
			go putCrypto(zk.Ip, zk.SshUserName, zk.SshPwd, zk.SshPort, zk.SshKey, ConfigDir(), TypeZookeeper, zk.NodeName, "", "", &wg)
		}
		for _, ord := range GlobalConfig.Orderers {
			wg.Add(1)
			orgName := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
			go putCrypto(ord.Ip, ord.SshUserName, ord.SshPwd, ord.SshPort, ord.SshKey, ConfigDir(), TypeOrder, ord.NodeName, orgName, "", &wg)
		}
		for _, peer := range GlobalConfig.Peers {
			wg.Add(1)
			orgName := fmt.Sprintf("%s.%s", peer.OrgId, GlobalConfig.Domain)
			go putCrypto(peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey, ConfigDir(), TypePeer, peer.NodeName, orgName, "", &wg)
		}
	}
	wg.Wait()
	return nil
}
