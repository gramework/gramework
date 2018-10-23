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
	"fmt"
)

// GET registers a handler for a GET request to the given route
func (r *SubRouter) GET(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodGET, route, handler, r.prefixes)
	}
	return r
}

func (r *SubRouter) determineHandler(handler interface{}) func(*Context) {
	return r.parent.determineHandler(handler)
}

// JSON register internal handler that sets json content type
// and serves given handler with GET method
func (r *SubRouter) JSON(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		h := r.parent.determineHandler(handler)
		r.parent.handleReg(MethodGET, route, jsonHandler(h), r.prefixes)
	}

	return r
}

// DELETE registers a handler for a DELETE request to the given route
func (r *SubRouter) DELETE(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodDELETE, route, handler, r.prefixes)
	}

	return r
}

// HEAD registers a handler for a HEAD request to the given route
func (r *SubRouter) HEAD(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodHEAD, route, handler, r.prefixes)
	}

	return r
}

// ServeFile serves a file on a given route
func (r *SubRouter) ServeFile(route, file string) *SubRouter {
	route = r.prefixedRoute(route)
	r.parent.handleReg(MethodGET, route, func(ctx *Context) {
		ctx.SendFile(file)
	}, r.prefixes)

	return r
}

// OPTIONS registers a handler for a OPTIONS request to the given route
func (r *SubRouter) OPTIONS(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodOPTIONS, route, handler, r.prefixes)
	}

	return r
}

// PUT registers a handler for a PUT request to the given route
func (r *SubRouter) PUT(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPUT, route, handler, r.prefixes)
	}

	return r
}

// POST registers a handler for a POST request to the given route
func (r *SubRouter) POST(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPOST, route, handler, r.prefixes)
	}

	return r
}

// PATCH registers a handler for a PATCH request to the given route
func (r *SubRouter) PATCH(route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(MethodPATCH, route, handler, r.prefixes)
	}

	return r
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut functions can be used.
// This function is intended for bulk loading and to allow the usage of less frequently used,
// non-standardized or custom methods (e.g. for internal communication with a proxy).
func (r *SubRouter) Handle(method, route string, handler interface{}) *SubRouter {
	route = r.prefixedRoute(route)
	if r.parent != nil {
		r.parent.handleReg(method, route, handler, r.prefixes)
	}
	return r
}

func (r *SubRouter) handleReg(method, route string, handler interface{}, prefixes []string) {
	r.parent.handleReg(method, route, handler, prefixes)
}

// Sub let you quickly register subroutes with given prefix
// like app.Sub("v1").Sub("users").GET("view/:id", "hi").DELETE("delete/:id", "hi"),
// that give you /v1/users/view/:id and /v1/users/delete/:id registered
func (r *SubRouter) Sub(path string) *SubRouter {
	return &SubRouter{
		parent:   r,
		prefix:   r.prefixedRoute(path),
		prefixes: append(r.prefixes, path),
	}
}

func (r *SubRouter) prefixedRoute(route string) string {
	if r.prefix[len(r.prefix)-1] != '/' && route[0] != '/' {
		return fmt.Sprintf("%s/%s", r.prefix, route)
	}

	return fmt.Sprintf("%s%s", r.prefix, route)
}

// HTTP returns SubRouter for http requests with given r.prefix
func (r *SubRouter) HTTP() *SubRouter {
	switch parent := r.parent.(type) {
	case *SubRouter:
		return parent.HTTP()
	case *Router:
		return &SubRouter{
			parent: parent,
			prefix: r.prefix,
		}
	default:
		Errorf("[HIGH SEVERITY BUG]: unreachable case found! Expected *SubRouter or *Router, got %T! Returning nil!", parent)
		Errorf("Please report the bug on https://github.com/gramework/gramework ASAP!")
		return nil
	}
}

// HTTPS returns SubRouter for https requests with given r.prefix
func (r *SubRouter) HTTPS() *SubRouter {
	switch parent := r.parent.(type) {
	case *SubRouter:
		return parent.HTTPS()
	case *Router:
		return &SubRouter{
			parent: parent,
			prefix: r.prefix,
		}
	default:
		Errorf("[HIGH SEVERITY BUG]: unreachable case found! Expected *SubRouter or *Router, got %T! Returning nil!", parent)
		Errorf("Please report the bug on https://github.com/gramework/gramework ASAP!")
		return nil
	}
}

// ToTLSHandler returns handler that redirects user to HTTP scheme
func (r *SubRouter) ToTLSHandler() func(*Context) {
	return func(ctx *Context) {
		ctx.ToTLS()
	}
}

// Forbidden is a shortcut for ctx.Forbidden
func (r *SubRouter) Forbidden(ctx *Context) {
	ctx.Forbidden()
}

// Redir sends 301 redirect to the given url
//
// it's equivalent to
//
//     ctx.Redirect(url, 301)
func (r *SubRouter) Redir(route, url string) {
	r.GET(route, func(ctx *Context) {
		ctx.Redirect(route, redirectCode)
	})
}
