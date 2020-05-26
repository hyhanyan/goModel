package main

import (
	"fmt"
	"net/http"
	"time"

	"websocketProcess/impl"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *impl.Connection
		data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		/*
			// 一对一发送
			if err = conn.WriteMessage(data); err != nil {
				goto ERR
			}
		*/

		// 广播
		connMap := impl.GetConnPool()
		fmt.Println(impl.GetLengthConnPool())
		for _, ws := range connMap {
			if err = ws.WriteMessage(data); err != nil {
				goto ERR
			}
		}
	}

ERR:
	conn.Close()

}

func main() {

	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
