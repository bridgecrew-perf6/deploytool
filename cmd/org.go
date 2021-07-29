package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/peersafe/deployFabricTool/tpl"
	"strings"
)

//生成新组织configtx配置文件
func CreateNewOrgConfigTxFile(orgid string) error {
	if orgid == "" {
		return fmt.Errorf("orgid is empty")
	}
	if err := CheckOrgNameIsExist(orgid) ; err != nil {
		return err
	}
	var orgInfo ConfigOrgInfo
	orgInfo.MspID = orgid

	for _, peer := range GlobalConfig.Peers {
		if peer.OrgId == orgid && peer.Id == "0" {
			orgInfo.MspPath = ConfigDir() + fmt.Sprintf("crypto-config/peerOrganizations/%s.%s/msp", orgid, peer.Domain)
			orgInfo.PeerAddress = peer.NodeName
			orgInfo.PeerPort = peer.ExternalPort
			break
		}
	}

	return tpl.Handler(orgInfo, TplPath(TplConfigtxOrg), UpdateConfigDir()+"configtx.yaml")
}

//添加组织到配置块
func AddOrgToConfigBlock(orgid, channelName string) error {
	if orgid == "" {
		return fmt.Errorf("orgid is empty")
	}
	if err := CheckOrgNameIsExist(orgid) ; err != nil {
		return err
	}
	if channelName == "" {
		return fmt.Errorf("channelName is empty")
	}
	ordererAddress := ""
	order_tls_path := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/orderer0.%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, dirPath, dirPath)
		break
	}
	type OrgInfo struct {
		MspId      string `json:"mspid"`
		OrgName    string `json:"orgname"`
		AnchorPeer string `json:"anchorpeer"`
		Port       string `json:"port"`
	}
	var OrgInfoList []*OrgInfo
	var info OrgInfo
	info.MspId = orgid
	info.OrgName = orgid
	for _, peer := range GlobalConfig.Peers {
		if peer.OrgId == orgid && peer.Id == "0" {
			info.AnchorPeer = peer.NodeName
			info.Port = peer.ExternalPort
			OrgInfoList = append(OrgInfoList, &info)
			break
		}
	}
	str_value, _ := json.Marshal(OrgInfoList)
	orgListStr := strings.Replace(string(str_value), ",", `\,`, -1)
	cmd := NewLocalFabCmd("add_org.py")
	OperationOrgId := ""
	for _, peer := range GlobalConfig.Peers {
		OperationOrgId = peer.OrgId
		break
	}
	err := cmd.RunShow("add_org_new", BinPath(), ConfigDir(), OperationOrgId, orgListStr, ordererAddress, order_tls_path, GlobalConfig.Domain, channelName)
	if err != nil {
		return err
	}
	return nil
}

//从配置块删除组织
func RmOrgFromConfigBlock(orgid, channelName string) error {
	if orgid == "" {
		return fmt.Errorf("orgid is empty")
	}
	if err := CheckOrgNameIsExist(orgid) ; err != nil {
		return err
	}
	if channelName == "" {
		return fmt.Errorf("channelName is empty")
	}
	ordererAddress := ""
	order_tls_path := ""
	for _, ord := range GlobalConfig.Orderers {
		ordererAddress = fmt.Sprintf("%s:%s", ord.NodeName, ord.ExternalPort)
		dirPath := fmt.Sprintf("%s.%s", ord.OrgId, GlobalConfig.Domain)
		order_tls_path = ConfigDir() + fmt.Sprintf("crypto-config/ordererOrganizations/%s/orderers/orderer0.%s/msp/tlscacerts/tlsca.%s-cert.pem", dirPath, dirPath, dirPath)
		break
	}
	type OrgInfo struct {
		OrgId string `json:"orgid"`
	}
	var OrgInfoList []*OrgInfo
	var info OrgInfo
	info.OrgId = orgid
	OrgInfoList = append(OrgInfoList,&info)
	//check orgid exist
	isExist := false
	for _, peer := range GlobalConfig.Peers {
		if peer.OrgId == orgid {
			isExist = true
			break
		}
	}
	if !isExist {
		return fmt.Errorf("the orgid: %s ,not exist node.json file", orgid)
	}
	str_value, _ := json.Marshal(OrgInfoList)
	orgListStr := strings.Replace(string(str_value), ",", `\,`, -1)
	cmd := NewLocalFabCmd("remove_org.py")
	OperationOrgId := ""
	for _, peer := range GlobalConfig.Peers {
		OperationOrgId = peer.OrgId
		break
	}
	err := cmd.RunShow("delete_org_new", BinPath(), ConfigDir(), OperationOrgId, orgListStr, ordererAddress, order_tls_path, GlobalConfig.Domain, channelName)
	if err != nil {
		return err
	}
	return nil
}
