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

import "testing"

func TestWhitelist(t *testing.T) {
	type tcase struct {
		port     string
		expected bool
	}
	cases := []tcase{
		{
			port:     "80",
			expected: true,
		},
		{
			port:     "443",
			expected: true,
		},
		{
			port:     "80-8080",
			expected: true,
		},
		{
			port:     "980-3330",
			expected: false,
		},
		{
			port:     "0",
			expected: false,
		},
		{
			port:     "",
			expected: false,
		},
		{
			port:     "21",
			expected: false,
		},
	}

	for _, testcase := range cases {
		if IsPortInRange(testcase.port) != testcase.expected {
			t.Errorf("unexpected result: for port %q expected %v", testcase.port, testcase.expected)
		}
	}
}
