package gramework

// LogHeaders for debug
func (ctx *Context) LogHeaders() {
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		ctx.Logger.Debugf("%s = [%s]\n", k, v)
	})
}
