# Golang Networking Library [![Go Report Card](https://goreportcard.com/badge/github.com/GrantEthanEckstein/gonet)](https://goreportcard.com/report/github.com/GrantEthanEckstein/gonet) [![GoDoc](https://godoc.org/github.com/GrantEthanEckstein/gonet?status.svg)](https://godoc.org/github.com/GrantEthanEckstein/gonet)
*Simplistic Golang networking*


### Overview
This is a client-server networking library I built for another project. It is intended to be a extensible base for networking projects.

### Example Usage
##### Client
```go
package main

import (
	"gonet"
)

func main() {
		conn := gonet3.NewDataConnection(gonet3.Host{
			Protocol:   "tcp",
			Address:    "localhost",
			Port:       "8080",
		}, gonet3.Host{})

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

	conn := gonet3.NewDataConnection(gonet3.Host{}, gonet3.Host{
		Protocol: "tcp",
		Address:  "localhost",
		Port:     "8080",
	})

	data := conn.Recv()
	fmt.Printf("Recieved \"%s\"\n", data)
}
```
