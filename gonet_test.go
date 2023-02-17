package gonet

import (
	"fmt"
	"log"
	"net/http"
	"os"

	client "github.com/Grant-Eckstein/gonet/client"
	server "github.com/Grant-Eckstein/gonet/server"
)

func ExampleServe() {
	h1 := server.NewHandler("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("This is an example server.\n"))
		if err != nil {
			panic(err)
		}
	})

	h2 := server.NewHandler("/kill", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("kill") == "true" {
			os.Exit(0)
		} else {
			log.Println("Didn't get kill ")
		}
	})

	err := server.Serve([]string{"localhost"}, []server.Handler{h1, h2}, 443)
	if err != nil {
		panic(err)
	}
}

func ExampleSend() {
	res, err := client.Send("localhost:443", []byte("hello, world"), 50)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server replied: '%v'\n", string(res))
}
