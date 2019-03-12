package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface.
func New() shim.Chaincode {
	fmt.Println("[honeybadgerscc] Chaincode Started.")
	return &honeybadgerscc{}
}

type honeybadgerscc struct{}

// Init implements the chaincode shim interface
func (s *honeybadgerscc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("[honeybadgerscc]Init is called")
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	// Set up any variables or assets here by calling stub.PutState()
	// We store the key and the value on the ledger
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}

	err = stub.PutState("readToken", []byte("0"))
	if err != nil {
		return shim.Error("Unable to set readToken")
	}
	return shim.Success(nil)
}

// Invoke implements the chaincode shim interface
func (s *honeybadgerscc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	var err error
	if fn == "hbmpc" && len(args) >= 3 {
		fmt.Println("[honeybadgerscc]HBMPC")
		key := args[1]
		namespace := args[2]
		if !authenticateRequest(namespace, key) {
			return shim.Success([]byte("Failed to authenticate request"))
		}
		share := StartHbmpc(args[0])
		if share != "None" {
			err := stub.PutState(key, []byte("NA"))
			dbPut(key, share)
			if err != nil {
				fmt.Errorf("Failed to set asset: %s", args[0])
			}
		}
		return shim.Success([]byte("Ran HBMPC-hbavss"))
	} else if fn == "getKey" && len(args) >= 1 {
		key := args[0]
		share := dbGet(key)
		return shim.Success([]byte(share))
	} else if fn == "pubRecon" && len(args) >= 2 {
		nodeID := args[0]
		key := args[1]
		share := dbGet(key)
		value := StartPubRec(nodeID, share)
		err := stub.PutState(key, []byte(value))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(value))
	} else if fn == "initRecon" && len(args) >= 2 {
		key := args[0]
		namespace := args[0]
		res := InitiatePubRec(key, namespace)
		return shim.Success([]byte(res))
	} else if fn == "get" && len(args) > 1 {
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return shim.Success([]byte("Key does not exist yet"))
		}
		return shim.Success([]byte(value))
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte("Invalid endpoint"))
}

func main() {}
