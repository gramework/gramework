// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package client

import (
	"github.com/gramework/gramework"
)

// Handler returns gramework handler
func (client *Instance) Handler() func(*gramework.Context) error {
	return client.handleHTTP
}

func (client *Instance) handleHTTP(ctx *gramework.Context) error {
	api, err := client.nextServer()
	if err != nil {
		ctx.Logger.Errorf("error %s", err)
		return err
	}

	bytes := buffer.Get()
	defer buffer.Put(bytes)
	statusCode, body, err := api.HostClient.Get(bytes.B, api.Addr)
	if err != nil {
		ctx.Logger.Errorf("error while .Do() the request %s", err)
		return err
	}

	ctx.SetStatusCode(statusCode)
	_, err = ctx.Write(body)
	return err
}
