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

// func TestResponseGet(t *testing.T) {
// 	c := Client{}
// 	c.responses = make(chan []byte, 2)
// 	c.responses <- dummySuccessfulResponse
// 	c.responses <- dummyPartialResponse1
// 	if reflect.DeepEqual(resp, dummySuccessfulResponse) != true {
// 		t.Fail()
// 	}
// }

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

// func TestResponseSortingMultipleResponse(t *testing.T) {
// 	c := newClient()
//
// 	c.sortResponse(dummyPartialResponse1Marshalled)
// 	c.sortResponse(dummyPartialResponse2Marshalled)
//
// 	// if reflect.DeepEqual(x interface{}, y interface{})
// }
