package fasthttprouter

import (
	"github.com/gramework/gramework"
	"github.com/valyala/fasthttp"
)

func (r *Router) handle(path, method string, ctx *gramework.Context, handler func(ctx *gramework.Context), redirectTrailingSlashs bool, isRootRouter bool) (handlerFound bool) {
	if r.router.PanicHandler != nil {
		defer r.router.Recv(ctx)
	}
	if root := r.router.Trees[method]; root != nil {
		if f, tsr := root.GetValue(path, ctx, string(ctx.Method())); f != nil {
			f(ctx)
			return true
		} else if method != CONNECT && path != PathSlash {
			code := redirectCode // Permanent redirect, request with GET method
			if method != GET {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = temporaryRedirectCode
			}

			if tsr && r.router.RedirectTrailingSlash {
				var uri string
				if len(path) > one && path[len(path)-one] == SlashByte {
					uri = path[:len(path)-one]
				} else {
					uri = path + PathSlash
				}
				if uri != emptyString {
					ctx.SetStatusCode(code)
					ctx.Response.Header.Add("Location", uri)
				}
				return false
			}

			// Try to fix the request path
			if r.router.RedirectFixedPath {
				fixedPath, found := root.FindCaseInsensitivePath(
					CleanPath(path),
					r.router.RedirectTrailingSlash,
				)

				if found && len(fixedPath) > 0 {
					queryBuf := ctx.URI().QueryString()
					if len(queryBuf) > zero {
						fixedPath = append(fixedPath, QuestionMark...)
						fixedPath = append(fixedPath, queryBuf...)
					}
					uri := string(fixedPath)
					ctx.SetStatusCode(code)
					ctx.Response.Header.Add("Location", uri)
					return true
				}
			}
		}
	}

	if !isRootRouter {
		return false
	}

	if method == OPTIONS {
		// Handle OPTIONS requests
		if r.router.HandleOPTIONS {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(HeaderAllow, allow)
				return true
			}
		}
	} else {
		// Handle 405
		if r.router.HandleMethodNotAllowed {
			if allow := r.router.Allowed(path, method); len(allow) > zero {
				ctx.Response.Header.Set(HeaderAllow, allow)
				if r.router.MethodNotAllowed != nil {
					r.router.MethodNotAllowed(ctx)
				} else {
					ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
					ctx.SetContentTypeBytes(DefaultContentType)
					ctx.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
				}
				return true
			}
		}
	}

	return false
}
