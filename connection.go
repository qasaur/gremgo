package gremgo

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type dialer interface {
	connect() error
	write([]byte) error
	read() ([]byte, error)
}

/////
/*
WebSocket Connection
*/
/////

// Ws is the dialer for a WebSocket connection
type Ws struct {
	Host string
	conn *websocket.Conn
}

func (ws *Ws) connect() (err error) {
	d := websocket.Dialer{}
	ws.conn, _, err = d.Dial(ws.Host, http.Header{})
	return
}

func (ws *Ws) write(msg []byte) (err error) {
	err = ws.conn.WriteMessage(2, msg)
	return
}

func (ws *Ws) read() (msg []byte, err error) {
	_, msg, err = ws.conn.ReadMessage()
	return
}

/////

func (c *Client) writeWorker() { // writeWorker works on a loop and dispatches messages as soon as it recieves them
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

func (c *Client) readWorker() { // readWorker works on a loop and sorts messages as soon as it recieves them
	for {
		msg, err := c.conn.read()
		if err != nil {
			log.Fatal(err)
		}
		if msg != nil {
			go c.handleResponse(msg) // Send message for sorting and retrieval
		}
	}
}
