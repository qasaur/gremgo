package gremgo

import "testing"

func TestPanicOnMissingAuthCredentials(t *testing.T) {
	c := newClient()
	ws := new(Ws)
	c.conn = ws

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	c.conn.getAuth()
}