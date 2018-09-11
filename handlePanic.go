package gramework

import "strings"

var errGotPanic = struct {
	Code    int    `json:"code" xml:"code" csv:"code"`
	Message string `json:"message" xml:"message" csv:"message"`
}{
	Code:    500,
	Message: "Internal Server Error",
}

// DefaultPanicHandler serves error page or error response depending on ctx.ContentType()
func DefaultPanicHandler(ctx *Context, panicReason interface{}) {
	ctx.SetStatusCode(500)
	if strings.HasPrefix(string(ctx.Request.Header.Peek("Accept")), "text/html") || strings.Contains(ctx.ContentType(), "text/html") {
		_, err := ctx.HTML().WriteString(handledPanic)
		if err != nil {
			ctx.Logger.WithError(err).Error("could not serve default panic page")
			return
		}
		if !ctx.App.PanicHandlerNoPoweredBy {
			ctx.WriteString(poweredBy)
		}
		if len(ctx.App.PanicHandlerCustomLayout) > 0 {
			ctx.WriteString(ctx.App.PanicHandlerCustomLayout)
		}
		return
	}

	ctx.Encode(errGotPanic)
}

const handledPanic = `<!doctype html>
<html>
<meta charset=utf-8>
<meta name=viewport content="width=device-width,initial-scale=1">
<title>Internal Server Error</title>
<meta http-equiv="Content-Security-Policy" content="default-src 'none'; base-uri 'self'; connect-src 'self'; form-action 'self'; img-src 'self' data:; script-src 'self'; style-src 'unsafe-inline'">
<style>
html {
	position: relative;
	font-family: sans-serif;
	-webkit-font-smoothing: antialiased;
	text-rendering: optimizeLegibility;
}
body, html,
.mainWrapper {
	height: 100%;
	width: 100%;
	margin: 0;
	padding: 0;
}
.mainWrapper {
	background: #000;
	background-image:
		linear-gradient(
			0deg,
			rgba(0,0,0,0) 24%,
			rgba(255,255,255,.05) 25%,
			rgba(255,255,255,.05) 26%,
			rgba(0,0,0,0) 27%,
			rgba(0,0,0,0) 74%,
			rgba(255,255,255,.05) 75%,
			rgba(255,255,255,.05) 76%,
			rgba(0,0,0,0) 77%,
			rgba(0,0,0,0)
		),
		linear-gradient(
			90deg,
			rgba(0,0,0,0) 24%,
			rgba(255,255,255,.05) 25%,
			rgba(255,255,255,.05) 26%,
			rgba(0,0,0,0) 27%,
			rgba(0,0,0,0) 74%,
			rgba(255,255,255,.05) 75%,
			rgba(255,255,255,.05) 76%,
			rgba(0,0,0,0) 77%,
			rgba(0,0,0,0)
		);
	background-size: 25px 25px;
	background-position: left 30px top 31px;
	display: flex;
	align-items: center;
	justify-content: center;
	flex-direction: column;
}
h1 {
	color: #1e2c7f;
	color: #fafafa;
	font-size: 72px;
	margin-top: 0;
}
p {
	margin: 0;
}
div {
	color: #fff;
}
.poweredBy {
	position: absolute;
	opacity: .65;
	width: 100%;
	bottom: 0;
	text-align: center;
	font-size: 10px;
}
a {
	font-weight: bold;
	color: #fff !important;
}
.poweredBy p {
	padding: 10px;
}
</style>
<div class="mainWrapper">
	<h1>500</h1>
	<p>Sorry, our service is currently unavailable.</p>
	<p>Wait a minute and try again.</p>
</div>
`

const poweredBy = `<div class="poweredBy">
	<p>Powered by <a target=_blank href="https://github.com/gramework/gramework">Gramework</a>.</p>
</div>`
