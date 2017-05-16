package gramework

import (
	"github.com/valyala/fasthttp"
)

// Proxy request to given url
func (ctx *Context) Proxy(url string) error {
	proxyReq := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(proxyReq)
	proxyReq.SetRequestURI(url)
	return fasthttp.Do(proxyReq, &ctx.Response)
}
