// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

// LogHeaders logs all request headers for debug
func (ctx *Context) LogHeaders() {
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		ctx.Logger.Debugf("%s = [%s]\n", k, v)
	})
}
