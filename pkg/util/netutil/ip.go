package netutil

import (
	"net"
)

// StripPort removes the port from an address string.
func StripPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
