package websocket

import (
	"bytes"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"fmt"
	"demo/lib/test"
	"github.com/golang/protobuf/proto"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	Conns *Conn
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) read() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		fmt.Println(message)
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("read", message)
		c.Conns.broadcast <- message
		t := &test.TestMsg{
			Channel: 1,
			MsgType: 2,
			Msg: []string{"test"},
		}
		data, err := proto.Marshal(t)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		newTest := &test.TestMsg{}
		err = proto.Unmarshal(data, newTest)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}
		// Now test and newTest contain the same data.
		if t.GetChannel() != newTest.GetChannel() {
			log.Fatalf("data mismatch %q != %q", t.GetChannel(), newTest.GetChannel())
		}
		c.conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg,ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			fmt.Println("write" ,msg)
			w.Write(msg)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func connHandle (conn *Conn, w http.ResponseWriter, r *http.Request) {
	newconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Conns: conn, conn: newconn, send: make(chan []byte, 256)}
	client.Conns.register <- client

	go client.read()
	go client.write()
}


