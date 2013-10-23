package main

import (
	"bufio"
	"fmt"
	"github.com/mattbaird/turnpike"
	"os"
	"strconv"
	"strings"
)

const (
	WELCOME     = iota // iota is reset to 0
	PREFIX      = iota
	CALL        = iota
	CALLRESULT  = iota
	CALLERROR   = iota
	SUBSCRIBE   = iota
	UNSUBSCRIBE = iota
	PUBLISH     = iota
	EVENT       = iota
)

func main() {
	c := turnpike.NewClient()
	fmt.Print("Server address (default: localhost:9091/ws)\n> ")
	read := bufio.NewReader(os.Stdin)
	// if _, err := fmt.Scanln(&server); err != nil {
	server, err := read.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}
	server = strings.TrimSpace(server)
	if server == "" {
		server = "localhost:9091/ws"
	}
	if err := c.Connect("ws://"+server, "http://localhost:8070"); err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	fmt.Print(
		"----------------------------------------------------------------------\n",
		"Connected to server at: ", server, "\n",
		"With session id: ", c.SessionId, "\n",
		"Enter WAMP message, parameters separated by spaces\n",
		"PREFIX=1, CALL=2, SUBSCRIBE=5, UNSUBSCRIBE=6, PUBLISH=7\n",
		"----------------------------------------------------------------------\n")

	for {
		fmt.Print(c.ServerIdent, "> ")
		// if _, err := fmt.Scanln(&msgType, &args); err != nil {
		line, err := read.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// fmt.Println(line)

		// get the type
		params := strings.SplitN(line, " ", 2)
		line = params[1]
		msgType, err := strconv.Atoi(params[0])
		if err != nil {
			fmt.Println("Error parsing message type:", params[0])
			continue
		}

		err = nil
		switch msgType {
		case PREFIX:
			var prefix, URI string
			fmt.Sscan(line, &prefix, &URI)
			err = c.Prefix(prefix, URI)
		case CALL:
			args := strings.Split(line, " ")
			resultCh := make(chan turnpike.CallResult)
			resultCh = c.Call(args[0], args[1:])
			r := <-resultCh
			if r.Error != nil {
				fmt.Println(r.Error)
			} else {
				fmt.Printf("Result %s\n", r.Result)
			}
		case SUBSCRIBE:
			eventCh := make(chan interface{})
			err = c.Subscribe(line, func(uri string, event interface{}) {
				eventCh <- event
			})
			if err != nil {
				fmt.Printf("Failed to subscribe :%v", err)
			}
			go receive(eventCh)
		case UNSUBSCRIBE:
			err = c.Unsubscribe(line)
		case PUBLISH:
			args := strings.Split(line, " ")
			if len(args) > 2 {
				err = c.Publish(args[0], args[1], args[2:])
			} else {
				err = c.Publish(args[0], args[1])
			}
			if err != nil {
				fmt.Printf("Failed to publish :%v", err)
			}
		default:
			fmt.Println("Invalid message type:", msgType)
			continue
		}

		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
func receive(eventCh chan (interface{})) {
	for {
		select {
		case msg := <-eventCh:
			fmt.Println("Got event:", msg)
		}
	}
}
