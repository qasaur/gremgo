package gremgo

import (
	"encoding/json"
	"errors"
)

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
	c.saveResponse(resp)
	return
}

// marshalResponse creates a response struct for every incoming response for further manipulation
func marshalResponse(msg []byte) (resp response, err error) {
	var j map[string]interface{}
	err = json.Unmarshal(msg, &j)
	if err != nil {
		return
	}

	status := j["status"].(map[string]interface{})
	result := j["result"].(map[string]interface{})
	code := status["code"].(float64)

	resp.code = int(code)
	err = responseDetectError(resp.code)
	if err != nil {
		resp.data = err // Modify response vehicle to have error (if exists) as data
	} else {
		resp.data = result["data"]
	}
	err = nil
	resp.requestid = j["requestId"].(string)
	return
}

// saveResponse makes the response available for retrieval by the requester. Mutexes are used for thread safety.
func (c *Client) saveResponse(resp response) {
	c.respMutex.Lock()
	container := c.results[resp.requestid]  // Retrieve old data container (for requests with multiple responses)
	newdata := append(container, resp.data) // Create new data container with new data
	c.results[resp.requestid] = newdata     // Add new data to buffer for future retrieval
	if resp.code == 200 {
		if c.responseNotifyer[resp.requestid] == nil {
			c.responseNotifyer[resp.requestid] = make(chan int, 1)
		}
		c.responseNotifyer[resp.requestid] <- 1
	}
	c.respMutex.Unlock()
	return
}

// retrieveResponse retrieves the response saved by saveResponse.
func (c *Client) retrieveResponse(id string) (data []interface{}) {
	n := <-c.responseNotifyer[id]
	if n == 1 {
		data = c.results[id]
		close(c.responseNotifyer[id])
		delete(c.responseNotifyer, id)
		c.deleteResponse(id)
	}
	return
}

// deleteRespones deletes the response from the container. Used for cleanup purposes by requester.
func (c *Client) deleteResponse(id string) {
	delete(c.results, id)
	return
}

// responseDetectError detects any possible errors in responses from Gremlin Server and generates an error for each code
func responseDetectError(code int) (err error) {
	switch {
	case code == 200:
		break
	case code == 204:
		break
	case code == 206:
		break
	case code == 401:
		err = errors.New("UNAUTHORIZED")
	case code == 407:
		err = errors.New("AUTHENTICATE")
	case code == 498:
		err = errors.New("MALFORMED REQUEST")
	case code == 499:
		err = errors.New("INVALID REQUEST ARGUMENTS")
	case code == 500:
		err = errors.New("SERVER ERROR")
	case code == 597:
		err = errors.New("SCRIPT EVALUATION ERROR")
	case code == 598:
		err = errors.New("SERVER TIMEOUT")
	case code == 599:
		err = errors.New("SERVER SERIALIZATION ERROR")
	default:
		err = errors.New("UNKNOWN ERROR")
	}
	return
}
