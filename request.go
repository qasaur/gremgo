package gremgo

import (
	"encoding/json"
	"log"

	"github.com/satori/go.uuid"
)

// Request is a container for all request parameters to be sent to the Gremlin Server
type Request struct {
	Requestid string  `json:"requestId"`
	Op        string  `json:"op"`
	Processor string  `json:"processor"`
	Args      ReqArgs `json:"args"`
}

// ReqArgs define the arguments for the Gremlin request
type ReqArgs struct {
	Gremlin  string            `json:"gremlin"`
	Language string            `json:"language"`
	Bindings map[string]string `json:"bindings"`
}

func prepareMessage(j []byte) (msg []byte) {
	mimetype := []byte("application/json")
	mimetypelen := byte(len(mimetype))
	msg = append(msg, mimetypelen)
	msg = append(msg, mimetype...)
	msg = append(msg, j...)
	return
}

// Execute formats a raw Gremlin query, sends it to Gremlin Server, and returns the result
func (c *Client) Execute(query string, bindings map[string]string) (r map[string]interface{}, err error) {
	var req Request
	var args ReqArgs

	args.Gremlin = query
	args.Language = "gremlin-groovy"
	args.Bindings = bindings

	req.Args = args
	req.Requestid = uuid.NewV4().String()
	req.Op = "eval"
	req.Processor = ""

	j, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	msg := prepareMessage(j)

	c.reschan[req.Requestid] = make(chan int) // Create channel for data arrival notification
	c.reqchan <- msg                          // Send query to write worker

	<-c.reschan[req.Requestid]       // Wait for data to arrive
	r = c.results[req.Requestid]     // Set return value to data
	delete(c.results, req.Requestid) // Delete data from sorter

	return
}
