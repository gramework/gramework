// Copyright 2015-2017 Metthew Holt
// Copyright 2017 Kirill Danshin

package gramework

import (
	"bytes"
	"net"
	"net/url"
	"strings"
)

const (
	localhost            = "localhost"
	brackets             = "[]"
	localIPV6            = "::1"
	localIPV6Brackets    = "[::1]"
	localIPV4SubnetShort = "127."
	col                  = ":"
)

var (
	localhostColB         = []byte("localhost:")
	localIPV4SubnetShortB = []byte(localIPV4SubnetShort)
)

// IsLoopback returns true if the hostname of addr looks
// explicitly like a common local hostname. addr must only
// be a host or a host:port combination.
func IsLoopback(addr string) bool {
	if addr == localhost || (len(addr) >= 10 && bytes.Equal([]byte(addr[0:9]), localhostColB)) ||
		addr == localIPV6Brackets {
		return true
	}
	var host = addr
	if strings.Contains(addr, col) {
		h, _, err := net.SplitHostPort(addr)
		if err == nil {
			host = h // happens if the addr is not just a hostname
		}
	}
	if isLoopbackHost(host) {
		return true
	}
	u, err := url.Parse(addr)
	if err != nil {
		return false
	}
	return isLoopbackHost(u.Host)
}

func isLoopbackHost(host string) bool {
	if host == localhost || strings.Trim(host, brackets) == localIPV6 { // skip net.ParseIP if possible
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && (len(host) > 4 && bytes.Equal([]byte(host[0:4]), localIPV4SubnetShortB))
}
