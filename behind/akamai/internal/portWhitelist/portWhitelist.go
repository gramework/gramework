// Package portWhitelist provides a parser
// for the Akamai's CIDR blocks Port field syntax.
// This is not an official Akamai-supported implementation.
// If you having any issues with this package, please
// consider to contact Gramework support first.
// Akamai doesn't provide any official support nor guaranties
// about this package.
//
// Akamai is a trademark of Akamai Technologies, Inc.
//
// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
package portWhitelist

import (
	"strconv"
	"strings"
)

var allowedPorts = []string{"80", "443"}
var allowedPortsInts = []int{80, 443}

// IsPortInRange checks if port is in whitelist.
func IsPortInRange(port string) bool {
	if strings.Contains(port, ",") {
		ports := strings.Split(port, ",")
		for _, p := range ports {
			if strings.Contains(p, "-") && IsPortInRange(p) {
				return true
			}
			for _, allowedPort := range allowedPorts {
				if p == allowedPort {
					return true
				}
			}
		}

		return false
	}

	if strings.Contains(port, "-") {
		portRange := strings.Split(port, "-")
		if len(portRange) != 2 {
			return false
		}

		min, err := strconv.Atoi(portRange[0])
		if err != nil {
			return false
		}

		max, err := strconv.Atoi(portRange[1])
		if err != nil {
			return false
		}
		for _, p := range allowedPortsInts {
			if min <= p && max >= p {
				return true
			}
		}
		return false
	}

	for _, p := range allowedPorts {
		if port == p {
			return true
		}
	}

	return false
}
