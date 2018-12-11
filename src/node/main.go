/*
 * Implementation of a Divvy node
 * This is where all sub-services are integrated with each other
 */

package main

import (
	//"fmt"
	"log"
	//"net"

	//"google.golang.org/grpc"

	//"github.com/Divvy/src/node"
	"github.com/google/uuid"
)

/*
*  Data structure that holds all the information. There should be only one NodeT
*  instantiations throughout the program. All services should use only this
*  object
 */

var Node NodeT

// NodeT stores information about this node
type NodeT struct {
	ID      uuid.UUID
	netMgr  *NetworkManager
	fileMgr *FileManager
	config  *Configuration

	// List of Divvy peers
	peers []PeerT
}

// Initialize a new NodeT
func newNodeT() NodeT {
	Node := NodeT{}
	Node.ID = uuid.New()

	// Make an empty list of peers
	Node.peers = make([]PeerT, 0)

	return Node
}

func getNode() *NodeT {
	return &Node
}

// Initialize everything about this node
func initNode(Node *NodeT) {
	log.Printf("[Node] Initializing Divvy node...")
	log.Printf("[Node] ID: %v", Node.ID.String())

	// Read configuration file
	Node.config = ReadConfigFile("config.json")
	log.Printf("[Node] Network interface: %v", Node.config.NetworkInterface)

	Node.netMgr = NewNetworkManager()
	// Redundant but saves computation
	Node.netMgr.ID = Node.ID

	//Create file manager and pass the path to shared directory
	Node.fileMgr = NewFileManager(Node.config.SharedDirectory)

	log.Printf("[Node] Divvy node initialized!")
}

// Main function that handles all requests from sub-services
func main() {
	Run()
}
