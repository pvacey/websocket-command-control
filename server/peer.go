package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Peer struct {
	// messages going to the peer
	inbox chan string
	// a collection of the messages returned from the peer
	outbox []string
	// the controller to send messages
	controller *Controller
	// the actuall connection to the peer
	websocket *websocket.Conn
}

func initPeer(w *websocket.Conn, c *Controller) *Peer {
	return &Peer{
		inbox:      make(chan string),
		outbox:     make([]string, 1),
		websocket:  w,
		controller: c,
	}
}

// this can likely be repurposed so that it reads messages and store them in
// the peer's outbox array. god only knows what to do with it then
// ignore pong messages, once those exist
func (p *Peer) reader() {
	for {
		_, payload, err := p.websocket.ReadMessage()
		// err is returned when the connection is closed by the client
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("incoming message: ", string(payload))
	}
}

// need to change the loop to watch for both messages in the peer inbox
// as well as to send a periodic ping message with time.After (use select)
func (p *Peer) writer() {
	for msg := range p.inbox {
		err := p.websocket.WriteMessage(1, []byte(msg))
		if err != nil {
			log.Println(err)
			break
		}
	}
	// drop the peer once the loop is broken by an error or some other condition
	p.controller.removePeer <- p
}
