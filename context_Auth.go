package gramework

import (
	"bytes"
	"encoding/base64"
	"errors"
)

var authHeaderName = []byte("Authorization")
var authSplitter = []byte(":")
var errInvalidAuth = errors.New("invalid auth request")

// GetPass lazy triggers parser and returns
// password or an error. Error will be persistent
// for current context
//
// Common typos: GetPassword
func (a *Auth) GetPass() (string, error) {
	// yep, we copying the code, but
	// we get one instead of two jumps
	if !a.parsed {
		a.parse()
	}

	return a.pass, a.err
}

// GetLogin lazy triggers parser and returns
// login or an error. Error will be persistent
// for current context
//
// Common typos: GetUser, GetUsername
func (a *Auth) GetLogin() (string, error) {
	if !a.parsed {
		a.parse()
	}

	return a.login, a.err
}

func (a *Auth) parse() {
	if a.err != nil {
		return // parsing already failed
	}
	auth := a.ctx.Request.Header.PeekBytes(authHeaderName)
	if len(auth) < 7 {
		a.err = errInvalidAuth
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(BytesToString(auth[6:]))
	if err != nil {
		a.err = err
		return
	}

	splitted := bytes.Split(decoded, authSplitter)

	if len(splitted) != 2 {
		a.err = errInvalidAuth
		return
	}

	a.login, a.pass = string(splitted[0]), string(splitted[1])
}

// Auth returns struct for simple basic auth handling
//
// useful to develop e.g. stage environment login,
// where high security is not required
func (ctx *Context) Auth() *Auth {
	if ctx.auth == nil {
		ctx.auth = &Auth{
			ctx: ctx,
		}
	}
	return ctx.auth
}
