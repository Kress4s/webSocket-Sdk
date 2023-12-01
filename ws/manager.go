/*
 * Copyright (c) Kress4s
 * All rights reserved.
 *
 * Filename: client.go
 * Description: to manage websocket clients methods
 *
 * Created by Kress4s at 2023/12/01 15:05
 */

package ws

import "sync"

var SocketManager = socketManager{
	clients:    make(map[string]WebsocketClient),
	Broadcast:  make(chan []byte),
	Register:   make(chan WebsocketClient),
	Unregister: make(chan WebsocketClient),
}

type socketManager struct {
	clients map[string]WebsocketClient
	// the user has subscribed ws, (eg system message, user message, etc)
	// notice    map[string]string
	Broadcast chan []byte
	// client register
	Register chan WebsocketClient
	// client login out
	Unregister chan WebsocketClient
	mtx        sync.Mutex
}

// add clint to manager
func (m *socketManager) Add(id string, client WebsocketClient) {
	m.mtx.Lock()
	m.clients[id] = client
	m.mtx.Unlock()
}

// func (m *socketManager) AddUser(userID, id string) {
// 	m.mtx.Lock()
// 	m.notice[userID] = id
// 	m.mtx.Unlock()
// }

func (m *socketManager) Get(id string) WebsocketClient {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if conn, ok := m.clients[id]; ok {
		return conn
	}
	return nil
}

func (m *socketManager) UnRegister(id string) {
	m.mtx.Lock()
	delete(m.clients, id)
	// for k, v := range m.notice {
	// 	if v == id {
	// 		delete(m.notice, k)
	// 		break
	// 	}
	// }
	m.mtx.Unlock()
}

// for each client to run f()
func (m *socketManager) Range(f func(string, WebsocketClient)) {
	m.mtx.Lock()
	for k, v := range m.clients {
		f(k, v)
	}
	m.mtx.Unlock()
}

// func (m *socketManager) RangeUser(f func(string, WebsocketClient)) {
// 	m.mtx.Lock()
// 	for k, v := range m.notice {
// 		if c, ok := m.clients[v]; ok {
// 			f(k, c)
// 		}
// 	}
// 	m.mtx.Unlock()
// }
