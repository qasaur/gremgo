package gremgo

import "testing"

type dummyConnector struct {
	msg []byte
}

func (c *dummyConnector) connect() (err error) {
	return
}

func (c *dummyConnector) write(msg []byte) (err error) {
	c.msg = msg
	return
}

func (c *dummyConnector) read() (msg []byte, err error) {
	msg = c.msg
	return
}

func TestStandardRequest(t *testing.T) {
	return
}
