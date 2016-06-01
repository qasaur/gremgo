# gremgo

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/qasaur/gremgo) [![Build Status](https://travis-ci.org/qasaur/gremgo.svg?branch=master)](https://travis-ci.org/qasaur/gremgo) [![Coverage Status](https://coveralls.io/repos/github/qasaur/gremgo/badge.svg?branch=master)](https://coveralls.io/github/qasaur/gremgo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/qasaur/gremgo)](https://goreportcard.com/report/github.com/qasaur/gremgo)

gremgo is a fast, efficient, and easy-to-use client for the TinkerPop graph database stack. It is a Gremlin language driver which uses WebSockets to interface with Gremlin Server and has a strong emphasis on concurrency and scalability. Functionality is limited to simple executions of commands with bindings at the moment, but there are plans to include session-based interactions and other more advanced features in the future.

Installation
==========
```
go get github.com/qasaur/gremgo
```

Documentation
==========

* [GoDoc](https://godoc.org/github.com/qasaur/gremgo)

Example
==========
```go
package main

import (
	"fmt"
	"log"

	"github.com/qasaur/gremgo"
)

func main() {
	dialer := Ws{Host: "127.0.0.1:8182"} // Returns a WebSocket dialer to connect to Gremlin Server
	g, err := gremgo.Dial(dialer) // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    	return
	}
	res, err := g.Execute( // Sends a query to Gremlin Server with bindings
		"g.V(x)",
		map[string]string{"x": "1234"}
	)
	if err != nil {
		fmt.Println(err)
    	return
	}
	fmt.Println(res)
}
```

License
==========

Copyright (c) 2016 Marcus Engvall

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
