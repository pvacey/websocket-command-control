package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Admin struct {
	// responses from peers
	inbox chan CommandResult
	// recieve triggers to send ping/keepalives
	ping chan bool
	// the controller to send messages
	controller *Controller
	// the actuall connection to the peer
	websocket *websocket.Conn
}

func initAdmin(w *websocket.Conn, c *Controller) *Admin {
	return &Admin{
		inbox:     make(chan CommandResult),
		ping:       make(chan bool),
		websocket:  w,
		controller: c,
	}
}

func (a *Admin) writer() {
Looper:
	for {
		select {
		case msg := <-a.inbox:
			log.Println("sending message to the admin")
			json_msg,_ := json.Marshal(msg)
			log.Println(string(json_msg))
			err := a.websocket.WriteMessage(1, []byte(json_msg))
			if err != nil {
				log.Println(err)
				break Looper
			}
		case <-a.ping:
			log.Println("sending ping to", a.websocket.RemoteAddr())
			err := a.websocket.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println(err)
				break Looper
			}
		}
	}
	// drop the peer once the loop is broken by an error or some other condition
	a.controller.removeAdmin <- a
}
