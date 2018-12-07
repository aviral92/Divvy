package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
    "strings"

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

func (netMgr *NetworkManager) GetFileList(ctx context.Context, empty *pb.Empty) (*pb.FileList, error) {
    // TODO: Call the File manager to get all files
    return &pb.FileList{}, nil
}

func (netMgr *NetworkManager) Search(ctx context.Context, query *pb.SearchQuery) (*pb.FileList, error) {
    // TODO: Call the File manager to get all the files matching name/hash
    return &pb.FileList{}, nil
}

func (netMgr *NetworkManager) DownloadFile(ctx context.Context, request *pb.DownloadRequest) (*pb.Success, error) {
    // TODO: The request should be forwarded to the download manager
    return &pb.Success{}, nil
}

/*
*  Discover other Divvy peers on the network
*/

func (netMgr *NetworkManager) AddNewNode(newNode pb.NewNode) {
	// Add the new node to the peers list
	log.Printf("[Network] New Divvy node (ID: %v, IP: %v)", newNode.NodeID, newNode.Address)
	var newPeer PeerT
	var err error
	newPeer.ID, err = uuid.Parse(newNode.NodeID)
	if err != nil {
		log.Printf("[Network] Unable to add new peer: %v", err)
	}

    newPeer.Address = net.ParseIP(newNode.Address)
    netMgr.peers = append(netMgr.peers, newPeer)
}

func (netMgr *NetworkManager) DiscoverPeers() int {
	// Send a broadcast message over the LAN
	addr, _ := net.ResolveUDPAddr("udp", broadcastAddress+discoveryPort)
	localAddress, _ := net.ResolveUDPAddr("ucp", "127.0.0.1:0")

	conn, err := net.DialUDP("udp", localAddress, addr)
	if err != nil {
		log.Fatalf("[Network] Unable to dial UDP %v", err)
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
    encoder.Encode(&pb.NewNode{NodeID: netMgr.ID.String(), Address: netMgr.address.String(), IsReply: false})
	conn.Write(buffer.Bytes())

	return 0
}

func (netMgr *NetworkManager) ListenForDiscoveryMessages() {
	udpData := make([]byte, 2048)
	var newNodeMessage pb.NewNode
    var buffer bytes.Buffer

	listenAddress, _ := net.ResolveUDPAddr("udp", discoveryPort)
	conn, err := net.ListenUDP("udp", listenAddress)
	defer conn.Close()

	if err != nil {
		log.Fatalf("[Network] Unable to listen for discovery: %v", err)
	}

	// Keep listening for new messages
	for {
		dataLen, peerAddr, _ := conn.ReadFromUDP(udpData)
		udpDataBuffer := bytes.NewBuffer(udpData[:dataLen])
		decoder := gob.NewDecoder(udpDataBuffer)
		decoder.Decode(&newNodeMessage)

        if newNodeMessage.NodeID == netMgr.ID.String() {
            continue
        }

		go netMgr.AddNewNode(newNodeMessage)

        // Respond to the peer if it is not a reply
        if newNodeMessage.IsReply == true {
            continue
        }
        peerIP := strings.Split(peerAddr.String(), ":")[0]
        log.Printf("New Peer IP: %v", peerIP)
        newPeerAddr, _ := net.ResolveUDPAddr("udp", peerIP+discoveryPort)
        encoder := gob.NewEncoder(&buffer)
        encoder.Encode(&pb.NewNode{NodeID: netMgr.ID.String(), Address: netMgr.address.String(), IsReply: true})
        conn.WriteToUDP(buffer.Bytes(), newPeerAddr)
	}
}
