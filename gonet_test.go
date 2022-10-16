package gonet

import (
	"log"
	"net/http"
	"os"
)

func ExampleServer() {
	h1 := NewHandler("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("This is an example server.\n"))
		if err != nil {
			panic(err)
		}
	})

	h2 := NewHandler("/kill", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("kill") == "true" {
			os.Exit(0)
		} else {
			log.Println("Didn't get kill ")
		}
	})

	err := Serve([]string{"localhost"}, []Handler{h1, h2})
	if err != nil {
		panic(err)
	}
}
