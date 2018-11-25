/*
* Implementation of a Divvy node
* This is where all sub-services are integrated with each other
 */

package main

import (
	"github.com/google/uuid"
	"log"
)

// PeerT stores information about other Divvy peers
type PeerT struct {
	ID      uuid.UUID
	Address string
}

// NodeT stores information about this node
type NodeT struct {
	ID      uuid.UUID
	Address string
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
func initNode(Node* NodeT) {
}

// Main function that handles all requests from sub-services
func main() {
	// Initialize node
	Node := newNodeT()

    // Initialize network manager
    // Read configuration
	// Initialize file manager

    initNode(&Node)

	log.Printf("Node ID: %v", Node.ID.String())

	// for loop that keeps listening for events
}
