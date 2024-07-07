package cidr

import (
	"fmt"
	"net"
)

// CIDR Struct describing Classless Inter-Domain Routing
type CIDR struct {
	ipNet *net.IPNet
	ip    net.IP
}

// NewCIDR Creates CIDR object
func NewCIDR(subnet string) (*CIDR, error) {
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("error while parsing subnet: %w", err)
	}

	return &CIDR{
		ip:    ip,
		ipNet: ipNet,
	}, nil
}

// Contains Returns true if ip is in CIDR subnet
func (c CIDR) Contains(ip string) bool {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return false
	}

	return c.ipNet.Contains(netIP)
}
