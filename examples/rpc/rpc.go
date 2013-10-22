package main

import (
	"flag"
	"fmt"
	"github.com/mattbaird/turnpike"
	"net/http"
)

var (
	runMode *string = flag.String("runMode", "server", "Run mode. [server | client]")
)

func handleTest(client, uri string, args ...interface{}) (interface{}, error) {
	return "hello world", nil
}

func main() {
	flag.Parse()
	if *runMode == "server" {
		s := turnpike.NewServer()
		s.RegisterRPC("rpc:test", handleTest)

		http.Handle("/ws", s.Handler)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}

	} else if *runMode == "client" {
		c := turnpike.NewClient()
		err := c.Connect("ws://127.0.0.1:8080/ws", "http://localhost/")
		if err != nil {
			panic("Error connecting:" + err.Error())
		}

		resultCh := c.Call("rpc:test")
		result := <-resultCh
		fmt.Printf("Call result is: %s\n", result.Result)
	} else {
		fmt.Printf("runMode must be one of server or client, you passed %s", runMode)
	}
}
