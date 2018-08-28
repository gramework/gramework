package gramework

import (
	"net"
	"sync"
	"sync/atomic"
)

// Protect enables Gramework Protection for routes registered after Protect() call.
//
// Protects all routes, that prefixed with given enpointPrefix.
// For example:
//
//		app := gramework.New()
//		app.GET("/internal/status", serveStatus) // will **not be** protected, .Protected() isn't called yet
//		app.Protect("/internal")
//		registerYourInternalRoutes(app.Sub("/internal")) // all routes here will be protected
//
// Any blacklisted ip can't access protected enpoints via any method.
// Blacklist can work automatically, manually or both. To disable automatic blacklist do App.MaxHackAttemts(-1).
// Automatic blacklist bans suspected IP after App.MaxHackAttempts(). This behaviour is disabled for whitelisted
// ip.
//
// See also App.Whitelist(), App.Untrust(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) Protect(endpointPrefix string) {
	if app.trustedIP == nil {
		app.trustedIP = &ipList{
			list: make(map[string]struct{}),
			mu:   &sync.RWMutex{},
		}
	}
	if app.untrustedIP == nil {
		app.untrustedIP = &ipList{
			list: make(map[string]struct{}),
			mu:   &sync.RWMutex{},
		}
	}
	if app.suspectedIP == nil {
		app.suspectedIP = &suspectsList{
			list: make(map[string]*suspect),
			mu:   &sync.RWMutex{},
		}
	}

	if app.protectedPrefixes == nil {
		app.protectedPrefixes = make(map[string]struct{})
		app.protectedEndpoints = make(map[string]struct{})
	}

	app.protectedPrefixes[endpointPrefix] = struct{}{}
}

func nilHijackHandler(c net.Conn) {
}

func (app *App) protectionMiddleware(handler func(*Context)) func(ctx *Context) {
	return func(ctx *Context) {
		if ctx.IsBlacklisted() {
			// force closing of the connection ASAP
			ctx.Hijack(nilHijackHandler)
			return
		}

		handler(ctx)
	}
}

// Whitelist adds given ip to Gramework Protection trustedIP list.
// To remove IP from whitelist, call App.Untrust()
//
// See also App.Protect(), App.Untrust(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) Whitelist(ip net.IP) (ok bool) {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() {
		return true
	}

	ipHash := app.prepareIPListKey(ip)

	// now we trust this ip
	app.trustedIP.mu.Lock()
	app.trustedIP.list[ipHash] = struct{}{}
	app.trustedIP.mu.Unlock()

	// unban this ip
	app.untrustedIP.mu.Lock()
	delete(app.untrustedIP.list, ipHash)
	app.untrustedIP.mu.Unlock()

	// whitelisted ip can't be suspected
	app.suspectedIP.mu.Lock()
	delete(app.suspectedIP.list, ipHash)
	app.suspectedIP.mu.Unlock()
	return true
}

// Untrust removes given ip from trustedIP list, that
// enables protection of Gramework Protection enabled endpoints for given ip too.
// Opposite of App.Whitelist().
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) Untrust(ip net.IP) (ok bool) {
	if ip == nil {
		return false
	}
	ipHash := app.prepareIPListKey(ip)

	// now we don't trust this ip
	app.trustedIP.mu.Lock()
	delete(app.trustedIP.list, ipHash)
	app.trustedIP.mu.Unlock()
	return true
}

