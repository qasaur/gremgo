package gremgo

import "time"

//DialerConfig is the struct for defining configuration for WebSocket dialer
type DialerConfig func(*Ws)

//SetAuthentication sets on dialer credentials for authentication
func SetAuthentication(username string, password string) DialerConfig {
	return func(c *Ws) {
		c.auth = &auth{username:username, password:password}
	}
}

//SetTimeout sets the dial timeout
func SetTimeout(seconds int) DialerConfig {
	return func(c *Ws) {
		c.timeout = time.Duration(seconds) * time.Second
	}
}

//SetPingInterval sets the interval of ping sending for know is
//connection is alive and in consequence the client is connected
func SetPingInterval(seconds int) DialerConfig {
	return func(c *Ws) {
		c.pingInterval = time.Duration(seconds)* time.Second
	}
}

//SetWritingWait sets the time for waiting that writing occur
func SetWritingWait(seconds int) DialerConfig {
	return func(c *Ws) {
		c.writingWait = time.Duration(seconds)* time.Second
	}
}

//SetReadingWait sets the time for waiting that reading occur
func SetReadingWait(seconds int) DialerConfig {
	return func(c *Ws) {
		c.readingWait = time.Duration(seconds)* time.Second
	}
}
