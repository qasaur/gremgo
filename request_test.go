package gremgo

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRequestPreparation(t *testing.T) {
	query := "g.V(x)"
	bindings := map[string]string{"x": "10"}
	req, id := prepareRequest(query, bindings)

	expectedRequest := request{
		Requestid: id,
		Op:        "eval",
		Processor: "",
		Args: map[string]interface{}{
			"gremlin":  "g.V(x)",
			"bindings": map[string]string{"x": "10"},
			"language": "gremlin-groovy",
		},
	}

	if reflect.DeepEqual(req, expectedRequest) != true {
		t.Fail()
	}
}

func TestRequestPackaging(t *testing.T) {
	testRequest := request{
		Requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
		Op:        "eval",
		Processor: "",
		Args: map[string]interface{}{
			"gremlin":  "g.V(x)",
			"bindings": map[string]string{"x": "10"},
			"language": "gremlin-groovy",
		},
	}

	msg, err := packageRequest(testRequest)
	if err != nil {
		t.Error(err)
	}

	j, err := json.Marshal(testRequest)
	if err != nil {
		t.Error(err)
	}

	var expected []byte

	mimetype := []byte("application/json")
	mimetypelen := byte(len(mimetype))

	expected = append(expected, mimetypelen)
	expected = append(expected, mimetype...)
	expected = append(expected, j...)

	if reflect.DeepEqual(msg, expected) != true {
		t.Fail()
	}
}

func TestRequestDispatch(t *testing.T) {
	testRequest := request{
		Requestid: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
		Op:        "eval",
		Processor: "",
		Args: map[string]interface{}{
			"gremlin":  "g.V(x)",
			"bindings": map[string]string{"x": "10"},
			"language": "gremlin-groovy",
		},
	}
	c := newClient()
	msg, err := packageRequest(testRequest)
	if err != nil {
		t.Error(err)
	}
	c.dispatchRequest(msg)
	req := <-c.requests
	if reflect.DeepEqual(msg, req) != true {
		t.Fail()
	}
}
