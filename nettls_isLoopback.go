// Copyright Metthew Holt 2015-2017

package gramework

import (
	"net"
	"strings"
)

// IsLoopback returns true if the hostname of addr looks
// explicitly like a common local hostname. addr must only
// be a host or a host:port combination.
func IsLoopback(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr // happens if the addr is just a hostname
	}
	return host == "localhost" ||
		strings.Trim(host, "[]") == "::1" ||
		strings.HasPrefix(host, "127.")
}
