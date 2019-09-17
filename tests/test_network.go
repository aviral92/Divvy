/*
*  Testcases related to the network manager
 */

package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/Divvy/src/pb"
)

// Colon is not a typo!
const (
	broadcastAddress = "255.255.255.255"
	discoveryPort    = ":2017"
)

// Send a broadcast message not to the local IP
func TestDiscoveryListener() {
	addr, _ := net.ResolveUDPAddr("udp", broadcastAddress+discoveryPort)
	localAddress, _ := net.ResolveUDPAddr("ucp", "127.0.0.1:0")

	conn, err := net.DialUDP("udp", localAddress, addr)
	if err != nil {
		log.Fatalf("[Test][Network] TestDiscoveryListener failed %v", err)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(&pb.NewNode{NodeID: "5f5c83fe-fb22-4896-bbf5-1f62211425b6", Address: "0.0.0.0"})
	conn.Write(buffer.Bytes())
}

func main() {
	TestDiscoveryListener()
}
