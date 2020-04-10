package gonet

import (
	"encoding/json"
	"io/ioutil"
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

// Root structure for gonet, stores information necessary for each connection.
type Host struct {
	Protocol string
	Address  string
	Port     string
}

// Purely for organizational needs. Provides logical seperation between
// DataConnection.Send and DataConnection.Recv.
type DataConnection struct {
	Target   Host
	Listener Host
}

// Returns active Conn from Host configuration.
func (h *Host) GetConnection() net.Conn {
	conn, err := net.Dial(h.Protocol, net.JoinHostPort(h.Address, h.Port))
	check(err)

	return conn
}

// DataConnection constructor
func NewDataConnection(target Host, listen Host) DataConnection {
	return DataConnection{Target: target, Listener: listen}
}

// Send a byteslice to target using specified configuration.
func (dConn *DataConnection) Send(data []byte) {
	settings := dConn.Target
	conn, err := net.Dial(settings.Protocol, net.JoinHostPort(settings.Address, settings.Port))
	check(err)

	_, err = conn.Write(data)
	check(err)
	conn.Close()
}

// Listen for a byteslice on a specific port. Must have localhost specification,
// this distinction exists for applications listening on multiple internal addresses.
func (dConn *DataConnection) Recv() []byte {
	settings := dConn.Listener

	ln, err := net.Listen(settings.Protocol, net.JoinHostPort(settings.Address, settings.Port))
	check(err)

	conn, err := ln.Accept()
	check(err)

	data, err := ioutil.ReadAll(conn)
	check(err)

	conn.Close()
	return data
}

/* Here is the trigger functionality */

// Message encapsulation object for trigger system.
type Message struct {
	Data    []byte
	Command string
}

// Implement new host types. This way there exists a incoming and outgoing trigger list.
type TriggerDataConnection struct {
	Target   TriggerHost
	Listener TriggerHost
}

// Host with map of commands and corresponding variables.
type TriggerHost struct {
	Protocol string
	Address  string
	Port     string
	Triggers map[string]func(data []byte) []byte
}

// DataConnection constructor
func NewTriggerDataConnection(target TriggerHost, listen TriggerHost) TriggerDataConnection {
	return TriggerDataConnection{
		Target:   target,
		Listener: listen,
	}
}

// NewTriggerHost constructor. Main use is to initialize the Triggers map.
func NewTriggerHost(protocol string, address string, port string) TriggerHost {
	return TriggerHost{
		Protocol: protocol,
		Address:  address,
		Port:     port,
		Triggers: make(map[string]func(data []byte) []byte),
	}
}

// A placeholder
func NullTriggerHost() TriggerHost {
	return TriggerHost{Triggers: make(map[string]func(data []byte) []byte),}
}

// Send a byteslice to target using specified configuration.
func (tdc *TriggerDataConnection) Send(command string, data []byte) {
	data = newMessage(command, data)
	settings := tdc.Target
	conn, err := net.Dial(settings.Protocol, net.JoinHostPort(settings.Address, settings.Port))
	check(err)

	_, err = conn.Write(data)
	check(err)
	conn.Close()
}

// Listen for a byteslice on a specific port. Must have localhost specification,
// this distinction exists for applications listening on multiple internal addresses.
// After recieving data, run command against trigger list. If a match is found, the output of
// that function is returned. The recieved data is used as input.
func (tdc *TriggerDataConnection) Recv() []byte {
	settings := tdc.Listener

	ln, err := net.Listen(settings.Protocol, net.JoinHostPort(settings.Address, settings.Port))
	check(err)

	conn, err := ln.Accept()
	check(err)

	data, err := ioutil.ReadAll(conn)
	check(err)
	conn.Close()

	// Process message
	var mesg Message
	json.Unmarshal(data, &mesg)
	return settings.resolveMessageTriggers(mesg)
}

// Add function for recieved data to be ran against by command.
func (th *TriggerHost) Trigger(command string, function func(data []byte) []byte) {
	th.Triggers[command] = function
}

// Return point for received messages. Processes trigger function with data as an input.
func (th *TriggerHost) resolveMessageTriggers(m Message) []byte {
	return th.lookupTriggers(m.Command)(m.Data)
}

// Lookup trigger functions by command for each TriggerHost.
func (th *TriggerHost) lookupTriggers(command string) func(data []byte) []byte {
	f := th.Triggers[command]
	if f != nil {
		return f
	} else {
		log.Fatal("Command does not exist")
		os.Exit(1)
		return nil
	}
}

// Include command in data json structure.
func newMessage(command string, data []byte) []byte {
	mesg := Message{Data: data, Command: command}
	out, err := json.Marshal(&mesg)
	check(err)
	return out
}

// Internal function to process json-encoded Messages
func decodeMessage(data []byte) Message {
	var m Message
	err := json.Unmarshal(data, &m)

	check(err)

	return m
}
