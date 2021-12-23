package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type CommandResult struct {
	HostInfo map[string]string
	Command string
	Result  string
	Err     string
}

type Peer struct {
	// messages going to the peer
	inbox chan string
	// a channel of live replies from the peer
	outbox  chan CommandResult
	// a collection of the messages returned from the peer
	history []CommandResult
	// recieve triggers to send ping/keepalives
	ping chan bool
	// the controller to send messages
	controller *Controller
	// the actuall connection to the peer
	websocket *websocket.Conn
	// store host information from the first connect message
	info map[string]string
}

func initPeer(w *websocket.Conn, c *Controller) *Peer {
	return &Peer{
		inbox:      make(chan string),
		outbox:    make(chan CommandResult),
		history:    make([]CommandResult, 0),
		ping:       make(chan bool),
		websocket:  w,
		controller: c,
		info:       map[string]string{},
	}
}

func (p *Peer) reader() {
	_, payload, _ := p.websocket.ReadMessage()
	// read the first message and set the peer's info based	
	// on the content of this first message
	cmdRes := CommandResult{}
	json.Unmarshal(payload, &cmdRes)
	p.info = cmdRes.HostInfo
	p.info["address"] = fmt.Sprint(p.websocket.RemoteAddr())
	log.Println("peer info: ", p.info)

	for {
		_, payload, err := p.websocket.ReadMessage()
		// err is returned when the connection is closed by the client
		if err != nil {
			log.Println(err)
			break
		}
		json.Unmarshal(payload, &cmdRes)
		log.Println("incoming reply to cmd:", cmdRes.Command)
		// return results to listening admins and track the history
		p.controller.results <- cmdRes
		p.history = append(p.history, cmdRes)
	}
}

func (p *Peer) writer() {
Looper:
	for {
		select {
		case msg := <-p.inbox:
			err := p.websocket.WriteMessage(1, []byte(msg))
			if err != nil {
				log.Println(err)
				break Looper
			}
		case <-p.ping:
			log.Println("sending ping to", p.websocket.RemoteAddr())
			err := p.websocket.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println(err)
				break Looper
			}
		}
	}
	// drop the peer once the loop is broken by an error or some other condition
	p.controller.removePeer <- p
}
