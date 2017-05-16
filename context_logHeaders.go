package gramework

func (ctx *Context) LogHeaders() {
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		ctx.Logger.Debugf("%s = [%s]\n", k, v)
	})
}
