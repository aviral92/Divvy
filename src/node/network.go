package main

import (
	"errors"
	"net"
    "log"

    context "golang.org/x/net/context"
    "google.golang.org/grpc"

    "github.com/Divvy/src/pb"
)

// Colon is not a typo!
const (
    controlPort = ":2018"
)

// NetworkManger implements the Divvy interface
type NetworkManager struct {
    address             net.IP
    availableToOthers   bool
    grpcServer          *grpc.Server
}

func NewNetworkManager() (*NetworkManager) {
    netMgr := &NetworkManager{}
    netMgr.getLocalAddress()

	log.Printf("[Network] IP: %v", netMgr.address)
    netMgr.grpcServer = grpc.NewServer()

    // This seems confusing. Is there a better way of doing this?
    pb.RegisterDivvyServer(netMgr.grpcServer, netMgr)

    return netMgr
}

func (netmgr *NetworkManager) Ping(ctx context.Context, empty *pb.Empty) (*pb.Success, error) {
    log.Printf("[Network] Received ping")
    return &pb.Success{}, nil
}

func (netMgr *NetworkManager) getLocalAddress() (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		goto ERROR
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		// handle err
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
            // Supports on IPv4
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
