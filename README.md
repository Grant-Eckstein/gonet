# Golang Networking Library
*Simplistic Golang networking*


### Overview
This is a client-server networking library I built for a project. It has two main features:
1. Simple send/receive functions
2. A command parsing system

### Example Usage
*Client*
```go
package main

import "gonet"

func main() {
	// Specify protocol and maximum message length
	c := gonet.NewConnection("tcp", 2048)
	
	// Specify destination address, trigger, and data to send
 	c.Send("127.0.0.1:8080", "hello", []byte("Jello"))
}
```

*Server*
```go
package main

import (
	"fmt"
	"gonet"
)

func hello(name []byte) []byte {
	return append([]byte("Hello, "), name...)
}

func main() {
	// Specify protocol and maximum message length
	l := gonet.NewConnection("tcp", 2048)
	
	// When a message with the trigger "hello" is sent, 
	// it's data will be ran through hello()
	l.Trigger("hello", hello)
	
	recv := l.Listen(":8080")
	
	// ResolveMessageTriggers will return the output 
	// of the appropriate trigger function
	message := string(l.ResolveMessageTriggers(recv))
	fmt.Println(message)
}
```
