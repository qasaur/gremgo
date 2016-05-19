package gremgo

import "encoding/json"

type response struct {
	Result    interface{}            `json:"result"`
	Requestid string                 `json:"requestId"`
	Status    map[string]interface{} `json:"status"`
}

func (c *Client) responseWorker() {
	for {
		select {
		case msg := <-c.responses:
			go c.handleResponse(msg)
		default:
		}
	}
}

// handleResponse classifies the data and sends it to the appropriate function
func (c *Client) handleResponse(msg []byte) (err error) {
	var r response
	err = unmarshalResponse(msg, &r)
	if err != nil {
		return
	}
	code := r.Status["code"].(string)
	resp := determineResponse(code)
	resp.process()
	c.saveResponse(resp)
	return
}

func unmarshalResponse(msg []byte, r *response) (err error) {
	err = json.Unmarshal(msg, r) // Unwrap message
	return
}

func determineResponse(code string) (resp responder) {
	switch {
	case code == "200":
		resp = successfulResponse{}
	case code == "204":
		resp = emptyResponse{}
	case code == "206":
		resp = partialResponse{}
	default:
		resp = erroneousResponse{}
	}
	return
}

func (c *Client) saveResponse(resp responder) {
	c.mutex.Lock()
	// c.results[resp.getId()] = resp.getData() TODO: Fix this
	c.mutex.Unlock()
}

func (c *Client) retrieveResponse(id string) (data interface{}) {
	for {
		c.mutex.Lock()
		data = c.results[id]
		if data != nil {
			delete(c.results, id)
			c.mutex.Unlock()
			break
		}
		c.mutex.Unlock()
	}
	return
}

/////

type responder interface {
	process() (interface{}, string, error)
}

/////

type successfulResponse struct {
	response
}

func (res successfulResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type emptyResponse struct {
	response
}

func (res emptyResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type partialResponse struct {
	response
}

func (res partialResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

/////

type erroneousResponse struct {
	response
}

func (res erroneousResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}
