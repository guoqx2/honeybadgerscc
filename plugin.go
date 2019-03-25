package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/syndtr/goleveldb/leveldb"
)

// New returns an implementation of the chaincode interface.
func New() shim.Chaincode {
	fmt.Println("[honeybadgerscc] Chaincode Started.")

	return &honeybadgerscc{}
}

type honeybadgerscc struct {
	db *leveldb.DB
}

var db *leveldb.DB
var db_name string

var rec_mutex sync.Mutex

type secretcell struct {
	ObjectType string `json:"docType"`
	CellName   string `json:"cellName"`
	IsWritten  bool   `json:"isWriten"`
	WriterKey  string `json:"WriterKey"`
	IsOpen     bool   `json:"IsOpen"`
	Value      string `json:"Value"`
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// Init implements the chaincode shim interface
func (s *honeybadgerscc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("[honeybadgerscc]Init is called")
	/// db_name = GetOutboundIP().String()
	s.db, _ = leveldb.OpenFile("db", nil)
	defer s.db.Close()
	if db == nil {
		fmt.Println("DB is nil")
		// s.Init(stub)
	} else {
		fmt.Println("DB is not nil ")
	}
	return shim.Success(nil)
}

func reconstructHelper(db *leveldb.DB, nodeID string, share string, key string) {
	dbPut(db, key+"_result", "None")
	rec_mutex.Lock()
	res := StartPubRec(nodeID, share)
	rec_mutex.Unlock()
	fmt.Println("In reconstructHelper")
	fmt.Println(res)
	dbPut(db, key+"_result", res)
}

// Invoke implements the chaincode shim interface
func (s *honeybadgerscc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fn, args := stub.GetFunctionAndParameters()
	var err error
	if fn == "hbmpc" && len(args) >= 3 {
		// TODO change the name of this endpoint to HBAVSS
		// args[0] = nodeid
		// args[1] = cellname or key for the key value store
		// args[2] = namespace
		// This starts AVSS on this node
		fmt.Println("[honeybadgerscc]HBAVSS")
		key := args[1]
		namespace := args[2]
		cellJSON, err := stub.GetState(key)
		if cellJSON == nil {
			return shim.Success([]byte("Failed to get cell, cell does not exist"))
		}
		if err != nil {
			return shim.Error(err.Error())
		}
		dbPut(nil, "nodeID", args[0])
		dbPut(nil, "mutex", "1")
		var cellInstance secretcell
		json.Unmarshal([]byte(cellJSON), &cellInstance)

		if !authenticateRequest(namespace, key) {
			return shim.Success([]byte("Failed to authenticate request"))
		}
		share := StartHbmpc(args[0])
		if share != "None" {
			cellInstance.IsWritten = true
			dbPut(nil, key, share) // Store the share in the private db
			cellJSON, err = json.Marshal(cellInstance)
			if err != nil {
				return shim.Error(err.Error())
			}

			err := stub.PutState(key, []byte(cellJSON))
			if err != nil {
				fmt.Errorf("Failed to set asset: %s", args[0])
			}
		} else {
			return shim.Success([]byte("Failed HBMPC-hbavss"))
		}
		return shim.Success([]byte("Ran HBMPC-hbavss"))
	} else if fn == "createCell" && len(args) >= 3 {
		// args[0] = nodeid
		// args[1] = public key of the writer
		// args[2] = namespace
		// Creates secretcell
		fmt.Println("[honeybadgerscc] createCell")
		key := args[0]
		writerKey := args[1]
		objType := "secretcell"
		cell := &secretcell{objType, key, false, writerKey, false, ""}
		cellJSON, err := json.Marshal(cell)
		if err != nil {
			return shim.Error(err.Error())
		}
		stub.PutState(key, cellJSON)
		return shim.Success([]byte("Created Cell"))
	} else if fn == "getCell" && len(args) >= 2 {
		// args[0] = cellname
		// args[1] = namespace
		fmt.Println("[honeybadgerscc] getCell")
		cell, err := stub.GetState(args[0])
		if err == nil {
			return shim.Success([]byte(cell))
		}
	} else if fn == "getResult" && len(args) >= 2 {
		// args[0] = key
		// args[1] = namespace
		key := args[0]
		fmt.Println("[honeybadgerscc] getResult")
		if dbGet(nil, key+"_result") == "None" {
			return shim.Success([]byte("Success"))
		} else {
			val := dbGet(nil, key+"_result")
			var cellInstance secretcell
			cellJSON, _ := stub.GetState(key)
			json.Unmarshal([]byte(cellJSON), &cellInstance)
			cellInstance.Value = val
			cellInstance.IsOpen = true
			cellJSON, err := json.Marshal(cellInstance)
			if err != nil {
				return shim.Error(err.Error())
			}
			stub.PutState(key, cellJSON)
			return shim.Success([]byte(dbGet(nil, key+"_result")))

		}
	} else if fn == "getKey" && len(args) >= 1 {
		// args[0] = key
		// TODO Add authentication to this
		key := args[0]
		share := dbGet(nil, key)
		return shim.Success([]byte(share))
	} else if fn == "pubRecon" && len(args) >= 2 {
		// args[0] = key
		// args[1] = namespace
		fmt.Println("In pubRecon")
		key := args[0]
		nodeID := dbGet(nil, "nodeID")
		cellJSON, err := stub.GetState(key)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println(nodeID)
		share := dbGet(nil, key)
		var cellInstance secretcell
		json.Unmarshal([]byte(cellJSON), &cellInstance)

		go reconstructHelper(nil, nodeID, share, key)
		value := "12"
		// err := stub.PutState(key, []byte(value))
		// if err != nil {
		//	return shim.Error(err.Error())
		// }
		return shim.Success([]byte(value))
	} else if fn == "initRecon" && len(args) >= 2 {
		// Deprecated - only for testing
		// TODO remove

		key := args[0]
		checkBytes, _ := stub.GetState(key)
		if checkBytes == nil {
			return shim.Success([]byte("NA"))
		}
		cellJSON, err := stub.GetState(key)
		if cellJSON == nil {
			return shim.Success([]byte("Failed to get cell, cell does not exist"))
		}
		if err != nil {
			return shim.Error(err.Error())
		}
		var cellInstance secretcell
		json.Unmarshal([]byte(cellJSON), &cellInstance)

		namespace := args[0]
		res := InitiatePubRec(key, namespace)
		stub.PutState(key, []byte(res))
		fmt.Println(res)

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
