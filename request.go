package gremgo

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/satori/go.uuid"
)

// Request is a container for all request parameters to be sent to the Gremlin Server.
type request struct {
	Requestid string  `json:"requestId"`
	Op        string  `json:"op"`
	Processor string  `json:"processor"`
	Args      reqArgs `json:"args"`
}

//reqArgs define the arguments for the Gremlin request.
type reqArgs struct {
	Gremlin  string            `json:"gremlin"`
	Language string            `json:"language"`
	Bindings map[string]string `json:"bindings"`
}

// prepareMessage formats a query into the standard accepted by Gremlin Server along with its bindings.
func prepareMessage(query string, bindings map[string]string) (msg []byte, requestid string, err error) {
	var req request
	var args reqArgs

	args.Gremlin = query
	args.Language = "gremlin-groovy"
	args.Bindings = bindings

	req.Args = args
	req.Requestid = uuid.NewV4().String() // Requestid will be used to identifiy the specific message and request when retrieving a response
	req.Op = "eval"
	req.Processor = ""

	j, err := json.Marshal(req) // Formats JSON into byte format
	if err != nil {
		return
	}

	mimetype := []byte("application/json")
	mimetypelen := byte(len(mimetype))
	msg = append(msg, mimetypelen)
	msg = append(msg, mimetype...)
	msg = append(msg, j...)

	return msg, req.Requestid, nil
}

// sendMessage sends the formatted and prepared message to Gremlin Server
func (c *Client) sendMessage(msg []byte, reqid string) (response map[string]interface{}, err error) {

	c.reschan[reqid] = make(chan int) // Create channel for data arrival notification
	c.reqchan <- msg                  // Send query to write worker

	<-c.reschan[reqid]          // Wait for data to arrive
	response = c.results[reqid] // Set return value to data
	delete(c.results, reqid)    // Delete data from sorter

	return
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result.
func (c *Client) Execute(query string, bindings map[string]string) (response map[string]interface{}, err error) {

	msg, reqid, err := prepareMessage(query, bindings) // Prepare message for request
	if err != nil {
		log.Fatal(err)
	}

	response, err = c.sendMessage(msg, reqid) // Send message to Gremlin Server and retrieve response
	if err != nil {
		log.Fatal(err)
	}

	return
}

// ExecuteFile takes a file path to a Gremlin script, sends it to Gremlin Server, and returns the result.
func (c *Client) ExecuteFile(path string, bindings map[string]string) (response map[string]interface{}, err error) {

	s, err := ioutil.ReadFile(path) // Read script
	if err != nil {
		return
	}

	msg, reqid, err := prepareMessage(string(s), bindings) // Prepare message for request
	if err != nil {
		log.Fatal(err)
	}

	response, err = c.sendMessage(msg, reqid) // Send message to Gremlin Server and retrieve response
	if err != nil {
		log.Fatal(err)
	}

	return
}
