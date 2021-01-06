/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// TestCC Chaincode implementation
type TestCC struct {
}

func (t *TestCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("TestCC Init")
	return shim.Success(nil)
}

func (t *TestCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "SaveData" {
		return t.invoke(stub, args)
	} else if function == "GetData" {
		return t.query(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"SaveData\" \"GetData\"")
}

func (t *TestCC) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.PutState(stub.GetTxID(), []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *TestCC) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	Avalbytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("get state error")
	}

	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(TestCC))
	if err != nil {
		fmt.Printf("Error starting TestCC chaincode: %s", err)
	}
}
