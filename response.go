package gremgo

import "encoding/json"

type response struct {
	data      interface{}
	requestid string
	code      int
}

func (c *Client) handleResponse(msg []byte) (err error) {
	resp, err := marshalResponse(msg)
	if err != nil {
		return
	}
	c.sortResponse(resp)
	return
}

func marshalResponse(msg []byte) (resp response, err error) {
	var j map[string]interface{}
	err = json.Unmarshal(msg, &j)

	status := j["status"].(map[string]interface{})
	result := j["result"].(map[string]interface{})
	code := status["code"].(float64)

	resp.code = int(code)
	resp.data = result["data"]
	resp.requestid = j["requestId"].(string)
	return
}

func (c *Client) sortResponse(resp response) {
	c.respMutex.RLock()
	container := c.results[resp.requestid]
	c.respMutex.RUnlock()
	data := append(container, resp.data)
	c.respMutex.Lock()
	c.results[resp.requestid] = data
	c.respMutex.Unlock()
	return
}

func (c *Client) retrieveResponse(id string) (data []interface{}) {
	data = c.results[id]
	return
}

func (c *Client) deleteResponse(id string) {
	delete(c.results, id)
	return
}
