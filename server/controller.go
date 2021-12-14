package main

import (
	"log"
	"time"
)

type Controller struct {
	// a collection of peers to loop over
	peers map[*Peer]bool
	// messages that need to be sent out to all peers
	outbox chan string
	// channel to add new peers
	addPeer chan *Peer
	// channel to remove new peers
	removePeer chan *Peer
}

func initController() *Controller {
	return &Controller{
		peers:      map[*Peer]bool{},
		outbox:     make(chan string),
		addPeer:    make(chan *Peer),
		removePeer: make(chan *Peer),
	}
}

func (c *Controller) run() {
	for {
		select {
		case msg := <-c.outbox:
			for p := range c.peers {
				p.inbox <- msg
			}
		case <-time.After(30 * time.Second):
			for p := range c.peers {
				p.ping <- true
			}
		case p := <-c.addPeer:
			log.Println("adding peer ", p.websocket.RemoteAddr())
			c.peers[p] = true
		case p := <-c.removePeer:
			log.Println("removing peer ", p.websocket.RemoteAddr())
			delete(c.peers, p)
			p.websocket.Close()
		}
	}
}
