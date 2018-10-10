// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"time"

	"github.com/valyala/fasthttp"
)

const defaultCookiePath = "/"

// GetCookieDomain returns previously configured cookie domain and if cookie domain
// was configured at all
func (ctx *Context) GetCookieDomain() (domain string, wasConfigured bool) {
	return ctx.App.cookieDomain, len(ctx.App.cookieDomain) > 0
}

func (ctx *Context) saveCookies() {
	ctx.Cookies.Mu.Lock()
	for k, v := range ctx.Cookies.Storage {
		c := fasthttp.AcquireCookie()
		c.SetKey(k)
		c.SetValue(v)
		if len(ctx.App.cookieDomain) > 0 {
			c.SetDomain(ctx.App.cookieDomain)
		}
		if len(ctx.App.cookieDomain) > 0 {
			c.SetPath(ctx.App.cookiePath)
		} else {
			c.SetPath(defaultCookiePath)
		}

		c.SetExpire(time.Now().Add(ctx.App.cookieExpire))
		ctx.Response.Header.SetCookie(c)
		fasthttp.ReleaseCookie(c)
	}
	ctx.Cookies.Mu.Unlock()
}

func (ctx *Context) loadCookies() {
	ctx.Cookies.Storage = make(map[string]string, zero)
	ctx.Request.Header.VisitAllCookie(ctx.loadCookieVisitor)
}

func (ctx *Context) loadCookieVisitor(k, v []byte) {
	ctx.Cookies.Set(string(k), string(v))
}

// Set a cookie with given key to the value
func (c *Cookies) Set(key, value string) {
	c.Mu.Lock()
	if c.Storage == nil {
		c.Storage = make(map[string]string, zero)
	}
	c.Storage[key] = value
	c.Mu.Unlock()
}

// Get a cookie by given key
func (c *Cookies) Get(key string) (string, bool) {
	c.Mu.Lock()
	if c.Storage == nil {
		c.Storage = make(map[string]string, zero)
		c.Mu.Unlock()
		return emptyString, false
	}
	if v, ok := c.Storage[key]; ok {
		c.Mu.Unlock()
		return v, ok
	}
	c.Mu.Unlock()
	return emptyString, false
}

// Exists reports if the given key exists for current request
func (c *Cookies) Exists(key string) bool {
	c.Mu.Lock()
	if c.Storage == nil {
		c.Storage = make(map[string]string, zero)
		c.Mu.Unlock()
		return false
	}
	if _, ok := c.Storage[key]; ok {
		c.Mu.Unlock()
		return ok
	}
	c.Mu.Unlock()
	return false
}
