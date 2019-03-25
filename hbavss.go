package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// StartHbmpc start honeybadgerMPC
// nodeID string
// namespace string
// key string
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

// InitiatePubRec initiates public reconstruction
func InitiatePubRec(key string, namespace string) string {
	cmd := exec.Command("python3.7", "-m", "honeybadgermpc.fabric_public_reconstruct", key, namespace)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()
	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
	// fmt.Println(cmd.Stdout)
	results := strings.Split(outb.String(), "\n")

	fmt.Println("TEST" + strings.TrimSpace(results[0]))
	fmt.Println("TEST2" + strings.TrimSpace(results[1]))

	fmt.Println(cmd.Stderr)
	return strings.TrimSpace(results[1])
}

// StartPubRec start public reconstruction
func StartPubRec(nodeID string, share string) string {
	cmd := exec.Command("python3.7", "-m", "honeybadgermpc.public_reconstruct", nodeID, "conf/hbavss.hyper.ini", share)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()
	//dbPut("mutex", "1")
	fmt.Println("In start pub rec ")

	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
	lines := strings.Split(outb.String(), "\n")
	fmt.Println(outb.String())
	fmt.Println(errb.String())
	for _, line := range lines {
		if strings.Contains(line, "Reconstructed Value:") {
			// Very hacky way of doing this.
			secretParts := strings.Split(line, "Reconstructed Value:")
			if len(secretParts) >= 2 {
				secret := secretParts[1]
				fmt.Println("The reconstructed secret is  is ", secret)
				return secret
			}
		}
	}

	return "None"
	// fmt.Println("out : ", outb.String(), "err: ", errb.String())

}

// TODO
// contact namespace associated with chaincode and authorize user
func authenticateRequest(namespace string, key string) bool {
	fmt.Println("[honeybadgerscc] Authenticating request")
	return true
}
