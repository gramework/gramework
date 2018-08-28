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
	"github.com/valyala/fasthttp"
)

// Proxy request to given url
func (ctx *Context) Proxy(url string) error {
	proxyReq := fasthttp.AcquireRequest()
	ctx.Request.CopyTo(proxyReq)
	proxyReq.SetRequestURI(url)
	return fasthttp.Do(proxyReq, &ctx.Response)
}
