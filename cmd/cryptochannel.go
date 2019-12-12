package cmd

import (
	"fmt"
	"github.com/peersafe/deployFabricTool/tpl"
	"os"
	"strings"
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
			orgName := fmt.Sprintf("org%s.%s", v.OrgId, GlobalConfig.Domain)
			err := obj.RunShow("generate_certs_to_ca", BinPath(), ConfigDir(), GlobalConfig.CryptoType, v.NodeType, v.NodeName, orgName, v.CaUrl, v.CaUrl, v.AdminName, v.AdminPw)
			if err != nil {
				return err
			}
		}
		for _, v := range GlobalConfig.Orderers {
			orgName := fmt.Sprintf("ord%s.%s", v.OrgId, GlobalConfig.Domain)
			err := obj.RunShow("generate_certs_to_ca", BinPath(), ConfigDir(), GlobalConfig.CryptoType, v.NodeType, v.NodeName, orgName, v.CaUrl, v.CaUrl, v.AdminName, v.AdminPw)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("-----Apply cert by Cryptogen---")
		err := obj.RunShow("generate_certs", BinPath(), ConfigDir(), ConfigDir(), GlobalConfig.CryptoType)
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
		return tpl.Handler(GlobalConfig, TplPath(TplCryptoConfig), ConfigDir()+"crypto-config.yaml")
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
			peer.DefaultNetwork = strings.Replace(peer.NodeName, ".", "", -1)
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
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		break
	}
	err := obj.RunShow("create_channel", BinPath(), ConfigDir(), ChannelPath(), channelName, ordererAddress, GlobalConfig.Domain, GlobalConfig.CryptoType)
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
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		break
	}
	for _, peer := range GlobalConfig.Peers {
		if peer.Id == "0" {
			obj := NewFabCmd("create_channel.py", peer.Ip, peer.SshUserName, peer.SshPwd)
			mspid := peer.OrgId
			err := obj.RunShow("update_anchor", BinPath(), ConfigDir(), ChannelPath(), channelName, mspid, ordererAddress, GlobalConfig.Domain, GlobalConfig.CryptoType)
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

func PutCryptoConfig() error {
	var wg sync.WaitGroup
	putCrypto := func(ip, sshuser, sshpwd, cfg, nodeTy, nodeName, orgName string, w1 *sync.WaitGroup) {
		obj := NewFabCmd("apply_cert.py", ip, sshuser, sshpwd)
		err := obj.RunShow("put_cryptoconfig", cfg, nodeTy, nodeName, orgName)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer w1.Done()
	}
	for _, kafka := range GlobalConfig.Kafkas {
		wg.Add(1)
		go putCrypto(kafka.Ip, kafka.SshUserName, kafka.SshPwd, ConfigDir(), TypeKafka, kafka.NodeName, "", &wg)
	}
	for _, zk := range GlobalConfig.Zookeepers {
		wg.Add(1)
		go putCrypto(zk.Ip, zk.SshUserName, zk.SshPwd, ConfigDir(), TypeZookeeper, zk.NodeName, "", &wg)
	}
	for _, ord := range GlobalConfig.Orderers {
		wg.Add(1)
		orgName := fmt.Sprintf("ord%s.%s", ord.OrgId, GlobalConfig.Domain)
		go putCrypto(ord.Ip, ord.SshUserName, ord.SshPwd, ConfigDir(), TypeOrder, ord.NodeName, orgName, &wg)
	}
	for _, peer := range GlobalConfig.Peers {
		wg.Add(1)
		orgName := fmt.Sprintf("org%s.%s", peer.OrgId, GlobalConfig.Domain)
		go putCrypto(peer.Ip, peer.SshUserName, peer.SshPwd, ConfigDir(), TypePeer, peer.NodeName, orgName, &wg)
	}
	wg.Wait()
	return nil
}
