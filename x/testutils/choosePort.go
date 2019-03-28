package testutils

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

var portsRegister = map[uint16]struct{}{}
var portsRegisterMu = &sync.Mutex{}

// PortChooser does not store any port information.
// It was created only for API purposes.
type PortChooser struct {
	nonRoot *bool
	unused  *bool
}

func Port() *PortChooser {
	return &PortChooser{}
}

func (pc *PortChooser) NonRoot() *PortChooser {
	pc.nonRoot = new(bool)
	*pc.nonRoot = true
	return pc
}

func (pc *PortChooser) Unused() *PortChooser {
	pc.unused = new(bool)
	*pc.unused = true
	return pc
}

func (pc *PortChooser) Root() *PortChooser {
	pc.nonRoot = new(bool)
	return pc
}

func (pc *PortChooser) Used() *PortChooser {
	pc.nonRoot = new(bool)
	return pc
}

func (pc *PortChooser) Acquire() int {
	port := pc.determinePort()

	portsRegisterMu.Lock()
	portsRegister[uint16(port)] = struct{}{}
	portsRegisterMu.Unlock()
	return port
}

func (pc *PortChooser) determinePort() int {
	minPort := 1
	maxPort := 65535

	if pc.nonRoot != nil {
		if *pc.nonRoot {
			minPort = 1025
		} else {
			maxPort = 1024
		}
	}

	choosenPort := 0
	if pc.unused != nil && !*pc.unused {
		portsRegisterMu.Lock()
		for port := range portsRegister {
			if int(port) > minPort && int(port) < maxPort {
				choosenPort = int(port)
			}
		}
		portsRegisterMu.Unlock()
		if choosenPort != 0 {
			return choosenPort
		}

		choosenPort = rand.Intn(maxPort-minPort) + minPort
		_, err := net.Listen("tcp4", fmt.Sprintf(":%d", choosenPort))
		_ = err // fixes linter warning
	} else {
		for {
			choosenPort = rand.Intn(maxPort-minPort) + minPort
			ln, err := net.Listen("tcp4", fmt.Sprintf(":%d", choosenPort))
			if err == nil {
				ln.Close()
				time.Sleep(200 * time.Millisecond)
				break
			}
		}
	}
	return choosenPort
}
