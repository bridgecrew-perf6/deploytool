package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	TplZookeeper    = "zookeeper.tpl"
	TplKafka        = "kafka.tpl"
	TplOrderer      = "orderer.tpl"
	TplPeer         = "peer.tpl"
	TplCryptoConfig = "crypto-config.tpl"
	TplConfigtx     = "configtx.tpl"
	TplApiClient    = "apiclient.tpl"
	TplApiDocker    = "apidocker.tpl"
	TplEventClient  = "eventclient.tpl"

	TypePeer         = "peer"
	TypeOrder        = "orderer"
	TypeKafka        = "kafka"
	TypeZookeeper    = "zookeeper"
	TypeApi          = "api"
	ZabbixServerIp   = "zabbix_server_ip"
	ZabbixServerPort = "zabbix_server_port"
	ZabbixAgentIp    = "zabbix_agent_ip"
	List             = "list"
	NodeType         = "node_type"
	ZkId             = "zk_id"
	IP               = "ip"
	APIIP            = "apiip"
	JMETER           = "jmeter"
	PeerId           = "peer_id"
	OrgId            = "org_id"
	OrderId          = "order_id"
	KfkId            = "kfk_id"
	ChanCounts       = "chan_counts"
)

var GlobalConfig *ConfigObj

type ConfigObj struct {
	FabricVersion  string         `json:"fabricVersion"`
	TestArgs       string         `json:"testArgs"`
	CCInit         string         `json:"ccInit"`
	CCPolicy       string         `json:"ccPolicy"`
	CCPath         string         `json:"ccPath"`
	CCName         string         `json:"ccName"`
	CCVersion      string         `json:"ccVersion"`
	CCInstallType  string         `json:"ccInstallType"`
	ConsensusType  string         `json:"consensusType"`
	BatchTime      string         `json:"batchTime"`
	BatchPreferred string         `json:"batchPreferred"`
	BatchSize      int            `json:"batchSize"`
	Zookeepers     []NodeObj      `json:"zookeepers"`
	Kafkas         []NodeObj      `json:"kafkas"`
	OrdList        map[string]int `json:"ordList"`
	OrgList        map[string]int `json:"orgList"`
	Expand
}

type Expand struct {
	SshUserName    string    `json:"sshUserName"`
	SshPwd         string    `json:"sshPwd"`
	SshKey         string    `json:"sshKey"`
	Log            string    `json:"log"`
	UseCouchdb     string    `json:"useCouchdb"`
	Domain         string    `json:"domain"`
	Orderers       []NodeObj `json:"orderers"`
	Peers          []NodeObj `json:"peers"`
	ImageTag       string    `json:"imageTag"`
	ImagePre       string    `json:"imagePre"`
	MountPath      string    `json:"mountPath"`
	CryptoType     string    `json:"cryptoType"`
	DefaultNetwork string    `json:"defaultNetwork"`
}

type NodeObj struct {
	Ip               string   `json:"ip"`
	ApiIp            string   `json:"apiIp"`
	Id               string   `json:"id"`
	OrgId            string   `json:"orgId"`
	Ports            []string `json:"ports"`
	ExternalPort     string   `json:"externalPort"`
	BootStrapAddress string   `json:"bootStrapAddress"`
	Expand
}

func ConfigDir() string {
	return os.Getenv("PWD") + "/config/"
}

func InputDir() string {
	return os.Getenv("PWD") + "/data/"
}

func TplPath(name string) string {
	return fmt.Sprintf("%s/templates/%s/%s", os.Getenv("PWD"), GlobalConfig.FabricVersion, name)
}

func BinPath() string {
	return fmt.Sprintf("%s/bin/%s/", os.Getenv("PWD"), GlobalConfig.FabricVersion)
}

func ChannelPath() string {
	return os.Getenv("PWD") + "/config/channel-artifacts/"
}

func ImagePath() string {
	return os.Getenv("PWD") + "/images/"
}

func ScriptPath() string {
	return os.Getenv("PWD") + "/scripts/"
}

func ParseJson(jsonfile string) (*ConfigObj, error) {
	var obj ConfigObj
	file := InputDir() + jsonfile
	fmt.Printf("json file %s\n", file)
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return &obj, err
	}

	obj.OrdList = make(map[string]int)
	obj.OrgList = make(map[string]int)

	peer0BootStrapMap := make(map[string]string)
	otherPeerBootStrapMap := make(map[string]string)

	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return &obj, err
	}
	for i, v := range obj.Peers {
		obj.OrgList[v.OrgId] = obj.OrgList[v.OrgId] + 1
		extPort, err := findExternalPort(v.Ports, "7051")
		if err != nil {
			return &obj, err
		}

		obj.Peers[i].ExternalPort = extPort
		if v.Id == "0" {
			otherPeerBootStrapMap[v.OrgId] = fmt.Sprintf("peer%s.org%s.%s:%s", v.Id, v.OrgId, obj.Domain, extPort)
		} else if peer0BootStrapMap[v.OrgId] == "" {
			peer0BootStrapMap[v.OrgId] = fmt.Sprintf("peer%s.org%s.%s:%s", v.Id, v.OrgId, obj.Domain, extPort)
		}
	}
	for i, v := range obj.Peers {
		if v.Id == "0" {
			obj.Peers[i].BootStrapAddress = peer0BootStrapMap[v.OrgId]
		} else {
			obj.Peers[i].BootStrapAddress = otherPeerBootStrapMap[v.OrgId]
		}
	}
	for i, v := range obj.Orderers {
		obj.OrdList[v.OrgId] = obj.OrdList[v.OrgId] + 1
		extPort, err := findExternalPort(v.Ports, "7050")
		if err != nil {
			return &obj, err
		}
		obj.Orderers[i].ExternalPort = extPort
	}
	for i, v := range obj.Kafkas {
		extPort, err := findExternalPort(v.Ports, "9092")
		if err != nil {
			return &obj, err
		}
		obj.Kafkas[i].ExternalPort = extPort
	}
	if obj.ImagePre == "" {
		obj.ImagePre = "peersafes"
	}
	if obj.MountPath == "" {
		obj.MountPath = "/data"
	}

	//fmt.Printf("config obj is %#v\n", obj)
	return &obj, nil
}

func GetJsonMap(jsonfile string) map[string]interface{} {
	var inputData map[string]interface{}
	var jsonData []byte
	var err error

	inputfile := InputDir() + jsonfile
	jsonData, err = ioutil.ReadFile(inputfile)
	if err != nil {
		return inputData
	}
	err = json.Unmarshal(jsonData, &inputData)
	if err != nil {
		return inputData
	}
	return inputData
}

func findExternalPort(list []string, destPort string) (string, error) {
	for _, v := range list {
		curLine := strings.Split(v, ":")
		if len(curLine) != 2 {
			return "", fmt.Errorf("findExternalPort err %s", v)
		}
		if curLine[1] == destPort {
			return curLine[0], nil
		}
	}
	return "", fmt.Errorf("findExternalPort err destPort %s not exist ", destPort)
}
