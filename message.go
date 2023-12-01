/*
 * Copyright (c) Kress4s
 * All rights reserved.
 *
 * Filename: message.go
 * Description: common message define
 *
 * Created by Kress4s at 2023/12/01 15:05
 */

package websocket_sdk

type Message struct {
	ID       string
	Content  string
	Sender   string
	Receiver string
}
