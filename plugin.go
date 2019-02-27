package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/syndtr/goleveldb/leveldb"
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
func StartHbmpc(nodeID string) string {
	cmd := exec.Command("python3.7", "-m", "honeybadgermpc.secretshare_hbavsslight", nodeID, "conf/hbavss.hyper.ini")
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()
	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
	lines := strings.Split(outb.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "[INFO]: Output available ") {
			// Very hacky way of doing this.
			shareParts := strings.Split(line, "[INFO]: Output available ")
			if len(shareParts) >= 2 {
				share := shareParts[1]
				fmt.Println("The share is ", share)
				return share
			}
		}
	}
	return "None"
	// fmt.Println("out : ", outb.String(), "err: ", errb.String())
}

func dbPut(key string, value string) {
	db, err := leveldb.OpenFile("db", nil)
	err = db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		fmt.Println("Error writing to database")
	}
	defer db.Close()
}

func dbGet(key string) string {
	db, err := leveldb.OpenFile("db", nil)
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Error getting from database")
	}
	defer db.Close()
	return string(data)
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
		share := StartHbmpc(args[0])
		if share != "None" {
			err := stub.PutState("key", []byte(share))
			dbPut("key", share)
			if err != nil {
				fmt.Errorf("Failed to set asset: %s", args[0])
			}
		}
		//err = db.Put([]byte("key"), []byte(share), nil)
		//if err != nil {
		//	fmt.Println("There has been an error writing to db")
		//	return shim.Error(err.Error())
		//}
		return shim.Success([]byte("Started HBMPC-hbavss"))
	} else if fn == "get" { // assume 'get' even if fn is nil
		result, err = get(stub, args)
	} else if fn == "test" {
		fmt.Println("[honetbadgerscc] I am here starting ")
	} else if fn == "getKey" {
		share := dbGet("key")
		//data, err := db.Get([]byte("key"), nil)
		//if err != nil {
		//	fmt.Println("error getting data ")
		//}
		//fmt.Println(data)
		return shim.Success([]byte(share))
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
