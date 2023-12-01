/*
 * Copyright (c) Kress4s
 * All rights reserved.
 *
 * Filename: client.go
 * Description: Single websocket client methods, this project can be used to connect to most users by user id.
 *
 * Created by Kress4s at 2023/12/01 15:05
 */

package ws

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketClient interface {
	Identifier() string
	// User() string
	ReadLoop()
	Write(message []byte)
	Close(id string) error
	Heart()
	WritePong() error
	IsClosed(id string) bool
}

func NewClient(conn *websocket.Conn, id string) WebsocketClient {
	return &websocketClient{
		ID:       id,
		conn:     conn,
		closed:   false,
		pingWait: time.Duration(PingWait) * time.Second,
		heart:    make(chan []byte),
	}

}

type websocketClient struct {
	// ws client connection
	conn *websocket.Conn
	// heart
	heart chan []byte
	// id
	ID string
	// user id -> can be used multi user
	// UserID string
	// ping ttl
	pingWait time.Duration
	// is closed or not
	closed bool
	// lock
	lock sync.Mutex
}

func (wc *websocketClient) Identifier() string {
	return wc.ID
}

// func (wc *websocketClient) User() string {
// 	return wc.UserID
// }

func (wc *websocketClient) ReadLoop() {
	defer func() {
		// occurred error, to exist
		SocketManager.Unregister <- wc
	}()
	for {
		_, msg, err := wc.conn.ReadMessage()
		if err != nil {
			log.Fatalf("client[%s] read message error, err is %s", wc.ID, err.Error())
		}
		switch string(msg) {
		case Close:
			log.Printf("client[%s] unregister", wc.ID)
			return
		case Ping:
			wc.heart <- []byte(Ping)
		}
		log.Printf("receive websocket message: [%s] ", string(msg))
	}
}

func (wc *websocketClient) Write(message []byte) {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if err := wc.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Fatalf("push message(%s) to client(%s) failed:[%s]", message, wc.ID, err)
	}
}

func (wc *websocketClient) Close(id string) error {
	wc.lock.Unlock()
	defer wc.lock.Unlock()
	if wc.closed {
		return nil
	}
	if err := wc.conn.Close(); err != nil {
		log.Fatalf("close client(%s) failed:[%s]", wc.ID, err)
		return err
	}
	wc.closed = true
	return nil
}

func (wc *websocketClient) Heart() {
	defer func() {
		SocketManager.Unregister <- wc
	}()
	for {
		select {
		case <-time.After(wc.pingWait):
			log.Fatalf("client(%s) pong timeout", wc.ID)
			return
		case <-wc.heart:
			if err := wc.WritePong(); err != nil {
				return
			}
		}
	}
}

func (wc *websocketClient) WritePong() error {
	wc.lock.Lock()
	defer wc.lock.Unlock()
	if err := wc.conn.WriteMessage(websocket.TextMessage, []byte(Pong)); err != nil {
		return err
	}
	return nil
}

func (wc *websocketClient) IsClosed(id string) bool {
	return wc.closed
}
