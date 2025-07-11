package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestWebsocketHandler_Upgrade tests that the handler upgrades the connection and stores it
func TestWebsocketHandler_Upgrade(t *testing.T) {
	wsh := NewWebsocketHandler()
	server := httptest.NewServer(http.HandlerFunc(wsh.WebsocketHandler))
	defer server.Close()

	// Convert http://127.0.0.1:port to ws://127.0.0.1:port
	url := "ws" + server.URL[len("http"):]

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	defer ws.Close()

	// Give the handler a moment to store the connection
	time.Sleep(100 * time.Millisecond)

	wsh.ConnectionMutex.Lock()
	defer wsh.ConnectionMutex.Unlock()
	if wsh.Connection == nil {
		t.Errorf("expected connection to be set")
	}
}

// TestWebsocketHandler_SendMessage tests sending a message over the websocket
func TestWebsocketHandler_SendMessage(t *testing.T) {
	wsh := NewWebsocketHandler()
	server := httptest.NewServer(http.HandlerFunc(wsh.WebsocketHandler))
	defer server.Close()

	url := "ws" + server.URL[len("http"):]

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("failed to dial websocket: %v", err)
	}
	defer ws.Close()

	// Wait for the handler to store the connection
	time.Sleep(100 * time.Millisecond)

	msg := map[string]any{"hello": "world"}
	err = wsh.SendMessage(msg)
	if err != nil {
		t.Errorf("SendMessage returned error: %v", err)
	}

	// Read the message from the websocket
	var received map[string]any
	if err := ws.ReadJSON(&received); err != nil {
		t.Errorf("failed to read JSON from websocket: %v", err)
	}
	if received["hello"] != "world" {
		t.Errorf("unexpected message: %+v", received)
	}
}

// TestWebsocketHandler_SendMessage_NoConnection tests SendMessage when no connection is present
func TestWebsocketHandler_SendMessage_NoConnection(t *testing.T) {
	wsh := NewWebsocketHandler()
	err := wsh.SendMessage(map[string]any{"test": 1})
	if err == nil {
		t.Errorf("expected error when no connection is available")
	}
}
