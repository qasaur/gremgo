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
	var container []interface{}
	existingData, ok := c.results.Load(resp.requestid) // Retrieve old data container (for requests with multiple responses)
	if ok {
		container = existingData.([]interface{})
	}
	newdata := append(container, resp.data)  // Create new data container with new data
	c.results.Store(resp.requestid, newdata) // Add new data to buffer for future retrieval
	respNotifier, _ := c.responseNotifyer.LoadOrStore(resp.requestid, make(chan int, 1))
	if resp.code != 206 {
		respNotifier.(chan int) <- 1
	}
	c.respMutex.Unlock()
}

// retrieveResponse retrieves the response saved by saveResponse.
func (c *Client) retrieveResponse(id string) (data []interface{}) {
	resp, _ := c.responseNotifyer.Load(id)
	n := <-resp.(chan int)
	if n == 1 {
		if dataI, ok := c.results.Load(id); ok {
			data = dataI.([]interface{})
			close(resp.(chan int))
			c.responseNotifyer.Delete(id)
			c.deleteResponse(id)
		}
	}
	return
}

// deleteRespones deletes the response from the container. Used for cleanup purposes by requester.
func (c *Client) deleteResponse(id string) {
	c.results.Delete(id)
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
