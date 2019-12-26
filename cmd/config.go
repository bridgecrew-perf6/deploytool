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
	TplCa           = "ca.tpl"
	TplPeer         = "peer.tpl"
	TplCryptoConfig = "crypto-config.tpl"
	TplConfigtx     = "configtx.tpl"
	TplApiClient    = "apiclient.tpl"
	TplApiDocker    = "apidocker.tpl"
	TplEventClient  = "eventclient.tpl"

	TypePeer         = "peer"
	TypeCa           = "ca"
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

var (
	GlobalConfig *ConfigObj
	HostMapList  = make(map[string]NodeObj)
)

type ConfigObj struct {
	FabricVersion  string         `json:"fabricVersion"`
	TestArgs       string         `json:"testArgs"`
	CaType         string         `json:"caType"`
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
	OrdList        map[string]int `json:"ordList"`
	OrgList        map[string]int `json:"orgList"`
	Zookeepers     []NodeObj      `json:"zookeepers"`
	Kafkas         []NodeObj      `json:"kafkas"`
	Orderers       []NodeObj      `json:"orderers"`
	Peers          []NodeObj      `json:"peers"`
	Cas            []NodeObj      `json:"cas"`
	TPLExpand
}

type TPLExpand struct {
	SshUserName    string `json:"sshUserName"`
	SshPwd         string `json:"sshPwd"`
	SshKey         string `json:"sshKey"`
	SshPort        string `json:"sshPort"`
	Log            string `json:"log"`
	UseCouchdb     string `json:"useCouchdb"`
	Domain         string `json:"domain"`
	ImageTag       string `json:"imageTag"`
	ImagePre       string `json:"imagePre"`
	MountPath      string `json:"mountPath"`
	CryptoType     string `json:"cryptoType"`
	DefaultNetwork string `json:"defaultNetwork"`
}

type NodeObj struct {
	Ip               string       `json:"ip"`
	ApiIp            string       `json:"apiIp"`
	Id               string       `json:"id"`
	NodeType         string       `json:"nodeType"`
	OrgId            string       `json:"orgId"`
	Ports            []string     `json:"ports"`
	ExternalPort     string       `json:"externalPort"`
	NodeName         string       `json:"nodeName"`
	ImageName        string       `json:"imageName"`
	CaUrl            string       `json:"caUrl"`
	CertType         string       `json:"certType"`
	AdminName        string       `json:"adminName"`
	AdminPw          string       `json:"adminPw"`
	BootStrapAddress string       `json:"bootStrapAddress"`
	ExtHosts         []ExtraHosts `json:"extra_hosts"`
	TPLExpand
}

var allPeerHostIp, allOrdererHostIp []ExtraHosts

