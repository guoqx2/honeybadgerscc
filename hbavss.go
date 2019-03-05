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
