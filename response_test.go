package gremgo

import (
	"reflect"
	"testing"
)

var dummySuccessfulResponse = []byte(`{"result":{"data":[{"id": 2,"label": "person","type": "vertex","properties": [
  {"id": 2, "value": "vadas", "label": "name"},
  {"id": 3, "value": 27, "label": "age"}]}
  ], "meta":{}},
 "requestId":"1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
 "status":{"code":200,"attributes":{},"message":""}}`)

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

var dummySuccessfulResponseMarshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      200,
	data:      "testData",
}

var dummyPartialResponse1Marshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      206,
	data:      "testPartialData1",
}

var dummyPartialResponse2Marshalled = response{
	requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
	code:      200,
	data:      "testPartialData2",
}

func TestResponseHandling(t *testing.T) {
	c := newClient()

	c.handleResponse(dummySuccessfulResponse)

	var expected []interface{}
	expected = append(expected, dummySuccessfulResponseMarshalled.data)

	if reflect.TypeOf(expected).String() != reflect.TypeOf(c.retrieveResponse(dummySuccessfulResponseMarshalled.requestid)).String() {
		t.Error("Expected data type does not match actual.")
	}
}

func TestResponseMarshalling(t *testing.T) {
	resp, err := marshalResponse(dummySuccessfulResponse)
	if err != nil {
		t.Error(err)
	}
	if dummySuccessfulResponseMarshalled.requestid != resp.requestid || dummySuccessfulResponseMarshalled.code != resp.code {
		t.Error("Expected requestid and code does not match actual.")
	} else if reflect.TypeOf(resp.data).String() != "[]interface {}" {
		t.Error("Expected data type does not match actual.")
	}
}

func TestResponseSortingSingleResponse(t *testing.T) {
	c := newClient()

	c.sortResponse(dummySuccessfulResponseMarshalled)

	var expected []interface{}
	expected = append(expected, dummySuccessfulResponseMarshalled.data)

	if reflect.DeepEqual(c.results[dummySuccessfulResponseMarshalled.requestid], expected) != true {
		t.Fail()
	}
}

func TestResponseSortingMultipleResponse(t *testing.T) {
	c := newClient()

	c.sortResponse(dummyPartialResponse1Marshalled)
	c.sortResponse(dummyPartialResponse2Marshalled)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled.data)
	expected = append(expected, dummyPartialResponse2Marshalled.data)

	if reflect.DeepEqual(c.results[dummyPartialResponse1Marshalled.requestid], expected) != true {
		t.Fail()
	}
}

func TestResponseRetrieval(t *testing.T) {
	c := newClient()

	c.sortResponse(dummyPartialResponse1Marshalled)
	c.sortResponse(dummyPartialResponse2Marshalled)

	resp := c.retrieveResponse(dummyPartialResponse1Marshalled.requestid)

	var expected []interface{}
	expected = append(expected, dummyPartialResponse1Marshalled.data)
	expected = append(expected, dummyPartialResponse2Marshalled.data)

	if reflect.DeepEqual(resp, expected) != true {
		t.Fail()
	}
}

func TestResponseDeletion(t *testing.T) {
	c := newClient()

	c.sortResponse(dummyPartialResponse1Marshalled)
	c.sortResponse(dummyPartialResponse2Marshalled)

	c.deleteResponse(dummyPartialResponse1Marshalled.requestid)

	if len(c.results[dummyPartialResponse1Marshalled.requestid]) != 0 {
		t.Fail()
	}
}
