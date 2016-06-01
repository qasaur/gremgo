package gremgo

import (
	"reflect"
	"testing"
)

func TestNewDialer(t *testing.T) {
	dialer := NewDialer("127.0.0.1")
	expected := Ws{host: "127.0.0.1"}
	if reflect.DeepEqual(dialer, expected) != true {
		t.Fail()
	}
}
