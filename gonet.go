package gonet

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		//log.Fatal(e)
		panic(e)
	}
}

type Connection struct {
	Protocol string
	Length   int
	Triggers map[string]func(data []byte) []byte
}

type Message struct {
	Cmd  string
	Data []byte
}

func NewMessage(cmd string, data []byte) Message {
	return Message{Cmd: cmd, Data: data}
}

func (m *Message) Export() ([]byte, error) {
	return json.Marshal(&m)
}

func NewConnection(protocol string, length int) Connection {
	conn := Connection{Protocol: protocol,
		Length:   length,
		Triggers: make(map[string]func(data []byte) []byte),
	}

	return conn
}

func decodeMessage(data []byte) Message {

	var m Message
	err := json.Unmarshal(bytes.Trim(data, "\x00"), &m)
	check(err)
	return m
}

func (c *Connection) Trigger(cmd string, function func(data []byte) []byte) {
	c.Triggers[cmd] = function
}

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

func (c *Connection) Send(address string, cmd string, data []byte) {
	m := NewMessage(cmd, data)

	j, _ := m.Export()

	conn, err := net.Dial(c.Protocol, address)
	check(err)

	_, err = conn.Write(j)
	conn.Close()
}

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
func (c *Connection) ResolveMessageTriggers(m Message) []byte {
	return c.lookupTriggers(m.Cmd)(m.Data)
}
