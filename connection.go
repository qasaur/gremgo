package gremgo

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type connector interface {
	connect() error
	write([]byte) error
	read() ([]byte, error)
}

/////
/*
WebSocket Connection
*/
/////

type ws struct {
	host string
	conn *websocket.Conn
}

func (ws *ws) connect() (err error) {
	d := websocket.Dialer{}
	ws.conn, _, err = d.Dial(ws.host, http.Header{})
	return
}

func (ws *ws) write(msg []byte) (err error) {
	err = ws.conn.WriteMessage(2, msg)
	return
}

func (ws *ws) read() (msg []byte, err error) {
	_, msg, err = ws.conn.ReadMessage()
	return
}

/////

func (c *Client) writeWorker() {
	for {
		select {
		case msg := <-c.requests: // Wait for message send request
			err := c.conn.write(msg) // Write message
			if err != nil {
				log.Fatal(err)
			}
		default:
		}
	}
}

func (c *Client) readWorker() {
	for {
		msg, err := c.conn.read()
		if err != nil {
			log.Fatal(err)
		}
		if msg != nil {
			c.responses <- msg
		}
	}
}
