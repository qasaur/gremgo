package gremgo

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type dialer interface {
	connect() error
	write([]byte) error
	read() ([]byte, error)
	close()
}

/////
/*
WebSocket Connection
*/
/////

// Ws is the dialer for a WebSocket connection
type Ws struct {
	host string
	conn *websocket.Conn
}

func (ws *Ws) connect() (err error) {
	d := websocket.Dialer{
		WriteBufferSize:  8192,
		ReadBufferSize:   8192,
		HandshakeTimeout: 5 * time.Second, // Timeout or else we'll hang forever and never fail on bad hosts.
	}
	ws.conn, _, err = d.Dial(ws.host, http.Header{})
	if err != nil {

		// As of 3.2.2 the URL has changed.
		// https://groups.google.com/forum/#!msg/gremlin-users/x4hiHsmTsHM/Xe4GcPtRCAAJ
		ws.host = ws.host + "/gremlin"
		ws.conn, _, err = d.Dial(ws.host, http.Header{})
	}
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

func (ws *Ws) close() (err error) {
	err = ws.conn.Close()
	return
}

/////

func (c *Client) writeWorker() { // writeWorker works on a loop and dispatches messages as soon as it recieves them
	for {
		select {
		case msg := <-c.requests:
			err := c.conn.write(msg)
			if err != nil {
				log.Println(err)
				c.Errored = true
				break
			}
		}
	}
}

func (c *Client) readWorker() { // readWorker works on a loop and sorts messages as soon as it recieves them
	for {
		msg, err := c.conn.read()
		if err != nil {
			log.Println(err)
			c.Errored = true
			break
		}
		if msg != nil {
			c.handleResponse(msg)
		}
	}
}
