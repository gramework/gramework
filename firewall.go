package gramework

import (
	"sync/atomic"
	"time"
)

func (app *App) initFirewall() {
	go app.firewall.cleanupRequestsCountLoop()
	go app.firewall.releaseBlockedRemoteAddrLoop()
}

// NewRequest tells the firewall, that a new request happened.
// True is returned, if this request should be blocked,
// because the IP is on the block list.
// The remote address is always returned for logging purpose.
func (fw *firewall) NewRequest(ctx *Context) (shouldBeBlocked bool, remoteAddr string) {
	if ctx == nil || ctx.RemoteAddr().String() == "0.0.0.0" {
		return false, ""
	}
	// Get the remote addresse of the request
	remoteAddr = ctx.RemoteIP().String()

	// Check if this remote address is blocked
	if fw.isBlocked(remoteAddr) {
		return true, remoteAddr
	}

	// Register the new request in a new goroutine
	go fw.addRequest(remoteAddr)

	return false, remoteAddr
}

// isBlocked checks if the remote address is blocked
func (fw *firewall) isBlocked(remoteAddr string) bool {
	fw.blockListMutex.Lock()

	// Check if the remote address exists in the blocked map
	_, exists := fw.blockList[remoteAddr]

	fw.blockListMutex.Unlock()
	return exists
}

func (fw *firewall) addRequest(remoteAddr string) {
	fw.requestCounterMutex.Lock()
	defer fw.requestCounterMutex.Unlock()

	count, ok := fw.requestCounter[remoteAddr]
	if !ok {
		count = 1
	}

	// Add the remote address to the block map, if the count
	// reached the limit.
	if count > atomic.LoadInt64(fw.MaxReqPerMin) {
		// Remove the remote address from the request counter map
		delete(fw.requestCounter, remoteAddr)

		// Get the current timestamp
		timestamp := time.Now().Unix()

		// Lock the mutex
		fw.blockListMutex.Lock()
		defer fw.blockListMutex.Unlock()

		// Add the remote address with the timestamp to the block map
		fw.blockList[remoteAddr] = timestamp

		return
	}

	// Save the incremented count to the map
	fw.requestCounter[remoteAddr] = count + 1
}

func (fw *firewall) cleanupRequestsCountLoop() {
	for {
		time.Sleep(time.Minute)

		fw.requestCounterMutex.Lock()

		// Clear the map if not empty
		if len(fw.requestCounter) > 0 {
			fw.requestCounter = make(map[string]int64)
		}
		fw.requestCounterMutex.Unlock()
	}
}

func (fw *firewall) releaseBlockedRemoteAddrLoop() {
	for {
		time.Sleep(time.Minute)

		fw.blockListMutex.Lock()

		if len(fw.blockList) == 0 {
			// nothing to do here
			fw.blockListMutex.Unlock()
			return
		}

		releaseTimestamp := time.Now().Unix() - int64(atomic.LoadInt64(fw.BlockTimeout))

		// removing expired blocks
		for key, timestamp := range fw.blockList {
			if timestamp < releaseTimestamp {
				delete(fw.blockList, key)
			}
		}
		fw.blockListMutex.Unlock()
	}
}
