package gramework

import (
	"net/http"
)

const (
	// MethodDELETE is the HTTP DELETE method
	MethodDELETE = http.MethodDelete

	// MethodGET is the HTTP GET method
	MethodGET = http.MethodGet

	// MethodHEAD is the HTTP HEAD method
	MethodHEAD = http.MethodHead

	// MethodOPTIONS is the HTTP OPTIONS method
	MethodOPTIONS = http.MethodOptions

	// MethodPATCH is the HTTP PATCH method
	MethodPATCH = http.MethodPatch

	// MethodPOST is the HTTP POST method
	MethodPOST = http.MethodPost

	// MethodPUT is the HTTP PUT method
	MethodPUT = http.MethodPut
)

const (
	emptyString = ""

	fmtV = "%v"

	htmlCT  = "text/html; charset=utf8"
	jsonCT  = "application/json;charset=utf8"
	xmlCT   = "text/xml"
	plainCT = "text/plain"

	acceptHeader = "Accept"
)

const (
	badRequest                        = "Bad Request"
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsCType                         = "Content-Type, *"
	forbidden                         = "Forbidden"
	hOrigin                           = "Origin"
	https                             = "https"
	methods                           = "GET,PUT,POST,DELETE"
	trueStr                           = "true"
	forbiddenCode                     = http.StatusForbidden
	redirectCode                      = http.StatusMovedPermanently
	temporaryRedirectCode             = http.StatusTemporaryRedirect
	one                               = 1
	zero                              = 0
)
