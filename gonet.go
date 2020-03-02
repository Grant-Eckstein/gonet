package main

import (
	"fmt"
	"log"
	"os"
)

type Listener struct {
	Name   string
	CmdMap map[string]func(data []byte) []byte
}

func newListener(Name string) Listener {
	listener := Listener{Name: Name}
	listener.CmdMap = make(map[string]func(data []byte) []byte)
	return listener
}

func (l *Listener) trigger(cmd string, function func(data []byte) []byte) {
	l.CmdMap[cmd] = function
}

func (l *Listener) lookup(cmd string) func(data []byte) []byte {
	f := l.CmdMap[cmd]
	if f != nil {
		return f
	} else {
		log.Fatal("Command does not exist")
		os.Exit(1)
		return nil
	}

}

func (l *Listener) check(cmd string, data []byte) []byte {
	return l.lookup(cmd)(data)
}

func test(data []byte) []byte {
	return []byte(string(data)[:1])
}

func main() {
	l := newListener("test")
	l.trigger("hello", test)

	output := l.check("hello", []byte("hello, world!"))
	fmt.Println(string(output))

}

