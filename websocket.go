/*
 * Copyright (c) 1998-2023 江苏斯菲尔电气股份有限公司
 * All rights reserved.
 *
 * Filename: websocket.go
 * Description:
 *
 * Created by xiayoushuang at 2023/12/01 16:32
 * http://www.sfere-elec.com/
 */

package websocket_sdk

import (
	"log"

	"github.com/Kress4s/websocket-sdk/ws"

	"github.com/goccy/go-json"
)

func Start() {
	for {
		select {
		case conn := <-ws.SocketManager.Register:
			msg := "websocket is connected"
			sendMsg, err := json.Marshal(&Message{Content: msg})
			if err != nil {
				conn.Write([]byte(msg))
			} else {
				conn.Write(sendMsg)
			}
			ws.SocketManager.Add(conn.Identifier(), conn)
			// ws.SocketManager.AddUser(conn.User(), conn.Identifier())
		case c := <-ws.SocketManager.Unregister:
			if conn := ws.SocketManager.Get(c.Identifier()); conn != nil {
				if err := conn.Close(conn.Identifier()); err != nil {
					log.Fatalf("close client(%s) failed:[%s]", conn.Identifier(), err)
				}
				ws.SocketManager.UnRegister(conn.Identifier())
			}
		case message := <-ws.SocketManager.Broadcast:
			ws.SocketManager.Range(func(_ string, conn ws.WebsocketClient) {
				go conn.Write(message)
			})
		}
	}
}
