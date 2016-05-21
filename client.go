package gremgo

import (
	"io/ioutil"
	"sync"
)

// Client is a container for the gremgo client.
type Client struct {
	conn      connector
	requests  chan []byte
	responses chan []byte
	results   map[string]interface{}
	mutex     *sync.Mutex
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func Dial(conn connector) (c Client, err error) {

	// Initializes client
	c.conn = conn
	c.requests = make(chan []byte, 3)
	c.responses = make(chan []byte, 3)
	c.results = make(map[string]interface{})
	c.mutex = &sync.Mutex{}

	// Connects to Gremlin Server
	err = conn.connect()
	if err != nil {
		return
	}

	go c.writeWorker()
	go c.readWorker()

	return
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
func (c *Client) Execute(query string, bindings map[string]string) (response interface{}, err error) {
	req := evalRequest{query: query, bindings: bindings}
	c.sendRequest(&req)
	response = c.retrieveResponse(&req)
	return
}

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteFile(path string, bindings map[string]string) (response interface{}, err error) {
	d, err := ioutil.ReadFile(path) // Read script
	if err != nil {
		return
	}
	q := string(d)
	req := evalRequest{query: q, bindings: bindings} // TODO: Make this cleaner
	c.sendRequest(&req)
	response = c.retrieveResponse(&req)
	return
}