// Blacklist adds given ip to untrustedIP list, if it is not whitelisted. Any ip blacklisted with
// Gramework Protection can't access protected enpoints via any method.
//
// See also App.Protect(), App.Whitelist(), App.Untrust(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) Blacklist(ip net.IP) (ok bool) {
	if ip == nil {
		return false
	}

	ipHash := app.prepareIPListKey(ip)

	app.trustedIP.mu.RLock()
	if _, ok := app.trustedIP.list[ipHash]; ok {
		app.trustedIP.mu.RUnlock()
		return false
	}
	app.trustedIP.mu.RUnlock()

	// ban this ip
	app.untrustedIP.mu.Lock()
	app.untrustedIP.list[ipHash] = struct{}{}
	app.untrustedIP.mu.Unlock()

	// we don't need to suspect already banned ip
	app.suspectedIP.mu.Lock()
	delete(app.suspectedIP.list, ipHash)
	app.suspectedIP.mu.Unlock()
	return true
}

// Suspect adds given ip to Gramework Protection suspectedIP list.
//
// See also App.Protect(), App.Untrust(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) Suspect(ip net.IP) (ok bool) {
	ipHash := app.prepareIPListKey(ip)

	app.trustedIP.mu.RLock()
	if _, ok := app.trustedIP.list[ipHash]; ok {
		app.trustedIP.mu.RUnlock()
		return false
	}
	app.trustedIP.mu.RUnlock()

	// suspect this ip
	app.suspectedIP.mu.Lock()
	app.suspectedIP.list[ipHash] = &suspect{
		hackAttempts: 0,
	}
	app.suspectedIP.mu.Unlock()
	return true
}

// MaxHackAttempts sets new max hack attempts for blacklist triggering in the Gramework Protection.
// If 0 passed, MaxHackAttempts just returns current value
// without setting a new one.
// If -1 passed, automatic blacklist disabled.
// This function is threadsafe and atomic.
//
// See `ctx.Whitelist()`, `ctx.Blacklist()` and `ctx.Suspect()` for manual Gramework Protection control.
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (app *App) MaxHackAttempts(attempts int32) (oldValue int32) {
	oldValue = atomic.LoadInt32(app.maxHackAttempts)
	if attempts != 0 && atomic.CompareAndSwapInt32(app.maxHackAttempts, oldValue, attempts) {
		app.internalLog.
			WithField("old", oldValue).
			WithField("new", attempts).
			Infof("[Gramework Protection] Updated max hack attemts")
	}
	return
}

// IsWhitelisted checks if we have current client in Gramework Protection trustedIP list.
// Use ctx.Whitelist() to add current client to trusted list.
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(),
// App.MaxHackAttempts(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) IsWhitelisted() (isWhitelisted bool) {
	if ctx.RemoteIP().IsLoopback() {
		return true
	}
	ctx.App.trustedIP.mu.RLock()
	_, isWhitelisted = ctx.App.trustedIP.list[ctx.remoteIPHash()]
	ctx.App.trustedIP.mu.RUnlock()
	return
}

// IsBlacklisted checks if we have current client in Gramework Protection untrustedIP list.
// Use ctx.Blacklist() to add current client to untrustedIP list.
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(),
// App.MaxHackAttempts(), Context.IsWhitelisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) IsBlacklisted() (isBlacklisted bool) {
	if ctx.IsWhitelisted() {
		return false
	}
	ctx.App.untrustedIP.mu.RLock()
	_, isBlacklisted = ctx.App.untrustedIP.list[ctx.remoteIPHash()]
	ctx.App.untrustedIP.mu.RUnlock()
	return
}

// IsSuspect checks if we have current client in Gramework Protection suspectedIP list.
// Use ctx.Suspect() to add current client to suspectedIP list.
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(),
// App.MaxHackAttempts(), Context.IsWhitelisted(), Context.IsBlacklisted(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) IsSuspect() (isSuspect bool) {
	if ctx.IsWhitelisted() {
		return false
	}
	ctx.App.suspectedIP.mu.RLock()
	_, isSuspect = ctx.App.suspectedIP.list[ctx.remoteIPHash()]
	ctx.App.suspectedIP.mu.RUnlock()
	return
}

