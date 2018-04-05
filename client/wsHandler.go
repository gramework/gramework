package client

import (
	"time"

	"github.com/fasthttp-contrib/websocket"
	"github.com/gramework/gramework"
)

// WSHandler returns gramework handler
func (client *Instance) WSHandler() func(*gramework.Context) error {
	return func(ctx *gramework.Context) error {
		if websocket.IsWebSocketUpgrade(ctx.RequestCtx) {
			return websocket.Upgrade(ctx.RequestCtx, func(conn *websocket.Conn) {
				for {
					v := <-client.watch(ctx)
					conn.WriteMessage(websocket.TextMessage, v)
				}
			}, 0, 0)
		}

		return client.handleHTTP(ctx)
	}
}

func (client *Instance) watch(ctx *gramework.Context) chan []byte {
	c := make(chan []byte)
	go func() {
		for {
			api, err := client.nextServer()
			if err != nil {
				ctx.Logger.Errorf("error: %s", err)
				time.Sleep(client.conf.WatcherTickTime)
				continue
			}
			bytes := buffer.Get()
			_, body, err := api.HostClient.Get(bytes.B, api.Addr)
			if err != nil {
				ctx.Logger.Errorf("error while .Do() the request %s", err)
				time.Sleep(client.conf.WatcherTickTime)
				buffer.Put(bytes)
				continue
			}
			c <- body
			time.Sleep(client.conf.WatcherTickTime)
			buffer.Put(bytes)
		}
	}()

	return c
}
