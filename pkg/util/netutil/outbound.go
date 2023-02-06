package netutil

import (
	"net"
	"net/netip"
)

// OutboundAddr returns the systems preferred outbound IP address.
// This IP address is likely but not guaranteed to be the IP address that others can use to communicate with this system.
func OutboundAddr() (netip.Addr, error) {
	// inspired by https://stackoverflow.com/a/37382208/317404

	// note: the IP here doesn't actually need to exist and we don't open a connection
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return netip.Addr{}, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip, _ := netip.AddrFromSlice(localAddr.IP)
	return ip, nil
}
