package gremgo

import (
	"log"
	"reflect"
	"testing"
)

/*
Dummy responses for mocking
*/

var dummySuccessfulResponse = []byte(`{"result":{"data":[{"id": 2,"label": "person","type": "vertex","properties": [
  {"id": 2, "value": "vadas", "label": "name"},
  {"id": 3, "value": 27, "label": "age"}]}
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":200,"attributes":{},"message":""}}`)

var dummyNeedAuthenticationResponse = []byte(`{"result":{},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":407,"attributes":{},"message":""}}`)

var dummyPartialResponse1 = []byte(`{"result":{"data":[{"id": 2,"label": "person","type": "vertex","properties": [
  {"id": 2, "value": "vadas", "label": "name"},
  {"id": 3, "value": 27, "label": "age"}]},
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":206,"attributes":{},"message":""}}`)

var dummyPartialResponse2 = []byte(`{"result":{"data":[{"id": 4,"label": "person","type": "vertex","properties": [
  {"id": 5, "value": "quant", "label": "name"},
  {"id": 6, "value": 54, "label": "age"}]},
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":200,"attributes":{},"message":""}}`)

var dummySuccessfulResponseMarshalled = Response{
	RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	Status:    Status{Code: 200},
	Result:    Result{Data: []interface{}{"testData"}},
}

var dummyNeedAuthenticationResponseMarshalled = Response{
	RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	Status:    Status{Code: 407},
	Result:    Result{Data: []interface{}{""}},
}

var dummyPartialResponse1Marshalled = Response{
	RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	Status:    Status{Code: 206}, // Code 206 indicates that the response is not the terminating response in a sequence of responses
	Result:    Result{Data: []interface{}{"testPartialData1"}},
}

var dummyPartialResponse2Marshalled = Response{
	RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	Status:    Status{Code: 200},
	Result:    Result{Data: []interface{}{"testPartialData2"}},
}

// TestResponseHandling tests the overall response handling mechanism of gremgo
func TestResponseHandling(t *testing.T) {
	c := newClient()

	c.handleResponse(dummySuccessfulResponse)

	var expected []Response
	expected = append(expected, dummySuccessfulResponseMarshalled)

	if reflect.TypeOf(expected).String() != reflect.TypeOf(c.retrieveResponse(dummySuccessfulResponseMarshalled.RequestID)).String() {
		t.Error("Expected data type does not match actual.")
	}
}

func TestResponseAuthHandling(t *testing.T) {
	c := newClient()
	ws := new(Ws)
	ws.auth = &auth{username: "test", password: "test"}
	c.conn = ws

	c.handleResponse(dummyNeedAuthenticationResponse)

	req, err := prepareAuthRequest(dummyNeedAuthenticationResponseMarshalled.RequestID, "test", "test")
	if err != nil {
		return
	}

	sampleAuthRequest, err := packageRequest(req)
	if err != nil {
		log.Println(err)
		return
	}

	authRequest := <-c.requests //Simulate that client send auth challenge to server

	if !reflect.DeepEqual(authRequest, sampleAuthRequest) {
		t.Error("Expected data type does not match actual.")
	}

	c.handleResponse(dummySuccessfulResponse) //If authentication is successful the server returns the origin petition

	var expectedSuccessful []Response
	expectedSuccessful = append(expectedSuccessful, dummySuccessfulResponseMarshalled)

	if reflect.TypeOf(expectedSuccessful).String() != reflect.TypeOf(c.retrieveResponse(dummySuccessfulResponseMarshalled.RequestID)).String() {
		t.Error("Expected data type does not match actual.")
	}
}

// TestResponseMarshalling tests the ability to marshal a response into a designated response struct for further manipulation
func TestResponseMarshalling(t *testing.T) {
	resp, err := marshalResponse(dummySuccessfulResponse)
	if err != nil {
		t.Error(err)
	}
	if dummySuccessfulResponseMarshalled.RequestID != resp.RequestID || dummySuccessfulResponseMarshalled.Status.Code != resp.Status.Code {
		t.Error("Expected requestId and code does not match actual.")
	} else if reflect.TypeOf(resp.Result.Data).String() != "[]interface {}" {
		t.Error("Expected data type does not match actual.")
	}
}

// TestResponseSortingSingleResponse tests the ability for sortResponse to save a response received from Gremlin Server
func TestResponseSortingSingleResponse(t *testing.T) {

	c := newClient()

	c.saveResponse(dummySuccessfulResponseMarshalled)

	var expected []interface{}
	expected = append(expected, dummySuccessfulResponseMarshalled)

	result, _ := c.results.Load(dummySuccessfulResponseMarshalled.RequestID)
	if reflect.DeepEqual(result.([]interface{}), expected) != true {
		t.Fail()
	}
}

// TestResponseSortingMultipleResponse tests the ability for the sortResponse function to categorize and group responses that are sent in a stream
func TestResponseSortingMultipleResponse(t *testing.T) {

	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled)
	expected = append(expected, dummyPartialResponse2Marshalled)

	results, _ := c.results.Load(dummyPartialResponse1Marshalled.RequestID)
	if reflect.DeepEqual(results.([]interface{}), expected) != true {
		t.Fail()
	}
}

// TestResponseRetrieval tests the ability for a requester to retrieve the response for a specified requestId generated when sending the request
func TestResponseRetrieval(t *testing.T) {
	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	resp := c.retrieveResponse(dummyPartialResponse1Marshalled.RequestID)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled)
	expected = append(expected, dummyPartialResponse2Marshalled)

	if reflect.DeepEqual(resp, expected) != true {
		t.Fail()
	}
}

// TestResponseDeletion tests the ability for a requester to clean up after retrieving a response after delivery to a client
func TestResponseDeletion(t *testing.T) {
	c := newClient()

	c.saveResponse(dummyPartialResponse1Marshalled)
	c.saveResponse(dummyPartialResponse2Marshalled)

	c.deleteResponse(dummyPartialResponse1Marshalled.RequestID)

	if _, ok := c.results.Load(dummyPartialResponse1Marshalled.RequestID); ok {
		t.Fail()
	}
}

var codes = []struct {
	code int
}{
	{200},
	{204},
	{206},
	{401},
	{407},
	{498},
	{499},
	{500},
	{597},
	{598},
	{599},
	{3434}, // Testing unknown error code
}

// Tests detection of errors and if an error is generated for a specific error code
func TestResponseErrorDetection(t *testing.T) {
	for _, co := range codes {
		dummyResponse := Response{
			RequestID: "",
			Status:    Status{Code: co.code},
			Result:    Result{},
		}
		err := dummyResponse.detectError()
		switch {
		case co.code == 200:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		case co.code == 204:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		case co.code == 206:
			if err != nil {
				t.Log("Successful response returned error.")
			}
		default:
			if err == nil {
				t.Log("Unsuccessful response did not return error.")
			}
		}
	}
}
