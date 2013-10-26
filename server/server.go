// Copyright (c) 2013 Joshua Elliott
// Released under the MIT License
// http://opensource.org/licenses/MIT

package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/mattbaird/turnpike"
	"log"
	"net/http"
)

func main() {
	s := turnpike.NewServer(false)
	http.Handle("/", websocket.Handler(s.HandleWebsocket))
	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
