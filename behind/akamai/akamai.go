// Package akamai provides a gramework.Behind implementation
// developed for Gramework.
// This is not an official Akamai-supported implementation.
// If you having any issues with this package, please
// consider to contact Gramework support first.
// Akamai doesn't provide any official support nor guaranties
// about this package.
//
// Akamai is a trademark of Akamai Technologies, Inc.
//
// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
package akamai

import (
	"bytes"
	"encoding/csv"
	"net"
	"net/url"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/gramework/gramework"
	"github.com/gramework/gramework/behind/akamai/internal/portWhitelist"
)

const (
	DefaultAkamaiIPHeader = "True-Client-IP"
)

type Option func(a *Unwrapper)

type Unwrapper struct {
	ipHeader        string
	whitelistedCIDR []*net.IPNet

	initCache    sync.Once
	disableCache bool
	cache        *fastcache.Cache
}

func (a *Unwrapper) OnAppActivation(app *gramework.App) {
	log := app.Logger.WithField("package", "gramework/behind/akamai")
	log.Info("Gramework is running behind Akamai.")
	log = log.WithField("CIDRs", len(a.whitelistedCIDR))
	if a.disableCache {
		log.WithField("cache", "disabled")
	} else {
		log.WithField("cache", "enabled")
	}
	log.Info("Activated")
}

// New creates an unwrapper, optimized for Akamai network
func New(opts ...Option) *Unwrapper {
	a := &Unwrapper{
		ipHeader: DefaultAkamaiIPHeader,
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// DisableCache sets
func DisableCache() Option {
	return func(a *Unwrapper) {
		a.disableCache = true
	}
}

// Service Name,CIDR Block,Port,Activation Date,CIDR Status
// [0]          [1]        [2]  [3]             [4]
type akamaiCSVRecord []string

func (rec akamaiCSVRecord) getCIDR() string {
	_ = rec[4]

	return rec[1]
}

func (rec akamaiCSVRecord) getPort() string {
	_ = rec[4]

	return rec[2]
}

func ParseCIDRBlocksCSV(csvDump []byte, usePortFilter, strict bool) (cidrBlocks []*net.IPNet, err error) {
	csvr := csv.NewReader(bytes.NewReader(csvDump))

	records, err := csvr.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		cidr := akamaiCSVRecord.getCIDR(record)
		port := akamaiCSVRecord.getPort(record)
		needed := !usePortFilter || portWhitelist.IsPortInRange(port)
		if needed {
			_, parsedCIDR, err := net.ParseCIDR(cidr)
			if strict && err != nil {
				return nil, err
			}
			if err == nil {
				cidrBlocks = append(cidrBlocks, parsedCIDR)
			}
		}
	}

	return cidrBlocks, nil
}

func CIDRBlocks(blocks []*net.IPNet) Option {
	return func(a *Unwrapper) {
		a.whitelistedCIDR = blocks
	}
}

func IPHeader(name string) func(*Unwrapper) {
	return func(a *Unwrapper) {
		a.ipHeader = name
	}
}

type Addr struct {
	network string // name of the network (for example, "tcp", "udp")
	remote  string // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

func (addr *Addr) Network() string {
	return addr.network
}

func (addr *Addr) String() string {
	return addr.remote
}

func (addr *Addr) setRemote(actualIP, origRemote string) error {
	u, err := url.Parse(origRemote)
	if err != nil {
		return err
	}

	u.Host = actualIP
	addr.remote = u.String()
	return nil
}

func (a *Unwrapper) RemoteAddr(ctx *gramework.Context) net.Addr {
	rIP := ctx.RequestCtx.RemoteIP()
	orig := ctx.RequestCtx.RemoteAddr()
	if !a.lookupWhitelisted(rIP) {
		return orig
	}

	actualIP := a.RemoteIP(ctx)
	if actualIP == nil {
		return orig
	}
	addr := &Addr{
		network: orig.Network(),
	}

	err := addr.setRemote(actualIP.String(), orig.String())
	if err != nil {
		return orig
	}
	return addr
}

func (a *Unwrapper) RemoteIP(ctx *gramework.Context) net.IP {
	ip := net.ParseIP(string(ctx.Request.Header.Peek(a.ipHeader)))
	rIP := ctx.RequestCtx.RemoteIP()
	if len(ip) == 0 || !a.lookupWhitelisted(rIP) {
		return rIP
	}

	return ip
}

func (a *Unwrapper) initCacheFunc() {
	if a.disableCache {
		return
	}

	a.cache = fastcache.New(32 * 1024 * 1024)
}

func (a *Unwrapper) cacheWrite(key net.IP, value []byte) {
	a.initCache.Do(a.initCacheFunc)
	if a.cache != nil {
		a.cache.Set(key, value)
	}
}

func (a *Unwrapper) cacheRead(key net.IP) (allowed, ok bool) {
	if a.cache != nil {
		result := a.cache.Get(nil, key)
		nonnil := result != nil
		return nonnil && result[0] == 1, nonnil
	}

	return false, false
}

func (a *Unwrapper) lookupWhitelisted(rIP net.IP) (allowed bool) {
	allowed, found := a.cacheRead(rIP)
	if found {
		return allowed
	}

	for _, cidr := range a.whitelistedCIDR {
		if cidr.Contains(rIP) {
			a.cacheWrite(rIP, []byte{1})
			return true
		}
	}

	return false
}
