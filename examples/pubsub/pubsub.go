package main

import (
	"flag"
	"fmt"
	"github.com/mattbaird/turnpike"
	"net/http"
	"time"
)

var (
	runMode *string = flag.String("runMode", "server", "Run mode. [server | client]")
)

func testHandler(uri string, event interface{}) {
	fmt.Printf("Received event: %s\n", event)
}

func main() {
	flag.Parse()
	if *runMode == "server" {

		s := turnpike.NewServer()

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

		c.Subscribe("event:test", testHandler)

		for {
			c.Publish("event:test", "test")
			<-time.After(time.Second)
		}
	} else {
		fmt.Printf("runMode must be one of server or client, you passed %s", runMode)
	}
}
