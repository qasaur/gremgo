package gremgo

import (
	"io/ioutil"
	"sync"
)

// Client is a container for the gremgo client.
type Client struct {
	conn      dialer
	requests  chan []byte
	responses chan []byte
	results   map[string][]interface{}
	respMutex *sync.RWMutex
}

// NewDialer returns a WebSocket dialer to use when connecting to Gremlin Server
func NewDialer(host string) (dialer Ws) {
	return Ws{host: "127.0.0.1"}
}

func newClient() (c Client) {
	c.requests = make(chan []byte, 3)  // c.requests takes any request and delivers it to the WriteWorker for dispatch to Gremlin Server
	c.responses = make(chan []byte, 3) // c.responses takes raw responses from ReadWorker and delivers it for sorting to handelResponse
	c.results = make(map[string][]interface{})
	c.respMutex = &sync.RWMutex{} // c.mutex ensures that sorting is thread safe
	return
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func Dial(conn dialer) (c Client, err error) {

	c = newClient()
	c.conn = conn

	// Connects to Gremlin Server
	err = conn.connect()
	if err != nil {
		return
	}

	go c.writeWorker()
	go c.readWorker()

	return
}

func (c *Client) executeRequest(query string, bindings map[string]string) (resp interface{}, err error) {
	req, id := prepareRequest(query, bindings)
	msg, err := packageRequest(req)
	if err != nil {
		return // TODO: Fix error handling
	}
	c.dispatchRequest(msg)
	resp = c.retrieveResponse(id)
	return
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
func (c *Client) Execute(query string, bindings map[string]string) (resp interface{}, err error) {
	resp, err = c.executeRequest(query, bindings)
	return
}

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteFile(path string, bindings map[string]string) (resp interface{}, err error) {
	d, err := ioutil.ReadFile(path) // Read script from file
	if err != nil {
		return
	}
	query := string(d)
	resp, err = c.executeRequest(query, bindings)
	return
}
