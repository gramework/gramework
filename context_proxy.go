package gramework

import (
	"github.com/valyala/fasthttp"
)

func (ctx *Context) Proxy(url string) {
	fasthttp.Do(&ctx.Request, &ctx.Response)
}
