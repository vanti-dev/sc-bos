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

// InterfaceNameForAddr returns the name of the interface that the given IP address is associated with.
func InterfaceNameForAddr(addr netip.Addr) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var lastErr error
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			lastErr = err
			continue
		}
		for _, ifaceAddr := range addrs {
			ifacePrefix, err := netip.ParsePrefix(ifaceAddr.String())
			if err != nil {
				lastErr = err
				continue
			}
			if ifacePrefix.Addr() == addr {
				return iface.Name, nil
			}
		}
	}
	return "", lastErr
}
