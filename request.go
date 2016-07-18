package gremgo

import (
	"encoding/json"
	"strconv"
	"sync"
)

/////

type requester interface {
	prepare() error
	getID() string
	getRequest() request
}

/////

// request is a container for all evaluation request parameters to be sent to the Gremlin Server.
type request struct {
	Requestid string                 `json:"requestId"`
	Op        string                 `json:"op"`
	Processor string                 `json:"processor"`
	Args      map[string]interface{} `json:"args"`
}

/////

type UidGenerator struct {
	uid int
	mu  sync.Mutex
}

func (u *UidGenerator) Next() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.uid++
	return u.uid
}

var uidGen = &UidGenerator{}

// prepareRequest packages a query and binding into the format that Gremlin Server accepts
func prepareRequest(query string, bindings map[string]string) (req request, id string) {
	id = strconv.Itoa(uidGen.Next())

	req.Requestid = id
	req.Op = "eval"
	req.Processor = ""

	req.Args = make(map[string]interface{})
	req.Args["language"] = "gremlin-groovy"
	req.Args["gremlin"] = query
	req.Args["bindings"] = bindings

	return
}

/////

// formatMessage takes a request type and formats it into being able to be delivered to Gremlin Server
func packageRequest(req request) (msg []byte, err error) {

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

// dispactchRequest sends the request for writing to the remote Gremlin Server
func (c *Client) dispatchRequest(msg []byte) {
	c.requests <- msg
}

/////
