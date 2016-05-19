package gremgo

import "sync"

type client interface {
}

// Client is a container for the gremgo client.
type Client struct {
	conn      connector
	requests  chan []byte
	responses map[string]interface{}
	mutex     *sync.Mutex
	buffer    map[string]map[int]interface{}
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func Dial(host string) (c Client, err error) {

	// Initializes client
	c.requests = make(chan []byte, 3)
	c.responses = make(map[string]interface{})
	c.buffer = make(map[string]map[int]interface{})
	c.mutex = &sync.Mutex{}
	c.conn = &ws{host: host}

	err = c.conn.connect()
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
