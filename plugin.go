package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface.
func New() shim.Chaincode {
	fmt.Println("[honeybadgerscc] Hello")
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
	return shim.Success(nil)
}

// StartHbmpc start honeybadgerMPC
func StartHbmpc(nodeID string) {
	cmd := exec.Command("python3.7", "-m", "honeybadgermpc.secretshare_hbavsslight", nodeID, "conf/hbavss.hyper.ini")
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	errmsg := cmd.Run()
	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
}

// Invoke implements the chaincode shim interface
func (s *honeybadgerscc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "hbmpc" && len(args) >= 1 {
		fmt.Println("[honeybadgerscc]HBMPC")
		StartHbmpc(args[0])
		return shim.Success([]byte("Started HBMPC-hbavss"))
	} else { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	v, e := stub.GetState("a")
	if e != nil {
		fmt.Println("ERRORR")
	} else {
		fmt.Println("Not err", v)
	}
	fmt.Println("value is " + args[1])
	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func main() {}
