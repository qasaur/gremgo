package gremgo

import "encoding/json"

type response struct {
	Result    interface{}            `json:"result"`
	Requestid string                 `json:"requestId"`
	Status    map[string]interface{} `json:"status"`
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
	c.results[resp.getID()] = resp.getData() // TODO: Fix this
	c.mutex.Unlock()
}

func (c *Client) retrieveResponse(req requester) (data interface{}) {
	reqID := req.getID()
	for {
		c.mutex.Lock()
		data = c.results[reqID]
		if data != nil {
			delete(c.results, reqID)
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
	getID() string
	getData() responseData
}

/////

type responseData struct {
}

/////

type successfulResponse struct {
	response
}

func (res successfulResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

func (res successfulResponse) getID() (id string) {
	return
}

func (res successfulResponse) getData() (data responseData) {
	return
}

/////

type emptyResponse struct {
	response
}

func (res emptyResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

func (res emptyResponse) getID() (id string) {
	return
}

func (res emptyResponse) getData() (data responseData) {
	return
}

/////

type partialResponse struct {
	response
}

func (res partialResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

func (res partialResponse) getID() (id string) {
	return
}

func (res partialResponse) getData() (data responseData) {
	return
}

/////

type erroneousResponse struct {
	response
}

func (res erroneousResponse) process() (data interface{}, id string, err error) {
	return res.Result, res.Requestid, nil
}

func (res erroneousResponse) getID() (id string) {
	return
}

func (res erroneousResponse) getData() (data responseData) {
	return
}
