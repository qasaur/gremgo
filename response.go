package gremgo

import "encoding/json"

type response struct {
	data      interface{}
	requestid string
	code      int
}

// func (c *Client) getResponse() (msg []byte) {
// 	msg = <-c.responses
// 	return
// }

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

// handleResponse classifies the data, sorts the data, and saves it for retrieval
// func (c *Client) handleResponse(msg []byte) (err error) {
// 	var r response
// 	err = json.Unmarshal(msg, &r) // Unwrap message
// 	if err != nil {
// 		return
// 	}
// 	code := r.Status["code"]
// 	resp := determineResponse(code)
// 	resp.process()
// 	c.saveResponse(resp)
// 	return
// }
