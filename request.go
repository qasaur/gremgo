package gremgo

import (
	"encoding/json"

	"github.com/satori/go.uuid"
)

// request is a container for all evaluation request parameters to be sent to the Gremlin Server.
type request struct {
	Requestid string                 `json:"requestId"`
	Op        string                 `json:"op"`
	Processor string                 `json:"processor"`
	Args      map[string]interface{} `json:"args"`
}

// formatMessage takes a request type and formats it into being able to be delivered to Gremlin Server
func formatRequest(req request) (msg []byte, err error) {

	j, err := json.Marshal(req) // Formats request into byte format
	if err != nil {
		return
	}

	mimetype := []byte("application/json")
	mimetypelen := byte(len(mimetype))

	msg = append(msg, mimetypelen)
	msg = append(msg, mimetype...)
	msg = append(msg, j...)

	return
}

/////

type requester interface {
	prepare() error
	format() error
	getID() string
}

/////

type evalRequest struct {
	request
	query    string
	bindings map[string]string
	prepared []byte
}

func (req *evalRequest) prepare() (err error) {
	req.request.Requestid = uuid.NewV4().String() // Requestid will be used to identify the specific message and request when retrieving a response
	req.request.Op = "eval"
	req.request.Processor = ""

	req.request.Args = make(map[string]interface{})

	req.request.Args["gremlin"] = req.query
	req.request.Args["language"] = "gremlin-groovy"
	req.request.Args["bindings"] = req.bindings

	return
}

func (req *evalRequest) format() (err error) {
	req.prepared, err = formatRequest(req.request)
	return
}

func (req *evalRequest) getID() (id string) {
	return req.Requestid
}

/////

func (c *Client) sendRequest(msg []byte) {
	c.requests <- msg // Send query to write worker
}

/////

func (c *Client) createRequest(req requester) (err error) {
	err = req.prepare()
	if err != nil {
		return
	}
	err = req.format()
	if err != nil {
		return
	}
	return
}
