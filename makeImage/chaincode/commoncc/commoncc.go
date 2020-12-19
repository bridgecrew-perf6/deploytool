package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type ResponseMsg struct {
	ErrCode string `json:"errcode"` //错误ID （"0"表示成功,"100"数据未发现，"999"为失败）
	ErrMsg  string `json:"errmsg"`  //错误消息
	Data    string `json:"data"`    //正文(json格式字符串)
	TxId    string `json:"txid"`    //交易id
}

const ERR_CODE_DATA_NOT_FOUND = "100"
const ERR_CODE = "999"
const SUCCESS_CODE = "0"

// CommonCC Chaincode implementation
type CommonCC struct {
}

func (t *CommonCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("CommonCC init")
	return shim.Success(nil)
}

func (t *CommonCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("Invoke function=%v,args=%v\n", function, args)
	if function == "saveData" {
		return t.saveData(stub, args)
	} else if function == "getDataByMsgSn" {
		return t.getDataByMsgSn(stub, args)
	}
	return response(ERR_CODE, "function not found", "", "")
}

//saveData
func (t *CommonCC) saveData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]
	data := args[1]
	if key == "" {
		str := "args[0]-key is empty"
		return response(ERR_CODE, str, "", "")
	}
	if data == "" {
		str := "args[1]-data is empty"
		return response(ERR_CODE, str, "", "")
	}
	txId := stub.GetTxID()
	if err := stub.PutState(key, []byte(data)); err != nil {
		str := "saveData key: [" + key + "] " + err.Error()
		return response(ERR_CODE, str, "", "")
	}
	if err := stub.PutState("txId_"+key, []byte(txId)); err != nil {
		str := "saveData key txId: [" + key + "] " + err.Error()
		return response(ERR_CODE, str, "", "")
	}
	return response(SUCCESS_CODE, "", "", txId)
}

//getDataByMsgSn
func (t *CommonCC) getDataByMsgSn(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	key := args[0]
	if key == "" {
		str := "args[0] key is empty"
		return response(ERR_CODE, str, "", "")
	}
	msgDataBytes, err := stub.GetState(key)
	if err != nil {
		str := "key: [" + key + "] " + err.Error()
		return response(ERR_CODE_DATA_NOT_FOUND, str, "", "")
	}
	if msgDataBytes == nil {
		str := "key: [" + key + "] not found"
		return response(ERR_CODE_DATA_NOT_FOUND, str, "", "")
	}
	txId, err := stub.GetState("txId_" + key)
	if err != nil {
		str := "txId:" + err.Error()
		return response(ERR_CODE_DATA_NOT_FOUND, str, "", "")
	}
	return response(SUCCESS_CODE, "", string(msgDataBytes), string(txId))
}

func response(errCode, errMsg, data, txId string) pb.Response {
	var responseMsg ResponseMsg
	responseMsg.ErrCode = errCode
	responseMsg.ErrMsg = errMsg
	responseMsg.Data = data
	responseMsg.TxId = txId
	responseMsgByte, _ := json.Marshal(responseMsg)
	return shim.Success(responseMsgByte)
}

func main() {
	err := shim.Start(new(CommonCC))
	if err != nil {
		fmt.Printf("Error starting CommonCC chaincode: %s", err)
	}
}
