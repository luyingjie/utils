package websocket

import "testing"

func TestMain(t *testing.T) {

}

// var d manage.IConnection = ws.NewWebSocket("", userID, c.Writer, c.Request, func(ws *ws.WebSocket) {
// 	for {
// 		msgType, message, err := ws.ReadMessage()
// 		if err != nil {
// 			fmt.Println("read msg fail", err.Error())
// 			break
// 		}
// 		fmt.Println("recv msg", string(message), msgType)
// 		ws.WriteMessage(1, message)
// 	}
// 	// select {}
// })
// global.WSConn.Add(d)
// d.Start()

// function conn() {
// 	var ws = new WebSocket("ws://"+window.location.host+"/ws");

// 	//连接打开时触发
// 	ws.onopen = function (evt) {
// 		console.log("Connection open ...");
// 		ws.send("Hello WebSockets test!");
// 	};

// 	//接收到消息时触发
// 	ws.onmessage = function (evt) {
// 		if (evt.data == "logout") {
// 			exit();
// 		}
// 	};
// 	//连接关闭时触发
// 	ws.onclose = function (evt) {
// 		close();
// 	};
// }
