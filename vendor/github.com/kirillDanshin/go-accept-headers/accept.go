// Copyright 2013 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package accept allows for easy handling of HTTP Accept headers.
// Accept-Ranges is currently not handled.
package accept

import (
	"sort"
	"strconv"
	"strings"
)

// Accept represents a parsed Accept(-Charset|-Encoding|-Language) header.
type Accept struct {
	Type, Subtype string
	Q             float64
	Extensions    map[string]string
}

// AcceptSlice is a slice of Accept.
type AcceptSlice []Accept

// Parse parses a HTTP Accept(-Charset|-Encoding|-Language) header and returns
// AcceptSlice, sorted in decreasing order of preference.  If the header lists
// multiple types that have the same level of preference (same specificity of
// type and subtype, same qvalue, and same number of extensions), the type
// that was listed in the header first comes first in the returned value.
//
// See http://www.w3.org/Protocols/rfc2616/rfc2616-sec14 for more information.
func Parse(header string) AcceptSlice {
	mediaRanges := strings.Split(header, ",")
	accepted := make(AcceptSlice, 0, len(mediaRanges))
	for _, mediaRange := range mediaRanges {
		rangeParams, typeSubtype, err := parseMediaRange(mediaRange)
		if err != nil {
			continue
		}

		accept := Accept{
			Type:       typeSubtype[0],
			Subtype:    typeSubtype[1],
			Q:          1.0,
			Extensions: make(map[string]string),
		}

		// If there is only one rangeParams, we can stop here.
		if len(rangeParams) == 1 {
			accepted = append(accepted, accept)
			continue
		}

		// Validate the rangeParams.
		validParams := true
		for _, v := range rangeParams[1:] {
			nameVal := strings.SplitN(v, "=", 2)
			if len(nameVal) != 2 {
				validParams = false
				break
			}
			nameVal[1] = strings.TrimSpace(nameVal[1])
			if name := strings.TrimSpace(nameVal[0]); name == "q" {
				qval, err := strconv.ParseFloat(nameVal[1], 64)
				if err != nil || qval < 0 {
					validParams = false
					break
				}
				if qval > 1.0 {
					qval = 1.0
				}
				accept.Q = qval
			} else {
				accept.Extensions[name] = nameVal[1]
			}
		}

		if validParams {
			accepted = append(accepted, accept)
		}
	}

	sort.Sort(accepted)
	return accepted
}

// Negotiate returns a type that is accepted by both the header declaration,
// and the list of types provided.  If no common types are found, an empty
// string is returned.
func Negotiate(header string, ctypes ...string) (string, error) {
	a := Parse(header)
	return a.Negotiate(ctypes...)
}

// Negotiate returns a type that is accepted by both the AcceptSlice, and the
// list of types provided.  If no common types are found, an empty string is
// returned.
func (accept AcceptSlice) Negotiate(ctypes ...string) (string, error) {
	if len(ctypes) == 0 {
		return "", nil
	}

	typeSubtypes := make([][]string, 0, len(ctypes))
	for _, v := range ctypes {
		_, ts, err := parseMediaRange(v)
		if err != nil {
			return "", err
		}
		if ts[0] == "*" && ts[1] == "*" {
			return v, nil
		}
		typeSubtypes = append(typeSubtypes, ts)
	}

	for _, a := range accept {
		for i, ts := range typeSubtypes {
			if ((a.Type == ts[0] || a.Type == "*") && (a.Subtype == ts[1] || a.Subtype == "*")) ||
				(ts[0] == "*" && ts[1] == a.Subtype) ||
				(ts[0] == a.Type && ts[1] == "*") {
				return ctypes[i], nil
			}
		}
	}
	return "", nil
}

// Accepts returns true if the provided type is accepted.
func (accept AcceptSlice) Accepts(ctype string) bool {
	t, err := accept.Negotiate(ctype)
	if t == "" || err != nil {
		return false
	}
	return true
}
