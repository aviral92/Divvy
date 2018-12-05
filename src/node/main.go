/*
 * Implementation of a Divvy node
 * This is where all sub-services are integrated with each other
 */

package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	//"google.golang.org/grpc"
	//"github.com/Divvy/src/pb"
)

// PeerT stores information about other Divvy peers
type PeerT struct {
	ID      uuid.UUID
	Address net.IP
}

// NodeT stores information about this node
type NodeT struct {
	ID     uuid.UUID
	netMgr *NetworkManager

	// List of Divvy peers
	peers []PeerT

	// Create file Manager
	fileMgr *FileManager
}

// Initialize a new NodeT
func newNodeT() NodeT {
	Node := NodeT{}
	Node.ID = uuid.New()
	// Make an empty list of peers
	Node.peers = make([]PeerT, 0)
	return Node
}

// Initialize everything about this node
func initNode(Node *NodeT) {
	log.Printf("[Node] Initializing Divvy node...")
	log.Printf("[Node] ID: %v", Node.ID.String())

	Node.netMgr = NewNetworkManager()
	Node.fileMgr = NewFileManager("/home/vagrant/go/src/github.com/Divvy/test")

	// Read configuration
	// Initialize file manager

	log.Printf("[Node] Divvy node initialized!")
}

// Main function that handles all requests from sub-services
func main() {
	fmt.Println("jnkjnjn")
	// Initialize node
	Node := newNodeT()

	initNode(&Node)

	log.Printf("path is %s", Node.fileMgr.files[0].Path)
	log.Printf("file exists is %v", Node.fileMgr.files[0].exists(Node.fileMgr.files[0].Path))
	log.Printf("hash is %s", Node.fileMgr.files[0].GetHash(Node.fileMgr.files[0].Path))
	Node.fileMgr.displayDirectory()
	/*
	   Once everything is setup start listening. This call is blocking
	   Do not put any logic after gRPC serve
	*/
	conn, err := net.Listen("tcp", controlPort)
	if err != nil {
		log.Fatalf("[Node] Failed to open port %v because %v", controlPort, err)
	}
	log.Printf("[Node] Listening on port %v", controlPort)
	err = Node.netMgr.grpcServer.Serve(conn)
	if err != nil {
		log.Fatalf("[Node] Failed to serve %v", err)
	}

	log.Printf("[Node] Bye from Divvy!")
}
