package gremgo

import "sync"

// Client is a container for the gremgo client.
type Client struct {
	conn      dialer
	requests  chan []byte
	responses chan []byte
	results   map[string]interface{}
	mutex     *sync.Mutex
}

// Dial returns a gremgo client for interaction with the Gremlin Server specified in the host IP.
func Dial(conn dialer) (c Client, err error) {

	// Initializes client
	c.conn = conn
	c.requests = make(chan []byte, 3)  // c.requests takes any request and delivers it to the WriteWorker for dispatch to Gremlin Server
	c.responses = make(chan []byte, 3) // c.responses takes raw responses from ReadWorker and delivers it for sorting to handelResponse
	c.results = make(map[string]interface{})
	c.mutex = &sync.Mutex{} // c.mutex ensures that response sorting is safe

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
// func (c *Client) Execute(query string, bindings map[string]string) (response interface{}, err error) {
// 	req := evalRequest{query: query, bindings: bindings} // Execute defaults to an evaulation request
// 	// c.sendRequest(&req)
// 	response = c.retrieveResponse(&req)
// 	return
// }

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
// func (c *Client) ExecuteFile(path string, bindings map[string]string) (response interface{}, err error) {
// 	d, err := ioutil.ReadFile(path) // Read script from file
// 	if err != nil {
// 		return
// 	}
// 	q := string(d)
// 	req := evalRequest{query: q, bindings: bindings}
// 	// c.sendRequest(&req)
// 	response = c.retrieveResponse(&req)
// 	return
// }
