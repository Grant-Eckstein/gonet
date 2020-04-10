# Golang Networking Library [![Go Report Card](https://goreportcard.com/badge/github.com/GrantEthanEckstein/gonet)](https://goreportcard.com/report/github.com/GrantEthanEckstein/gonet) [![GoDoc](https://godoc.org/github.com/GrantEthanEckstein/gonet?status.svg)](https://godoc.org/github.com/GrantEthanEckstein/gonet)
*Simplistic Golang networking*


### Overview
This is a client-server networking library I built for another project. It is intended to be a extensible base for networking projects.

### Example Usage - Without triggers
##### Client
```go
package main

import (
	"gonet"
)

func main() {
		conn := gonet.NewDataConnection(gonet.Host{
			Protocol:   "tcp",
			Address:    "localhost",
			Port:       "8080",
		}, gonet.Host{})

		conn.Send([]byte("Hello, world!"))
}
```
<br>

##### Server
```go
package main

import (
	"fmt"
	"gonet"
)

func main() {

	conn := gonet.NewDataConnection(gonet.Host{}, gonet.Host{
		Protocol: "tcp",
		Address:  "localhost",
		Port:     "8080",
	})

	data := conn.Recv()
	fmt.Printf("Recieved \"%s\"\n", data)
}
```
<br>

### Example Usage - With triggers
##### Client
```go
package main

import (
	"fmt"
	"gonet"
)

func main() {

	target := gonet.NewTriggerHost("tcp", "localhost", "8080")
	conn := gonet.NewTriggerDataConnection(target, gonet.NullTriggerHost())

	data := []byte("Grant")
	conn.Send("test-cmd", data)

	fmt.Printf("Sent \"%s\"\n", data)
}
```
<br>

##### Server
```go
package main

import (
	"gonet"
	"fmt"
)

func main() {

	listener := gonet.NewTriggerHost("tcp", "localhost", "8080")
	conn := gonet.NewTriggerDataConnection(gonet.NullTriggerHost(), listener)

	conn.Listener.Trigger("test-cmd", func(data []byte) []byte {
		return append([]byte("Hello, "), data...)
	})

	data := conn.Recv()

	fmt.Printf("Output of recieved data is \"%s\"\n", data)
}

```
