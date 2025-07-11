package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	Connection      *websocket.Conn
	ConnectionMutex *sync.Mutex
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{
		Connection:      nil,
		ConnectionMutex: &sync.Mutex{},
	}
}

func (wsh *WebsocketHandler) WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error in websocket upgrade: %v\n", err)
		return
	}
	wsh.ConnectionMutex.Lock()
	wsh.Connection = c
	wsh.ConnectionMutex.Unlock()

	// Wait for the client to disconnect
	for {
		if _, _, err := wsh.Connection.NextReader(); err != nil {
			wsh.ConnectionMutex.Lock()
			wsh.Connection = nil
			wsh.ConnectionMutex.Unlock()
			c.Close()
			break
		}
	}
}

func (wsh *WebsocketHandler) SendMessage(message any) error {
	if wsh.ConnectionMutex == nil {
		return fmt.Errorf("websocket connection mutex is not initialized")
	}
	wsh.ConnectionMutex.Lock()
	defer wsh.ConnectionMutex.Unlock()

	if wsh.Connection == nil {
		return fmt.Errorf("no WebSocket connection available")
	}

	if err := wsh.Connection.WriteJSON(message); err != nil {
		return fmt.Errorf("error sending WebSocket message: %w", err)
	}
	return nil
}
