package websocket

import "net/http"

func WsHander(w http.ResponseWriter, r *http.Request) {
	conn := newConn()
	go conn.run()
	connHandle(conn, w, r)
}