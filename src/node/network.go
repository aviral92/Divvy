package main

import (
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
	controlPort = ":2018"
)

// PeerT stores information about other Divvy peers
type PeerT struct {
	ID      uuid.UUID
	Address net.IP
}

// NetworkManger implements the Divvy interface
type NetworkManager struct {
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
