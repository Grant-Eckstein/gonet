# Golang Networking Library [![Go Report Card](https://goreportcard.com/badge/github.com/GrantEthanEckstein/gonet)](https://goreportcard.com/report/github.com/GrantEthanEckstein/gonet) [![GoDoc](https://godoc.org/github.com/GrantEthanEckstein/gonet?status.svg)](https://godoc.org/github.com/GrantEthanEckstein/gonet)
*Simplistic Golang networking*


### Overview
This is a client-server networking library I built for another project. It is an extensible base for other networking projects.

### Example Usage - Without triggers
##### Client
```go
package main

import (
	"gonet"
)

func main() {
	conn := gonet.NewTLS("tcp", true)

	conn.Send([]byte("Hello from client"), "localhost:4444")
	data := conn.Recv(32, "localhost:4445", "example.crt", "example.key")
	fmt.Printf("Recieved \"%s\" over secure websocket.\n", data)
}
```
<br>

##### Server
```go
package main

import (
	"fmt"
	"time"
	"gonet"
)

func main() {

	conn := gonet6.NewTLS("tcp")

	data := conn.Recv(128, "localhost:4444", "example.crt", "example.key")
	fmt.Printf("Recieved \"%s\" over secure websocket.\n", data)
	fmt.Println("Sending response...")

	// Typically this is unnecessary unless
	// running both client and server locally.
	time.Sleep(2*time.Second)
	conn.Send([]byte("Hello from server"), "localhost:4445")
}
```
