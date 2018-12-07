/*
 * Implementation of a Divvy node
 * This is where all sub-services are integrated with each other
 */

package main

import (
	//"fmt"
	"log"
	"net"

	//"google.golang.org/grpc"

	//"github.com/Divvy/src/node"
	"github.com/google/uuid"
)

// NodeT stores information about this node
type NodeT struct {
	ID      uuid.UUID
	netMgr  *NetworkManager
	fileMgr *FileManager

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

// Initialize everything about this node
func initNode(Node *NodeT) {
	log.Printf("[Node] Initializing Divvy node...")
	log.Printf("[Node] ID: %v", Node.ID.String())

	Node.netMgr = NewNetworkManager()
	// Redundant but saves computation
	Node.netMgr.ID = Node.ID

	Node.fileMgr = NewFileManager("/home/vagrant/go/src/github.com/Divvy/test")

	log.Printf("[Node] Divvy node initialized!")
}

// Main function that handles all requests from sub-services
func main() {

	// Initialize node
	Node := newNodeT()

	initNode(&Node)

	// Discovery listener. Do this before sending the discovery messages
	go Node.netMgr.ListenForDiscoveryMessages()

	Node.netMgr.DiscoverPeers()

	//go Run()
	// Once everything is setup start listening. This call is blocking
	// Do not put any logic after gRPC serve

	// gRPC server
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
