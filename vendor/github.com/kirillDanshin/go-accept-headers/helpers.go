// Copyright 2013 Ryan Rogers. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accept

import (
	"fmt"
	"strings"
)

var (
	errInvalidTypeSubtype = "accept: Invalid type '%s'."
)

// Len implements the Len() method of the Sort interface.
func (a AcceptSlice) Len() int {
	return len(a)
}

// Less implements the Less() method of the Sort interface.  Elements are
// sorted in order of decreasing preference.
func (a AcceptSlice) Less(i, j int) bool {
	// Higher qvalues come first.
	if a[i].Q > a[j].Q {
		return true
	} else if a[i].Q < a[j].Q {
		return false
	}

	// Specific types come before wildcard types.
	if a[i].Type != "*" && a[j].Type == "*" {
		return true
	} else if a[i].Type == "*" && a[j].Type != "*" {
		return false
	}

	// Specific subtypes come before wildcard subtypes.
	if a[i].Subtype != "*" && a[j].Subtype == "*" {
		return true
	} else if a[i].Subtype == "*" && a[j].Subtype != "*" {
		return false
	}

	// A lot of extensions comes before not a lot of extensions.
	if len(a[i].Extensions) > len(a[j].Extensions) {
		return true
	}

	return false
}

// Swap implements the Swap() method of the Sort interface.
func (a AcceptSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// parseMediaRange parses the provided media range, and on success returns the
// parsed range params and type/subtype pair.
func parseMediaRange(mediaRange string) (rangeParams, typeSubtype []string, err error) {
	rangeParams = strings.Split(mediaRange, ";")
	typeSubtype = strings.Split(rangeParams[0], "/")

	// typeSubtype should have a length of exactly two.
	if len(typeSubtype) > 2 {
		err = fmt.Errorf(errInvalidTypeSubtype, rangeParams[0])
		return
	} else {
		typeSubtype = append(typeSubtype, "*")
	}

	// Sanitize typeSubtype.
	typeSubtype[0] = strings.TrimSpace(typeSubtype[0])
	typeSubtype[1] = strings.TrimSpace(typeSubtype[1])
	if typeSubtype[0] == "" {
		typeSubtype[0] = "*"
	}
	if typeSubtype[1] == "" {
		typeSubtype[1] = "*"
	}

	return
}
