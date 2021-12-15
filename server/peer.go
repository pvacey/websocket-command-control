package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

type Peer struct {
	// messages going to the peer
	inbox chan string
	// a collection of the messages returned from the peer
	outbox []string
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
		outbox:     make([]string, 1),
		ping:       make(chan bool),
		websocket:  w,
		controller: c,
		info:       map[string]string{},
	}
}

// this can likely be repurposed so that it reads messages and store them in
// the peer's outbox array. god only knows what to do with it then
// ignore pong messages, once those exist
func (p *Peer) reader() {
	_, payload, _ := p.websocket.ReadMessage()
	log.Println(string(payload))

	p.info["address"] = fmt.Sprint(p.websocket.RemoteAddr())
	firstMsg := strings.Split(string(payload), " ")
	p.info["os"] = firstMsg[0]
	p.info["workingDir"] = firstMsg[1]
	p.info["hostname"] = firstMsg[2]
	p.info["username"] = firstMsg[3]

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
			log.Println("sending ping to ", p.websocket.RemoteAddr())
			err := p.websocket.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println(err)
				break Looper
			}

		}
	}
	// drop the peer once the loop is broken by an error or some other condition
	p.controller.removePeer <- p
	log.Println(p.info)
}
