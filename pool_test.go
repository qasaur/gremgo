package gremgo

import (
	"testing"
	"time"
)

func TestPurge(t *testing.T) {
	n := time.Now()

	// invalid has timedout and should be cleaned up
	invalid := &idleConnection{t: n.Add(-30 * time.Second), pc: &PooledConnection{Client: &Client{}}}
	// valid has not yet timed out and should remain in the idle pool
	valid := &idleConnection{t: n.Add(30 * time.Second), pc: &PooledConnection{Client: &Client{}}}

	// Pool has a 30 second timeout and an idle connection slice containing both
	// the invalid and valid idle connections
	p := &Pool{IdleTimeout: time.Second * 30, idle: []*idleConnection{invalid, valid}}

	if len(p.idle) != 2 {
		t.Errorf("Expected 2 idle connections, got %d", len(p.idle))
	}

	p.purge()

	if len(p.idle) != 1 {
		t.Errorf("Expected 1 idle connection after purge, got %d", len(p.idle))
	}

	if p.idle[0].t != valid.t {
		t.Error("Expected the valid connection to remain in idle pool")
	}

}

func TestPurgeErrorClosedConnection(t *testing.T) {
	n := time.Now()

	p := &Pool{IdleTimeout: time.Second * 30}

	valid := &idleConnection{t: n.Add(30 * time.Second), pc: &PooledConnection{Client: &Client{}}}

	client := &Client{}

	closed := &idleConnection{t: n.Add(30 * time.Second), pc: &PooledConnection{Pool: p, Client: client}}

	idle := []*idleConnection{valid, closed}

	p.idle = idle

	// Simulate error
	closed.pc.Client.Errored = true

	if len(p.idle) != 2 {
		t.Errorf("Expected 2 idle connections, got %d", len(p.idle))
	}

	p.purge()

	if len(p.idle) != 1 {
		t.Errorf("Expected 1 idle connection after purge, got %d", len(p.idle))
	}

	if p.idle[0] != valid {
		t.Error("Expected valid connection to remain in pool")
	}
}

func TestPooledConnectionClose(t *testing.T) {
	pool := &Pool{}
	pc := &PooledConnection{Pool: pool}

	if len(pool.idle) != 0 {
		t.Errorf("Expected 0 idle connection, got %d", len(pool.idle))
	}

	pc.Close()

	if len(pool.idle) != 1 {
		t.Errorf("Expected 1 idle connection, got %d", len(pool.idle))
	}

	idled := pool.idle[0]

	if idled == nil {
		t.Error("Expected to get connection")
	}

	if idled.t.IsZero() {
		t.Error("Expected an idled time")
	}
}

func TestFirst(t *testing.T) {
	n := time.Now()
	pool := &Pool{MaxActive: 1, IdleTimeout: 30 * time.Second}
	idled := []*idleConnection{
		&idleConnection{t: n.Add(-45 * time.Second), pc: &PooledConnection{Pool: pool, Client: &Client{}}}, // expired
		&idleConnection{t: n.Add(-45 * time.Second), pc: &PooledConnection{Pool: pool, Client: &Client{}}}, // expired
		&idleConnection{pc: &PooledConnection{Pool: pool, Client: &Client{}}},                              // valid
	}
	pool.idle = idled

	if len(pool.idle) != 3 {
		t.Errorf("Expected 3 idle connection, got %d", len(pool.idle))
	}

	// Get should return the last idle connection and purge the others
	c := pool.first()

	if c != pool.idle[0] {
		t.Error("Expected to get first connection in idle slice")
	}

	// Empty pool should return nil
	emptyPool := &Pool{}

	c = emptyPool.first()

	if c != nil {
		t.Errorf("Expected nil, got %T", c)
	}
}

func TestGetAndDial(t *testing.T) {
	n := time.Now()

	pool := &Pool{IdleTimeout: time.Second * 30}

	invalid := &idleConnection{t: n.Add(-30 * time.Second), pc: &PooledConnection{Pool: pool, Client: &Client{}}}

	idle := []*idleConnection{invalid}
	pool.idle = idle

	client := &Client{}
	pool.Dial = func() (*Client, error) {
		return client, nil
	}

	if len(pool.idle) != 1 {
		t.Error("Expected 1 idle connection")
	}

	if pool.idle[0] != invalid {
		t.Error("Expected invalid connection")
	}

	conn, err := pool.Get()

	if err != nil {
		t.Error(err)
	}

	if len(pool.idle) != 0 {
		t.Errorf("Expected 0 idle connections, got %d", len(pool.idle))
	}

	if conn.Client != client {
		t.Error("Expected correct client to be returned")
	}

	if pool.active != 1 {
		t.Errorf("Expected 1 active connection, got %d", pool.active)
	}

	// Close the connection and ensure it was returned to the idle pool
	conn.Close()

	if len(pool.idle) != 1 {
		t.Error("Expected connection to be returned to idle pool")
	}

	if pool.active != 0 {
		t.Errorf("Expected 0 active connections, got %d", pool.active)
	}

	// Get a new connection and ensure that it is the now idling connection
	conn, err = pool.Get()

	if err != nil {
		t.Error(err)
	}

	if conn.Client != client {
		t.Error("Expected the same connection to be reused")
	}

	if pool.active != 1 {
		t.Errorf("Expected 1 active connection, got %d", pool.active)
	}
}
