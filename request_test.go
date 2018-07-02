package gremgo

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestRequestPreparation tests the ability to package a query and a set of bindings into a request struct for further manipulation
func TestRequestPreparation(t *testing.T) {
	query := "g.V(x)"
	bindings := map[string]string{"x": "10"}
	rebindings := map[string]string{}
	req, id, err := prepareRequestWithBindings(query, bindings, rebindings)
	if err != nil {
		t.Error(err)
	}

	expectedRequest := request{
		RequestID: id,
		Op:        "eval",
		Processor: "",
		Args: map[string]interface{}{
			"gremlin":    query,
			"bindings":   bindings,
			"language":   "gremlin-groovy",
			"rebindings": rebindings,
		},
	}

	if reflect.DeepEqual(req, expectedRequest) != true {
		t.Fail()
	}
}

// TestRequestPackaging tests the ability for gremgo to format a request using the established Gremlin Server WebSockets protocol for delivery to the server
func TestRequestPackaging(t *testing.T) {
	testRequest := request{
		RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
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

	mimetype := []byte("application/vnd.gremlin-v2.0+json")
	mimetypelen := byte(len(mimetype))

	expected = append(expected, mimetypelen)
	expected = append(expected, mimetype...)
	expected = append(expected, j...)

	if reflect.DeepEqual(msg, expected) != true {
		t.Fail()
	}
}

// TestRequestDispatch tests the ability for a requester to send a request to the client for writing to Gremlin Server
func TestRequestDispatch(t *testing.T) {
	testRequest := request{
		RequestID: "1d6d02bd-8e56-421d-9438-3bd6d0079ff1",
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
	req := <-c.requests // c.requests is the channel where all requests are sent for writing to Gremlin Server, write workers listen on this channel
	if reflect.DeepEqual(msg, req) != true {
		t.Fail()
	}
}

// TestAuthRequestDispatch tests the ability for a requester to send a request to the client for writing to Gremlin Server
func TestAuthRequestDispatch(t *testing.T) {
	id := "1d6d02bd-8e56-421d-9438-3bd6d0079ff1"
	testRequest, err := prepareAuthRequest(id, "test", "root")

	c := newClient()
	msg, err := packageRequest(testRequest)
	if err != nil {
		t.Error(err)
	}
	c.dispatchRequest(msg)
	req := <-c.requests // c.requests is the channel where all requests are sent for writing to Gremlin Server, write workers listen on this channel
	if reflect.DeepEqual(msg, req) != true {
		t.Fail()
	}
}

// TestAuthRequestPreparation tests the ability to create successful authentication request
func TestAuthRequestPreparation(t *testing.T) {
	id := "1d6d02bd-8e56-421d-9438-3bd6d0079ff1"
	testRequest, err := prepareAuthRequest(id, "test", "root")
	if err != nil {
		t.Fail()
	}
	if testRequest.RequestID != id || testRequest.Processor != "trasversal" || testRequest.Op != "authentication" {
		t.Fail()
	}
	if len(testRequest.Args) != 1 || testRequest.Args["sasl"] == "" {
		t.Fail()
	}
}
