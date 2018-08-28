// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package client

import (
	"net/url"

	"github.com/gramework/gramework"
)

func (client *Instance) nextServer() (*requestInfo, error) {
	if len(client.conf.Addresses) == 0 {
		return nil, ErrNoServerAvailable
	}

	for i := 0; i < len(client.conf.Addresses); i++ {
		addr := client.conf.Addresses[client.balancer.next()]
		hostURL, err := url.Parse(addr)
		if err != nil {
			gramework.Errorf("error while parsing host url: %s", err)
			continue
		}

		return &requestInfo{
			HostClient: client.getHostClient(hostURL),
			Addr:       addr,
		}, nil
	}

	return nil, ErrNoServerAvailable
}
