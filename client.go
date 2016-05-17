package gremgo

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is a container for the gremgo client.
type Client struct {
	host       string
	connection bool
	requests   chan []byte
	responses  map[string]interface{}
	mutex      *sync.Mutex
	buffer     map[string]map[int]interface{}
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func Dial(host string) (c Client, err error) {

	// Initializes client

	c.host = "ws://" + host + ":8182"
	c.requests = make(chan []byte, 3)
	c.responses = make(map[string]interface{})
	c.buffer = make(map[string]map[int]interface{})

	c.mutex = &sync.Mutex{}

	// Connect to websocket

	d := websocket.Dialer{}
	ws, _, err := d.Dial(c.host, http.Header{})
	if err != nil {
		return
	}

	c.connection = true

	// Write worker
	go func() {
		for c.connection == true {
			select {
			case msg := <-c.requests: // Wait for message send request
				err = ws.WriteMessage(2, msg) // Write message
				if err != nil {
					log.Fatal(err)
				}
			default:
			}
		}
	}()

	// Read worker
	go func() {
		for c.connection == true {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			if msg != nil {
				go c.handleResponse(msg)
			}
		}
	}()

	return
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
func (c *Client) Execute(query string, bindings map[string]string) (response interface{}, err error) {
	req := evalRequest{query: query, bindings: bindings}
	c.createRequest(&req)
	c.sendRequest(req.prepared)
	response = c.retrieveResponse(req.request.Requestid)
	return
}

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteFile(path string, bindings map[string]string) (response map[string]interface{}, err error) {
	// s, err := ioutil.ReadFile(path) // Read script
	// if err != nil {
	// 	return
	// }
	// req := evalRequest{query: s, bindings: bindings}
	// c.createRequest(req)
	return
}
