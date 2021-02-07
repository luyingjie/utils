package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler func(*WebSocket)

type WebSocket struct {
	Key 	string
	Path    string
	Handler WebSocketHandler

	*websocket.Upgrader
	*websocket.Conn
}

func NewWebSocket(path string, key string, handler WebSocketHandler) *WebSocket {
	return &WebSocket{
		Key:  	  key,
		Path:     path,
		Upgrader: &websocket.Upgrader{},
		Handler:  handler,
	}
}

func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) error {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ws.Conn = conn
	ws.Handler(ws)
	return nil
}

func (ws *WebSocket) GetConnID() string {
	return ws.Key
}

func (ws *WebSocket) Stop() error {
	err := ws.Close()
	if err != nil {
		return err
	}
	return nil
}
