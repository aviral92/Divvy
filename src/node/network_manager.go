package main

import (
	"errors"
	"net"
)

// NetworkManager is responsible for all network related functions of the node
type NetworkManager struct {
	localAddress          net.IP
	nodeAvailableToOthers bool
}

func (netmgr *NetworkManager) getLocalAddress() (net.IP, error) {
	if netmgr.localAddress != nil {
		return netmgr.localAddress, nil
	}

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
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			netmgr.localAddress = ip
			netmgr.nodeAvailableToOthers = true
			return ip, nil
		}
	}

ERROR:
	netmgr.nodeAvailableToOthers = false
	if err != nil {
		return nil, err
	}
	return nil, errors.New("[Network Manager] Unable to get local address")
}
