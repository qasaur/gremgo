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
	c.respMutex.RLock()                     // Lock for reading
	container := c.results[resp.requestid]  // Retrieve old data container (for requests with multiple responses)
	c.respMutex.RUnlock()                   // Unlock for reading
	newdata := append(container, resp.data) // Create new data container with new data
	c.respMutex.Lock()                      // Lock for writing
	c.results[resp.requestid] = newdata     // Add new data to buffer for future retrieval
	c.responseNotifyer[resp.requestid] = 1
	c.respMutex.Unlock() // Unlock for writing
	return
}

// retrieveResponse retrieves the response saved by saveResponse.
func (c *Client) retrieveResponse(id string) (data []interface{}) {
	var recieved bool
	recieved = false
	for recieved == false {
		c.respMutex.RLock()
		if c.responseNotifyer[id] != 0 {
			recieved = true
		}
		c.respMutex.RUnlock()
	}
	data = c.results[id]
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
