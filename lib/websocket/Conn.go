package websocket

import "fmt"

type Conn struct {
	//在线用户
	clients map[*Client]bool
	//广播消息
	broadcast chan []byte


	//注册连接
	register chan *Client
	//释放连接
	unregister chan *Client
}

func newConn() *Conn {
	return &Conn{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister:  make(chan *Client),
	}
}

func (conn *Conn) run() {
	for {
		select {
		case client := <-conn.register:
			conn.clients[client] = true
		case client := <-conn.unregister:
			if _,ok := conn.clients[client]; ok {
				delete(conn.clients, client)
			}
		case message := <-conn.broadcast:
			fmt.Println("recv", message)
		}
	}
}