type ExtraHosts struct {
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
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
	ordererCaMap := make(map[string]string)
	peerCaMap := make(map[string]string)

	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return &obj, err
	}
	if obj.CaType != "fabric-ca" {
		obj.Cas = []NodeObj{}
	}
	for i, v := range obj.Cas {
		extPort, err := findExternalPort(v.Ports, "7054")
		if err != nil {
			return &obj, err
		}
		if v.CertType == TypeOrder {
			ordererCaMap[v.OrgId] = fmt.Sprintf("%s:%s", v.Ip, extPort)
			obj.Cas[i].NodeName = fmt.Sprintf("ca.ord%s.%s", v.OrgId, obj.Domain)
			obj.Cas[i].AdminName = fmt.Sprintf("Admin@ord%s.%s", v.OrgId, obj.Domain)
		} else if v.CertType == TypePeer {
			peerCaMap[v.OrgId] = fmt.Sprintf("%s:%s", v.Ip, extPort)
			obj.Cas[i].NodeName = fmt.Sprintf("ca.org%s.%s", v.OrgId, obj.Domain)
			obj.Cas[i].AdminName = fmt.Sprintf("Admin@org%s.%s", v.OrgId, obj.Domain)
		}
		obj.Cas[i].NodeType = TypeCa
		obj.Cas[i].AdminPw = "adminpw"
		obj.Cas[i].ImageName = fmt.Sprintf("%s/fabric-ca:%s", obj.ImagePre, obj.ImageTag)
		HostMapList[v.Ip] = obj.Cas[i]
	}
	for i, v := range obj.Peers {
		obj.OrgList[v.OrgId] = obj.OrgList[v.OrgId] + 1

		extPort, err := findExternalPort(v.Ports, "7051")
		if err != nil {
			return &obj, err
		}
		obj.Peers[i].ExternalPort = extPort
		obj.Peers[i].AdminName = fmt.Sprintf("Admin@org%s.%s", v.OrgId, obj.Domain)
		obj.Peers[i].AdminPw = "adminpw"
		obj.Peers[i].CaUrl = peerCaMap[v.OrgId]
		if v.Id == "0" {
			otherPeerBootStrapMap[v.OrgId] = fmt.Sprintf("peer%s.org%s.%s:%s", v.Id, v.OrgId, obj.Domain, extPort)
		} else if peer0BootStrapMap[v.OrgId] == "" {
			peer0BootStrapMap[v.OrgId] = fmt.Sprintf("peer%s.org%s.%s:%s", v.Id, v.OrgId, obj.Domain, extPort)
		}
		obj.Peers[i].NodeName = fmt.Sprintf("peer%s.org%s.%s", v.Id, v.OrgId, obj.Domain)
		obj.Peers[i].NodeType = TypePeer
		allPeerHostIp = append(allPeerHostIp, ExtraHosts{obj.Peers[i].NodeName, v.Ip})
		obj.Peers[i].ImageName = fmt.Sprintf("%s/fabric-peer:%s", obj.ImagePre, obj.ImageTag)
		HostMapList[v.Ip] = obj.Peers[i]
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
		obj.Orderers[i].AdminName = fmt.Sprintf("Admin@ord%s.%s", v.OrgId, obj.Domain)
		obj.Orderers[i].AdminPw = "adminpw"
		obj.Orderers[i].CaUrl = ordererCaMap[v.OrgId]
		obj.Orderers[i].NodeType = TypeOrder
		obj.Orderers[i].NodeName = fmt.Sprintf("orderer%s.ord%s.%s", v.Id, v.OrgId, obj.Domain)
		allOrdererHostIp = append(allOrdererHostIp, ExtraHosts{obj.Orderers[i].NodeName, v.Ip})
		obj.Orderers[i].ImageName = fmt.Sprintf("%s/fabric-orderer:%s", obj.ImagePre, obj.ImageTag)
		HostMapList[v.Ip] = obj.Orderers[i]
	}
	for i, v := range obj.Kafkas {
		extPort, err := findExternalPort(v.Ports, "9092")
		if err != nil {
			return &obj, err
		}
		obj.Kafkas[i].ExternalPort = extPort
		obj.Kafkas[i].NodeType = TypeKafka
		obj.Kafkas[i].NodeName = fmt.Sprintf("kafka%s", v.Id)
		obj.Kafkas[i].ImageName = fmt.Sprintf("%s/fabric-kafka:%s", obj.ImagePre, obj.ImageTag)
		HostMapList[v.Ip] = obj.Kafkas[i]
	}

	for i, v := range obj.Zookeepers {
		extPort, err := findExternalPort(v.Ports, "2888")
		if err != nil {
			return &obj, err
		}
		obj.Zookeepers[i].NodeType = TypeZookeeper
		obj.Zookeepers[i].ExternalPort = extPort
		obj.Zookeepers[i].NodeName = fmt.Sprintf("zk%s", v.Id)
		obj.Zookeepers[i].ImageName = fmt.Sprintf("%s/fabric-zookeeper:%s", obj.ImagePre, obj.ImageTag)
		HostMapList[v.Ip] = obj.Zookeepers[i]
	}

	if obj.ImagePre == "" {
		obj.ImagePre = "peersafes"
	}
	if obj.MountPath == "" {
		obj.MountPath = "/data"
	}

	SetExtalHost(&obj)
	//fmt.Printf("config obj is %#v\n", obj)
	return &obj, nil
}

func SetExtalHost(obj *ConfigObj) {
	for i, v := range obj.Peers {
		obj.Peers[i].ExtHosts = append(obj.Peers[i].ExtHosts, allOrdererHostIp...)
		for _, item := range allPeerHostIp {
			if item.Domain != v.NodeName {
				obj.Peers[i].ExtHosts = append(obj.Peers[i].ExtHosts, item)
			}
		}
	}
	for i, v := range obj.Orderers {
		obj.Orderers[i].ExtHosts = []ExtraHosts{}
		for _, item := range allOrdererHostIp {
			if item.Domain != v.NodeName {
				obj.Orderers[i].ExtHosts = append(obj.Orderers[i].ExtHosts, item)
			}
		}
	}
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
