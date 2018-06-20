package gremgo

import (
	"encoding/json"

	"github.com/satori/go.uuid"
	"encoding/base64"
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
	RequestId string                 `json:"requestId"`
	Op        string                 `json:"op"`
	Processor string                 `json:"processor"`
	Args      map[string]interface{} `json:"args"`
}

/////

// prepareRequest packages a query and binding into the format that Gremlin Server accepts
func prepareRequest(query string, bindings, rebindings map[string]string) (req request, id string, err error) {
	var uuID uuid.UUID
        uuID, err = uuid.NewV4()
        if err != nil {
                return
        }
        id = uuID.String()

	req.RequestId = id
	req.Op = "eval"
	req.Processor = ""

	req.Args = make(map[string]interface{})
	req.Args["language"] = "gremlin-groovy"
	req.Args["gremlin"] = query
	req.Args["bindings"] = bindings
	req.Args["rebindings"] = rebindings

	return
}

//prepareAuthRequest creates a ws request for Gremlin Server
func prepareAuthRequest(requestId string, username string, password string) (req request, err error) {
	req.RequestId = requestId
	req.Op = "authentication"
	req.Processor = "trasversal"

	var simpleAuth []byte
	user := []byte(username)
	pass := []byte(password)

	simpleAuth = append(simpleAuth, 0)
	simpleAuth = append(simpleAuth, user...)
	simpleAuth = append(simpleAuth, 0)
	simpleAuth = append(simpleAuth, pass...)

	req.Args = make(map[string]interface{})
	req.Args["sasl"] = base64.StdEncoding.EncodeToString(simpleAuth)

	return
}

/////

// formatMessage takes a request type and formats it into being able to be delivered to Gremlin Server
func packageRequest(req request) (msg []byte, err error) {

	j, err := json.Marshal(req) // Formats request into byte format
	if err != nil {
		return
	}
	mimeType := []byte("application/vnd.gremlin-v2.0+json")
	msg = append([]byte{0x21}, mimeType...) //0x21 is the fixed length of mimeType in hex
	msg = append(msg, j...)

	return
}

/////

// dispactchRequest sends the request for writing to the remote Gremlin Server
func (c *Client) dispatchRequest(msg []byte) {
	c.requests <- msg
}

/////
