// ws的连接包，基于 "github.com/gorilla/websocket" 封装，实现了IConnection接口，可以用于连接管理。
package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler func(*WebSocket)

type WebSocket struct {
	key     string
	Path    string
	Handler WebSocketHandler

	w http.ResponseWriter
	r *http.Request
	*websocket.Upgrader
	*websocket.Conn
}

func NewWebSocket(path string, key string, w http.ResponseWriter, r *http.Request, handler WebSocketHandler, checkOrigin ...bool) *WebSocket {
	var upg websocket.Upgrader
	// 设置 checkOrigin 是否允许跨域请求。默认不开启， true时候开启。
	if len(checkOrigin) > 0 && checkOrigin[0] {
		upg = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	} else {
		upg = websocket.Upgrader{}
	}
	return &WebSocket{
		key:      key,
		Path:     path,
		w:        w,
		r:        r,
		Upgrader: &upg,
		Handler:  handler,
	}
}

func (ws *WebSocket) Start() error {
	conn, err := ws.Upgrader.Upgrade(ws.w, ws.r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ws.Conn = conn
	ws.Handler(ws)
	return nil
}

func (ws *WebSocket) GetConnID() string {
	return ws.key
}

func (ws *WebSocket) Stop() error {
	err := ws.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ws *WebSocket) GetConnection() interface{} {
	return ws
}

func (ws *WebSocket) SendMsg(id int, data []byte) error {
	return ws.WriteMessage(id, data)
}
