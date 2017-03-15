// Copyright 2015-2017 Metthew Holt
// Copyright 2017 Kirill Danshin

package gramework

import (
	"net"
	"net/url"
	"strings"
)

// IsLoopback returns true if the hostname of addr looks
// explicitly like a common local hostname. addr must only
// be a host or a host:port combination.
func IsLoopback(addr string) bool {
	ip := net.ParseIP(addr)
	if _, err := url.Parse(addr); err != nil && ip == nil {
		return false
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr // happens if the addr is just a hostname
	}
	return host == "localhost" ||
		strings.Trim(host, "[]") == "::1" ||
		(ip != nil && strings.HasPrefix(host, "127."))
}
