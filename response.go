package gremgo

import (
	"encoding/json"
	"log"
)

type response struct {
	Result    interface{}            `json:"result"`
	Requestid string                 `json:"requestId"`
	Status    map[string]interface{} `json:"status"`
}

// handleResponse classifies the data and sends it to the appropriate function
func (c *Client) handleResponse(msg []byte) (err error) {
	var r response
	err = json.Unmarshal(msg, &r) // Unwrap message
	if err != nil {
		log.Fatal(err)
	}
	code := r.Status["code"]
	switch {
	case code == "200":
		resp := successfulResponse{response: r}
		c.processResponse(&resp)
	case code == "204":
		resp := emptyResponse{response: r}
		c.processResponse(&resp)
	case code == "206":
		resp := partialResponse{response: r}
		c.processResponse(&resp)
	default:
		resp := erroneousResponse{response: r}
		c.processResponse(&resp)
	}
	return
}

func (c *Client) processResponse(resp responseHandler) {
	data, id, err := resp.process()
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Fix processing of partial requests
	c.mutex.Lock()
	c.responses[id] = data
	c.mutex.Unlock()
	return
}

func (c *Client) retrieveResponse(id string) (data interface{}) {
	for {
		c.mutex.Lock()
		data = c.responses[id]
		if data != nil {
			delete(c.responses, id)
			c.mutex.Unlock()
			break
		}
		c.mutex.Unlock()
	}
	return
}

/////

type responseHandler interface {
	process() (interface{}, string, error)
}

/////

type successfulResponse struct {
	response
}

func (res *successfulResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type emptyResponse struct {
	response
}

func (res *emptyResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type partialResponse struct {
	response
}

func (res *partialResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type erroneousResponse struct {
	response
}

func (res *erroneousResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}
