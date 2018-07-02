package gremgo

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	statusSuccess                  = 200
	statusNoContent                = 204
	statusPartialContent           = 206
	statusUnauthorized             = 401
	statusAuthenticate             = 407
	statusMalformedRequest         = 498
	statusInvalidRequestArguments  = 499
	statusServerError              = 500
	statusScriptEvaluationError    = 597
	statusServerTimeout            = 598
	statusServerSerializationError = 599
)

// Status struct is used to hold properties returned from requests to the gremlin server
type Status struct {
	Message    string            `json:"message"`
	Code       int               `json:"code"`
	Attributes map[string]string `json:"attributes"`
}

// Result struct is used to hold properties returned for results from requests to the gremlin server
type Result struct {
	Data []interface{}     `json:"data"`
	Meta map[string]string `json:"meta"`
}

// Response structs holds the entire response from requests to the gremlin server
type Response struct {
	RequestID string `json:"requestId"`
	Status    Status `json:"status"`
	Result    Result `json:"result"`
}

func (c *Client) handleResponse(msg []byte) (err error) {
	resp, err := marshalResponse(msg)
	if err != nil {
		log.Printf("message: %s \n err: %s", msg, err)
		return
	}

	if resp.Status.Code == statusAuthenticate { //Server request authentication
		return c.authenticate(resp.RequestID)
	}

	c.saveResponse(resp)
	return
}

// marshalResponse creates a response struct for every incoming response for further manipulation
func marshalResponse(msg []byte) (resp Response, err error) {
	err = json.Unmarshal(msg, &resp)
	if err != nil {
		return
	}
	log.Printf("msg: %s", msg)
	log.Printf("json: %+v", resp)

	err = resp.detectError()
	return
}

// saveResponse makes the response available for retrieval by the requester. Mutexes are used for thread safety.
func (c *Client) saveResponse(resp Response) {
	c.respMutex.Lock()
	var container []interface{}
	existingData, ok := c.results.Load(resp.RequestID) // Retrieve old data container (for requests with multiple responses)
	if ok {
		container = existingData.([]interface{})
	}
	newdata := append(container, resp)       // Create new data container with new data
	c.results.Store(resp.RequestID, newdata) // Add new data to buffer for future retrieval
	respNotifier, load := c.responseNotifier.LoadOrStore(resp.RequestID, make(chan int, 1))
	_ = load
	if resp.Status.Code != statusPartialContent {
		respNotifier.(chan int) <- 1
	}
	c.respMutex.Unlock()
}

// retrieveResponse retrieves the response saved by saveResponse.
func (c *Client) retrieveResponse(id string) (data []Response) {
	resp, _ := c.responseNotifier.Load(id)
	n := <-resp.(chan int)
	if n == 1 {
		if dataI, ok := c.results.Load(id); ok {
			d := dataI.([]interface{})
			data = make([]Response, len(d))
			for i := range d {
				data[i] = d[i].(Response)
			}
			close(resp.(chan int))
			c.responseNotifier.Delete(id)
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
func (r *Response) detectError() (err error) {
	switch r.Status.Code {
	case statusSuccess, statusNoContent, statusPartialContent:
		break
	case statusUnauthorized:
		err = fmt.Errorf("UNAUTHORIZED - Response: %+v", r)
	case statusAuthenticate:
		err = fmt.Errorf("AUTHENTICATE - Response: %+v", r)
	case statusMalformedRequest:
		err = fmt.Errorf("MALFORMED REQUEST - Response: %+v", r)
	case statusInvalidRequestArguments:
		err = fmt.Errorf("INVALID REQUEST ARGUMENTS - Response: %+v", r)
	case statusServerError:
		err = fmt.Errorf("SERVER ERROR - Response: %+v", r)
	case statusScriptEvaluationError:
		err = fmt.Errorf("SCRIPT EVALUATION ERROR - Response: %+v", r)
	case statusServerTimeout:
		err = fmt.Errorf("SERVER TIMEOUT - %+v", r)
	case statusServerSerializationError:
		err = fmt.Errorf("SERVER SERIALIZATION ERROR - Response: %+v", r)
	default:
		err = fmt.Errorf("UNKNOWN ERROR - Response: %+v", r)
	}
	return
}
