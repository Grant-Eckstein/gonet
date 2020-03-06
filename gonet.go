package gonet

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
)

// Basic error handling
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

/*
	- This is the standard object for this system.
	- The general idea here is to create a "connection"
	  then either:
		1. Add trigger function(s) and then call Connection c.Listen(port string)
		2. Call Connection c.Send(address string, cmd string, data []byte)
 */
type Connection struct {
	Protocol string
	Length   int
	Triggers map[string]func(data []byte) []byte
}

// Creates Connection object
func NewConnection(protocol string, length int) Connection {
	conn := Connection{Protocol: protocol,
		Length:   length,
		Triggers: make(map[string]func(data []byte) []byte),
	}

	return conn
}


/*
	- This is the cool part of the library.
	- Each trigger is a command and a function.
	- When you call Connection c.ResolveMessageTriggers(m Message),
	  the message's command string is comapred aganst each entry
	  and the message data is used as input for the matching trigger's
	  function. c.ResolveMessageTriggers(m Message) return's this output.
 */
func (c *Connection) Trigger(cmd string, function func(data []byte) []byte) {
	c.Triggers[cmd] = function
}


// Returns byte slice of specified message length
func (c *Connection) Listen(port string) Message {
	ln, err := net.Listen(c.Protocol, port)
	check(err)

	conn, err := ln.Accept()
	data := make([]byte, c.Length)
	_, err = conn.Read(data)
	conn.Close()
	check(err)

	m := decodeMessage(data)

	return m
}

// Sends byte slice of specified message length
func (c *Connection) Send(address string, cmd string, data []byte) {
	m := NewMessage(cmd, data)

	j, _ := m.Export()

	conn, err := net.Dial(c.Protocol, address)
	check(err)

	_, err = conn.Write(j)
	conn.Close()
}

// Internal function to resolve message triggers
func (c *Connection) lookupTriggers(cmd string) func(data []byte) []byte {
	f := c.Triggers[cmd]
	if f != nil {
		return f
	} else {
		log.Fatal("Command does not exist")
		os.Exit(1)
		return nil
	}
}

// Called from a server to get message output when using the trigger system.
func (c *Connection) ResolveMessageTriggers(m Message) []byte {
	return c.lookupTriggers(m.Cmd)(m.Data)
}

/*
	Basic data type for message delivery
	- A few assumptions are made:
		1. Message.Cmd is a small string (ex. "save")
		2. Data is obviously a byte slice
	- Note that the message size is specifiec in NewConnection()
 */
type Message struct {
	Cmd  string
	Data []byte
}

/*
	Creates message object
	- This is currently only intended to be used internally
	  however I anticipate applications so it is exported
 */
func NewMessage(cmd string, data []byte) Message {
	return Message{Cmd: cmd, Data: data}
}


// Returns JSONifyed message object.
// Currently only used in Connection.Send(address string, cmd string, data []byte)
func (m *Message) Export() ([]byte, error) {
	return json.Marshal(&m)
}

// Parses raw response, returns Message object
func decodeMessage(data []byte) Message {
	var m Message
	err := json.Unmarshal(bytes.Trim(data, "\x00"), &m)
	check(err)
	return m
}


