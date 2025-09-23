package netutil

import (
	"net"
	"strings"
)

// StripPort removes the port from an address string.
func StripPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

// MergeHostPort merges multiple host:port strings, where later values override earlier ones.
// The host and port are merged separately, so the result will have the last non-empty host and the last non-empty port.
// If there is no valid hosts or ports, returns an empty string.
// If there is no valid port, returns a bare host without a colon.
// If there is no valid host, returns ":port".
func MergeHostPort(hostPorts ...string) (string, error) {
	var host, port string
	for _, hp := range hostPorts {
		// a bare host is either:
		//  - a hostname (contains no colon)
		//  - an IPv4 address (contains no colon)
		//  - a bracketed IPv6 address (starts with [ and ends with ])
		if !strings.Contains(hp, ":") || (strings.HasPrefix(hp, "[") && strings.HasSuffix(hp, "]")) {
			if hp != "" {
				host = hp
			}
			continue
		}
		h, p, err := net.SplitHostPort(hp)
		if err != nil {
			return "", err
		}

		if h != "" {
			host = h
		}
		if p != "" {
			port = p
		}
	}

	if host == "" && port == "" {
		return "", nil
	} else if port == "" {
		return host, nil
	} else {
		return net.JoinHostPort(host, port), nil
	}
}
