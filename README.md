# gremgo-neptune

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/schwartzmx/gremgo-neptune) [![Build Status](https://travis-ci.org/schwartzmx/gremgo-neptune.svg?branch=master)](https://travis-ci.org/schwartzmx/gremgo-neptune) [![Go Report Card](https://goreportcard.com/badge/github.com/schwartzmx/gremgo-neptune)](https://goreportcard.com/report/github.com/schwartzmx/gremgo-neptune)

gremgo-neptune is a fork of [qasaur/gremgo](https://github.com/qasaur/gremgo) with alterations to make it compatible with [AWS Neptune](https://aws.amazon.com/neptune/) which is a "Fast, reliable graph database built for the cloud".

gremgo is a fast, efficient, and easy-to-use client for the TinkerPop graph database stack. It is a Gremlin language driver which uses WebSockets to interface with Gremlin Server and has a strong emphasis on concurrency and scalability. Please keep in mind that gremgo is still under heavy development and although effort is being made to fully cover gremgo with reliable tests, bugs may be present in several areas.

**Modifications were made to `gremgo` in order to "support" AWS Neptune's lack of Gremlin-specific features,  like no support query bindings among others. See differences in Gremlin support here: [AWS Neptune Gremlin Implementation Differences](https://docs.aws.amazon.com/neptune/latest/userguide/access-graph-gremlin-differences.html)**

Installation
==========
```
go get github.com/schwartzmx/gremgo-neptune
dep ensure
```

Documentation
==========

* [GoDoc](https://godoc.org/github.com/schwartzmx/gremgo-neptune)

Example
==========
```go
package main

import (
	"fmt"
	"log"

	"github.com/schwartzmx/gremgo-neptune"
)

func main() {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs) // Example of connection error handling logic

	dialer := gremgo.NewDialer("ws://127.0.0.1:8182") // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.Dial(dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    	return
	}
	res, err := g.Execute( // Sends a query to Gremlin Server with bindings
		"g.V(x)"
	)
	if err != nil {
		fmt.Println(err)
    	return
	}
	fmt.Printf("%+v", res)
}
```

Authentication
==========
The plugin accepts authentication creating a secure dialer where credentials are setted.
If the server where are you trying to connect needs authentication and you do not provide the 
credentials the complement will panic.

```go
package main

import (
	"fmt"
	"log"

	"github.com/schwartzmx/gremgo-neptune"
)

func main() {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs) // Example of connection error handling logic

	dialer := gremgo.NewSecureDialer("127.0.0.1:8182", "username", "password") // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.Dial(dialer, errs) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    	return
	}
	res, err := g.Execute("g.V(x)")
	if err != nil {
		fmt.Println(err)
    	return
	}
	fmt.Printf("%+v", res)
}
```

License
==========

Copyright (c) 2016 Marcus Engvall

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
