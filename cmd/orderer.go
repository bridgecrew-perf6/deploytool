package cmd

import (
	"fmt"
)

//添加orderer到配置块
func HandleOrderToConfigBlock(nodeName, channelName, operation string) error {
	if nodeName == "" {
		return fmt.Errorf("nodeName is empty")
	}
	if channelName == "" {
		return fmt.Errorf("channelName is empty")
	}
	if err := CheckNodeNameIsExist(nodeName) ; err != nil {
		return err
	}
	ordererAddress := ""
	order_tls_path := ""
	orderer_mspId := ""
	for _, ord := range GlobalConfig.Orderers {
		orderer_mspId = ord.OrgId
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, ord.NodeName, dirPath)
		break
	}

	cmd := NewLocalFabCmd("handle_orderer.py")
	OperationPort := ""
	for _, orderer := range GlobalConfig.Orderers {
		if nodeName == orderer.NodeName {
			OperationPort = orderer.ExternalPort
			break
		}
	}
	if OperationPort == "" {
		return fmt.Errorf("param invalided, node_name: %s not exist\n",nodeName)
	}
	err := cmd.RunShow("handle_orderer", BinPath(), ConfigDir(), nodeName, OperationPort, orderer_mspId, ordererAddress, order_tls_path, GlobalConfig.Domain, channelName, operation)
	if err != nil {
		return err
	}
	return nil
}


//更新创世块
func UpdateGenesisBlock() error {
	ordererAddress := ""
	order_tls_path := ""
	orderer_mspId := ""
	for _, ord := range GlobalConfig.Orderers {
		orderer_mspId = ord.OrgId
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, ord.NodeName, dirPath)
		break
	}

	cmd := NewLocalFabCmd("handle_orderer.py")
	err := cmd.RunShow("update_genesis_block", BinPath(), ConfigDir(), orderer_mspId, ordererAddress, order_tls_path, GlobalConfig.Domain)
	if err != nil {
		return err
	}
	return nil
}
