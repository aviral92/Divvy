package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"

	"github.com/google/uuid"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/Divvy/src/pb"
)

// Colon is not a typo!
const (
	discoveryPort    = ":2017"
	controlPort      = ":2018"
	broadcastAddress = "255.255.255.255"
)

// PeerT stores information about other Divvy peers
type PeerT struct {
	ID      uuid.UUID
	Address net.IP
}

// NetworkManger implements the Divvy interface
type NetworkManager struct {
	ID                uuid.UUID
	address           net.IP
	availableToOthers bool
	grpcServer        *grpc.Server

	// List of Divvy peers
	peers []PeerT
}

func NewNetworkManager() *NetworkManager {
	netMgr := &NetworkManager{}
	netMgr.getLocalAddress()

	log.Printf("[Network] IP: %v", netMgr.address)
	netMgr.grpcServer = grpc.NewServer()

	// This seems confusing. Is there a better way of doing this?
	pb.RegisterDivvyServer(netMgr.grpcServer, netMgr)

	return netMgr
}

func (netMgr *NetworkManager) getLocalAddress() (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		goto ERROR
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			goto ERROR
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// Supports only IPv4
			ip = ip.To4()
			if ip == nil {
				continue
			}
			netMgr.address = ip
			netMgr.availableToOthers = true
			return ip, nil
		}
	}

ERROR:
	netMgr.availableToOthers = false
	if err != nil {
		return nil, err
	}
	return nil, errors.New("[Network Manager] Unable to get local address")
}

// Divvy service interface implementation

func (netMgr *NetworkManager) Ping(ctx context.Context, empty *pb.Empty) (*pb.Success, error) {
	log.Printf("[Network] Received ping")
	return &pb.Success{}, nil
}

func (netMgr *NetworkManager) NodeJoin(ctx context.Context, newNode *pb.NewNode) (*pb.Success, error) {
	// A new peer has appeared
	return &pb.Success{}, nil
}

func (netMgr *NetworkManager) AddNewNode(newNode pb.NewNode) {
	// Add the new node to the peers list
	log.Printf("[Network] New Divvy node (ID: %v, IP: %v)", newNode.NodeID, newNode.Address)
	var newPeer PeerT
	var err error
	newPeer.ID, err = uuid.Parse(newNode.NodeID)
	if err != nil {
		log.Printf("[Network] Unable to add new peer: %v", err)
	}

	// Broadcast message is sent to the sender as well. Ignore that message
	if newPeer.ID != netMgr.ID {
		newPeer.Address = net.ParseIP(newNode.Address)
		netMgr.peers = append(netMgr.peers, newPeer)
	}
}

// Discover other Divvy peers on the network
func (netMgr *NetworkManager) DiscoverPeers(nodeID uuid.UUID) int {
	// Send a broadcast message over the LAN
	addr, _ := net.ResolveUDPAddr("udp", broadcastAddress+discoveryPort)
	localAddress, _ := net.ResolveUDPAddr("ucp", "127.0.0.1:0")

	conn, err := net.DialUDP("udp", localAddress, addr)
	if err != nil {
		log.Fatalf("[Network] Unable to dial UDP %v", err)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(&pb.NewNode{NodeID: nodeID.String(), Address: netMgr.address.String()})
	conn.Write(buffer.Bytes())

	return 0
}

func (netMgr *NetworkManager) ListenForDiscoveryMessages() {
	udpData := make([]byte, 2048)
	var newNodeMessage pb.NewNode

	listenAddress, _ := net.ResolveUDPAddr("udp", discoveryPort)
	conn, err := net.ListenUDP("udp", listenAddress)
	defer conn.Close()

	if err != nil {
		log.Fatalf("[Network] Unable to listen for discovery: %v", err)
	}

	// Keep listening for new messages
	for {
		dataLen, _, _ := conn.ReadFromUDP(udpData)
		udpDataBuffer := bytes.NewBuffer(udpData[:dataLen])
		decoder := gob.NewDecoder(udpDataBuffer)
		decoder.Decode(&newNodeMessage)
		go netMgr.AddNewNode(newNodeMessage)
	}
}
