package gremgo

// type dummyConnector struct {
// 	expected response
// 	msg      []byte
// }
//
// func (c *dummyConnector) connect() (err error) {
// 	return
// }
//
// func (c *dummyConnector) write(msg []byte) (err error) {
// 	c.msg = msg
// 	return
// }
//
// func (c *dummyConnector) read() (msg []byte, err error) {
// 	// dummyID := "1d6d02bd-8e56-421d-9438-3bd6d0079ff1"
// 	c.msg, err = json.Marshal(c.expected)
// 	return
// }
//
// func TestStandardRequest(t *testing.T) {
// 	dialer := dummyConnector{expected: dummySuccessfulResponse}
// 	c, err := Dial(&dialer)
// 	res, err := c.Execute("g.V(x)", map[string]string{"x": "10"})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if res != "success" {
// 		t.Fail()
// 	}
// }
