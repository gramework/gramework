package gramework

import (
	"github.com/valyala/fasthttp"
)

func (ctx *Context) saveCookies() {
	ctx.Cookies.Mu.Lock()
	for k, v := range ctx.Cookies.Storage {
		c := fasthttp.AcquireCookie()
		c.SetKey(k)
		c.SetValue(v)
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