// Whitelist adds given ip to trustedIP list of the Gramework Protection.
// To remove IP from whitelist, call App.Untrust()
//
// See also App.Protect(), App.Untrust(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// App.Whitelist(), Context.Untrust(), Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) Whitelist() (ok bool) {
	return ctx.App.Whitelist(ctx.RemoteIP())
}

// Untrust deletes given ip from trustedIP list, that
// enables protection of Protect()'ed endpoints for given ip too.
// Opposite of Context.Whitelist().
//
// See also App.Protect(), App.Whitelist(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// App.Untrust(), Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) Untrust() (ok bool) {
	return ctx.App.Untrust(ctx.RemoteIP())
}

// Blacklist adds given ip to untrustedIP list, if it is not whitelisted. Any blacklisted ip can't
// access protected enpoints via any method.
//
// See also App.Protect(), App.Whitelist(), App.Untrust(), App.Suspect(), App.MaxHackAttempts(),
// App.Blacklist(), Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) Blacklist() (ok bool) {
	return ctx.App.Blacklist(ctx.RemoteIP())
}

// Suspect adds current client ip to Gramework Protection suspectedIP list.
//
// See also App.Protect(), App.Untrust(), App.Blacklist(), App.Suspect(), App.MaxHackAttempts(),
// Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Blacklist(), Context.Suspect(), Context.HackAttemptDetected(),
// Context.SuspectsHackAttempts()
func (ctx *Context) Suspect() (ok bool) {
	return ctx.App.Suspect(ctx.RemoteIP())
}

// remoteIPHash is a shortcut for app.prepareIPListKey() function that
// calculates hash for the []byte, which is what the ip is
func (ctx *Context) remoteIPHash() string {
	return ctx.App.prepareIPListKey(ctx.RemoteIP())
}

// HackAttemptDetected adds given ip to Gramework Protection suspectedIP list.
// Use it when you detected app-level hack attempt from current client.
//
// See also App.Protect(), App.Whitelist(), App.Untrust(), App.Suspect(), App.MaxHackAttempts(),
// App.Blacklist(), Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Suspect(), Context.Blacklist(),
// Context.SuspectsHackAttempts()
func (ctx *Context) HackAttemptDetected() {
	if ctx.IsWhitelisted() {
		return
	}
	ipHash := ctx.remoteIPHash()

	ctx.App.suspectedIP.mu.Lock()
	if s, ok := ctx.App.suspectedIP.list[ipHash]; !ok || s == nil {
		ctx.App.suspectedIP.list[ipHash] = &suspect{
			hackAttempts: 1,
		}
		ctx.App.suspectedIP.mu.Unlock()
		return
	}
	// release lock ASAP
	ctx.App.suspectedIP.mu.Unlock()
	atomic.AddInt32(&ctx.App.suspectedIP.list[ipHash].hackAttempts, 1)
}

// SuspectsHackAttempts returns hack attempts detected with Gramework Protection
// both automatically and manually by calling Context.HackAttemptDetected().
// For any whitelisted ip this function will return 0.
//
// See also App.Protect(), App.Whitelist(), App.Untrust(), App.Suspect(), App.MaxHackAttempts(),
// App.Blacklist(), Context.IsWhitelisted(), Context.IsBlacklisted(), Context.IsSuspect(),
// Context.Whitelist(), Context.Suspect(), Context.Blacklist(),
// Context.HackAttemptDetected()
func (ctx *Context) SuspectsHackAttempts() (attempts int32) {
	if ctx.IsWhitelisted() {
		return zero
	}
	ctx.App.suspectedIP.mu.RLock()
	if s, ok := ctx.App.suspectedIP.list[ctx.remoteIPHash()]; ok && s != nil {
		attempts = s.hackAttempts
	}
	ctx.App.suspectedIP.mu.RUnlock()
	return
}

func (app *App) prepareIPListKey(ip net.IP) string {
	if ip == nil {
		// we should ban any invalid remoteIP headers
		// to ban this type of attacks
		return ""
	}
	return ip.String()
}
