package routers

import (
	"demo/lib/websocket"
	"net/http"
)
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func InitRouter() Handler {
	//SayHello
	Handler := http.HandlerFunc(websocket.WsHander)
	return Handler
}