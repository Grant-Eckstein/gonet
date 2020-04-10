package gonet

import (
	"io/ioutil"
	"log"
	"net"
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
