package websocket

import (
	"net/http"

	myerror "utils/error"

	"github.com/gorilla/websocket"
)

type WebSocketHandler func(*WebSocket)

type WebSocket struct {
	Path    string
	Handler WebSocketHandler

	*websocket.Upgrader
	*websocket.Conn
}

func NewWebSocket(path string, handler WebSocketHandler) *WebSocket {
	return &WebSocket{
		Path:     path,
		Upgrader: &websocket.Upgrader{},
		Handler:  handler,
	}
}

func (ws *WebSocket) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		myerror.Log(2000, 3, err)
		return
	}
	defer conn.Close()

	ws.Conn = conn
	ws.Handler(ws)
}
