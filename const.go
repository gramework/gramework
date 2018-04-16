// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

const (
	// MethodDELETE is the HTTP DELETE method
	MethodDELETE = "DELETE"

	// MethodGET is the HTTP GET method
	MethodGET = "GET"

	// MethodHEAD is the HTTP HEAD method
	MethodHEAD = "HEAD"

	// MethodOPTIONS is the HTTP OPTIONS method
	MethodOPTIONS = "OPTIONS"

	// MethodPATCH is the HTTP PATCH method
	MethodPATCH = "PATCH"

	// MethodPOST is the HTTP POST method
	MethodPOST = "POST"

	// MethodPUT is the HTTP PUT method
	MethodPUT = "PUT"
)

const (
	acceptHeader                      = "Accept"
	badRequest                        = "Bad Request"
	corsAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	corsAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	corsAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	corsAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	corsCType                         = "Content-Type, *"
	emptyString                       = ""
	fmtV                              = "%v"
	forbidden                         = "Forbidden"
	hOrigin                           = "Origin"
	htmlCT                            = "text/html; charset=utf8"
	https                             = "https"
	jsonCT                            = "application/json;charset=utf8"
	methods                           = "GET,PUT,POST,DELETE"
	plainCT                           = "text/plain"
	trueStr                           = "true"
	xmlCT                             = "text/xml"
	forbiddenCode                     = 403
	redirectCode                      = 301
	temporaryRedirectCode             = 307
	one                               = 1
	zero                              = 0
)
