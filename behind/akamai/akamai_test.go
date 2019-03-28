// Package akamai provideds a gramework.Behind implementation
// developed for Gramework.
// This is not an official Akamai-supported implementation.
// If you having any issues with this package, please
// consider to contact Gramework support first.
// Akamai doesn't provide any official support nor guaranties
// about this package.
//
// Akamai is a trademark of Akamai Technologies, Inc.
//
// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package akamai

import (
	"net"
	"testing"
)

// Note: those IP CIDRs are fake.
// Please, download a fresh CIDRs you need to whitelist
// in your Luna Control Panel.
const csvData = `Service Name,CIDR Block,Port,Activation Date,CIDR Status
"Log Delivery","120.33.22.0/24","21","Tue Dec 18 2021 02:00:00 GMT+0200 (Москва, стандартное время)","current"
"Log Delivery","120.33.21.0/24","80,443","Tue Dec 18 2107 02:00:00 GMT+0200 (Москва, стандартное время)","current"
"Log Delivery","120.33.23.0/24","80-8080","Tue Dec 18 1507 02:00:00 GMT+0200 (Москва, стандартное время)","current"
"Log Delivery","120.33.24.0/24","980-3300","Tue Dec 18 5507 02:00:00 GMT+0200 (Москва, стандартное время)","current"
"Log Delivery","120.17.33.0/24","21","Tue Dec 18 6507 02:00:00 GMT+0200 (Москва, стандартное время)","current"
`

func TestParseCIDRBlocksCSV(t *testing.T) {
	cidrs, err := ParseCIDRBlocksCSV([]byte(csvData), true, true)
	if err != nil {
		t.Error(err)
	}

	type tcase struct {
		cidr     *net.IPNet
		expected bool
	}
	cases := []tcase{
		{
			expected: true,
			cidr:     parseCIDR("120.33.21.0/24"),
		},
		{
			expected: true,
			cidr:     parseCIDR("120.33.23.0/24"),
		},

		{
			expected: false,
			cidr:     parseCIDR("120.33.22.0/24"),
		},
		{
			expected: false,
			cidr:     parseCIDR("120.33.24.0/24"),
		},
		{
			expected: false,
			cidr:     parseCIDR("120.17.33.0/24"),
		},
	}

	for _, testcase := range cases {
		found := false

		for _, cidr := range cidrs {
			if cidr.String() == testcase.cidr.String() {
				found = true
				break
			}
		}

		if found != testcase.expected {
			t.Errorf("unexpected result: CIDR %q expected=%v", testcase.cidr.String(), testcase.expected)
			return
		}
	}
}

func parseCIDR(raw string) *net.IPNet {
	_, cidr, _ := net.ParseCIDR(raw)
	return cidr
}
