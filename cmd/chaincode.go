package cmd

import (
	"fmt"
	"sync"
)

func InstallChaincode(ccname, ccversion, channelName, ccpath string) error {
	if ccpath == "" {
		ccpath = GlobalConfig.CCPath
	}
	if ccversion == "" {
		ccversion = GlobalConfig.CCVersion
	}
	if ccname == "" {
		ccname = GlobalConfig.CCName
	}
	for _, peer := range GlobalConfig.Peers {
		//make cc pkg file even by chaincode path or pkg type
		obj := NewLocalFabCmd("chaincode.py")
		err := obj.RunShow("pkg_chaincode", GlobalConfig.FabricVersion, BinPath(), ConfigDir(), peer.OrgId, GlobalConfig.Domain, ccname, ccversion, ccpath, GlobalConfig.CCInstallType, GlobalConfig.CryptoType)
		if err != nil {
			return err
		}
		//only once
		break
	}

	var wg sync.WaitGroup
	for _, peer := range GlobalConfig.Peers {
		wg.Add(1)
		peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
		go func(binPath, configDir, peerAds, PeerId, OrgId, Pdn string) {
			defer wg.Done()
			obj := NewLocalFabCmd("chaincode.py")
			err := obj.RunShow("install_chaincode", GlobalConfig.FabricVersion, binPath, configDir, peerAds, PeerId, OrgId, Pdn, ccname, ccversion, ccpath, GlobalConfig.CCInstallType, GlobalConfig.CryptoType)
			if err != nil {
				panic(err)
			}
		}(BinPath(), ConfigDir(), peerAddress, peer.Id, peer.OrgId, GlobalConfig.Domain)
	}
	wg.Wait()

	//for _, peer := range GlobalConfig.Peers {
	//	peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
	//	obj := NewLocalFabCmd("chaincode.py")
	//	err := obj.RunShow("install_chaincode", GlobalConfig.FabricVersion, BinPath(), ConfigDir(), peerAddress, peer.Id, peer.OrgId, GlobalConfig.Domain, ccname, ccversion, ccpath, GlobalConfig.CCInstallType, GlobalConfig.CryptoType)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	ordererAddress := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("orderer%s.ord%s.%s:%s", ord.Id, ord.OrgId, GlobalConfig.Domain, ord.ExternalPort)
		break
	}
	if GlobalConfig.FabricVersion != "1.4" {
		for _, peer := range GlobalConfig.Peers {
			if peer.Id != "0" {
				continue
			}
			peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
			obj := NewLocalFabCmd("chaincode.py")
			if channelName == "" {
				panic("approve file channelName is empty")
			} else {
				err := obj.RunShow("approve_chaincode", BinPath(), ConfigDir(), peerAddress, ordererAddress, peer.Id, peer.OrgId, GlobalConfig.Domain, channelName, ccname, ccversion, GlobalConfig.CryptoType)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func RunChaincode(ccname, ccversion, channelName, opration string) error {
	if channelName == "" {
		return fmt.Errorf("channel is nil")
	}
	if ccversion == "" {
		ccversion = GlobalConfig.CCVersion
	}
	if ccname == "" {
		ccname = GlobalConfig.CCName
	}
	ordererAddress := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("orderer%s.ord%s.%s:%s", ord.Id, ord.OrgId, GlobalConfig.Domain, ord.ExternalPort)
		break
	}
	cmdParas := ""
	var wg sync.WaitGroup
	for _, peer := range GlobalConfig.Peers {
		peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
		if GlobalConfig.FabricVersion == "1.4" {
			wg.Add(1)
			go func(binPath, configDir, peerAds, PeerId, OrgId, Pdn string) {
				defer wg.Done()
				obj := NewFabCmd("chaincode.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
				initparam := fmt.Sprintf(`%s`, GlobalConfig.CCInit)
				policy := fmt.Sprintf("%s", GlobalConfig.CCPolicy)
				err := obj.RunShow("instantiate_chaincode", GlobalConfig.FabricVersion, BinPath(),
					opration, ConfigDir(), peerAds, ordererAddress, PeerId, OrgId, GlobalConfig.Domain,
					channelName, ccname, ccversion, initparam, policy, GlobalConfig.CryptoType, cmdParas)
				if err != nil {
					fmt.Println(err)
				}
			}(BinPath(), ConfigDir(), peerAddress, peer.Id, peer.OrgId, GlobalConfig.Domain)
		} else {
			if peer.Id == "0" {
				peerTlsCert := fmt.Sprintf("%s/crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/server.crt", ConfigDir(), peer.OrgId, peer.Domain, peer.Id, peer.OrgId, peer.Domain)
				cmdParas = cmdParas + fmt.Sprintf("  --peerAddresses %s --tlsRootCertFiles %s", peerAddress, peerTlsCert)
			}
		}
	}
	wg.Wait()
	if GlobalConfig.FabricVersion != "1.4" {
		for _, peer := range GlobalConfig.Peers {
			peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
			obj := NewFabCmd("chaincode.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			initparam := fmt.Sprintf(`%s`, GlobalConfig.CCInit)
			policy := fmt.Sprintf("%s", GlobalConfig.CCPolicy)
			err := obj.RunShow("instantiate_chaincode", GlobalConfig.FabricVersion, BinPath(), opration, ConfigDir(),
				peerAddress, ordererAddress, peer.Id, peer.OrgId, GlobalConfig.Domain, channelName,
				ccname, ccversion, initparam, policy, GlobalConfig.CryptoType, cmdParas)
			if err != nil {
				return err
			}
			err = obj.RunShow("test_chaincode", "2.0", "invoke", BinPath(), ConfigDir(), peerAddress, ordererAddress, peer.Id, peer.OrgId,
				GlobalConfig.Domain, channelName, ccname, initparam, GlobalConfig.CryptoType, cmdParas+" --isInit ")
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func TestChaincode(ccname, channelName, function, testArgs string) error {
	if channelName == "" {
		return fmt.Errorf("channel is nil")
	}
	if ccname == "" {
		ccname = GlobalConfig.CCName
	}
	if testArgs == "" {
		testArgs = GlobalConfig.TestArgs
	}
	ordererAddress := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("orderer%s.ord%s.%s:%s", ord.Id, ord.OrgId, GlobalConfig.Domain, ord.ExternalPort)
		break
	}
	cmdParas := ""
	var wg sync.WaitGroup
	for _, peer := range GlobalConfig.Peers {
		peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
		if GlobalConfig.FabricVersion == "1.4" {
			wg.Add(1)
			go func(binPath, configDir, peerAds, PeerId, OrgId, Pdn string) {
				defer wg.Done()
				obj := NewFabCmd("chaincode.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
				err := obj.RunShow("test_chaincode", "1.4", function, BinPath(), ConfigDir(), peerAds, ordererAddress, PeerId, OrgId, GlobalConfig.Domain, channelName, ccname, testArgs, GlobalConfig.CryptoType, "")
				if err != nil {
					fmt.Println(err)
				}
			}(BinPath(), ConfigDir(), peerAddress, peer.Id, peer.OrgId, GlobalConfig.Domain)
		} else {
			if peer.Id == "0" {
				peerTlsCert := fmt.Sprintf("%s/crypto-config/peerOrganizations/org%s.%s/peers/peer%s.org%s.%s/tls/server.crt", ConfigDir(), peer.OrgId, peer.Domain, peer.Id, peer.OrgId, peer.Domain)
				cmdParas = cmdParas + fmt.Sprintf("  --peerAddresses %s --tlsRootCertFiles %s", peerAddress, peerTlsCert)
			}
		}
	}
	wg.Wait()
	if GlobalConfig.FabricVersion != "1.4" {
		for _, peer := range GlobalConfig.Peers {
			peerAddress := fmt.Sprintf("peer%s.org%s.%s:%s", peer.Id, peer.OrgId, GlobalConfig.Domain, peer.ExternalPort)
			obj := NewFabCmd("chaincode.py", peer.Ip, peer.SshUserName, peer.SshPwd, peer.SshPort, peer.SshKey)
			testparam := fmt.Sprintf(`%s`, GlobalConfig.TestArgs)
			err := obj.RunShow("test_chaincode", "2.0", "invoke", BinPath(), ConfigDir(), peerAddress, ordererAddress, peer.Id, peer.OrgId,
				GlobalConfig.Domain, channelName, ccname, testparam, GlobalConfig.CryptoType, cmdParas)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
