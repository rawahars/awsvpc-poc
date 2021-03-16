package main

import "net"

func GetIPAddressFromString(ipAddress string) (*net.IPNet, error) {
	address, prefix, err := net.ParseCIDR(ipAddress)
	if err != nil {
		return nil, err
	}

	prefix.IP = address

	return prefix, nil
}

func GetSubnetPrefix(ipAddress *net.IPNet) *net.IPNet {
	return &net.IPNet{
		IP:   ipAddress.IP.Mask(ipAddress.Mask),
		Mask: ipAddress.Mask,
	}
}