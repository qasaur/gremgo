# gremgo - Gremlin client for Golang

gremgo is a fast, efficient, and easy-to-use client for the TinkerPop graph database stack. Functionality is limited to simple executions of commands with bindings at the moment, but there are plans to include session-based interactions and other more advanced features in the future.

Installation
==========
```
go get github.com/qasaur/gremgo
```

Example
==========
```
package main

import (
	"fmt"
	"log"

	"github.com/qasaur/gremgo"
)

func main() {
	c, err := gremgo.Dial("127.0.0.1:8182") // Returns a gremgo client to interact with
	if err != nil {
		fmt.Println(err)
    return
	}
	res, err := c.Execute("g.V(x)", map[string]string{"x": "1234"}) // Sends a query to Gremlin Server with bindings
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
