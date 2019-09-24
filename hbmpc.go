package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/syndtr/goleveldb/leveldb"
)

// StartMPC start MPC operation on the peer
// appName is the name of the MPC application e.g cmp, equals etc.
func StartMPC(nodeID string, appName string, shares ...string) string {
	fmt.Println("In StartMPC ")
	fmt.Println("appName is " + appName + " Node ID is " + nodeID)
	params := []string{"-m", "honeybadgermpc.fabric_mpc_runner", nodeID, "conf/hbavss.hyper.ini", appName}
	params = append(params, shares...)
	cmd := exec.Command("python3.7", params...)

	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()

	if errmsg != nil {
		log.Fatalf(errb.String())
		log.Fatalf("[honeybadgerscc]cmd.Run()  failed with %s\n", errmsg)
		log.Fatalf("Another")
		return "None"
	}
	fmt.Println("Finished printing error")
	fmt.Println(outb.String())
	fmt.Println(errb.String())
	fmt.Println("Finished printing stdout and stderr")
	return outb.String()
}

// mpcStarter runs in background the mpc operation and waits until output is available
func mpcStarter(db *leveldb.DB, nodeID string, appName string, instance string, shares ...string) {
	// instance is the instance of the mpc computation
	// shares contains all the shares from the cells
	dbPut(db, instance+"_result", "None")
	rec_mutex.Lock() // TODO change this accomodate multiple MPC ops
	res := StartMPC(nodeID, appName, shares...)
	fmt.Println("Finished StartMPC call")
	rec_mutex.Unlock()
	fmt.Println("Completed MPC operation.")
	dbPut(db, instance+"_result", res)
}
