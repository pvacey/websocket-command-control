package main

import (
	"fmt"
	"log"
	"time"
)

type Controller struct {
	// a collection of peers
	peers map[*Peer]bool
	// a collection of admins
	admins  map[*Admin]bool
	// messages that need to be sent out to all peers
	commands chan string
	// replies to send to all admins
	results chan CommandResult
	// channel to add new peers
	addPeer chan *Peer
	// channel to add new admins
	addAdmin chan *Admin
	// channel to remove new peers
	removePeer chan *Peer
	// channel to remove admins
	removeAdmin chan *Admin
}

func initController() *Controller {
	return &Controller{
		peers:      map[*Peer]bool{},
		admins:      map[*Admin]bool{},
		commands:     make(chan string),
		results:     make(chan CommandResult),
		addPeer:    make(chan *Peer),
		addAdmin:    make(chan *Admin),
		removePeer: make(chan *Peer),
		removeAdmin: make(chan *Admin),
	}
}

func (c *Controller) run() {
	for {
		select {
		case cmd := <-c.commands:
			for p := range c.peers {
				p.inbox <- cmd
			}
		case res := <-c.results:
			for a := range c.admins {
				a.inbox <- res
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
		case a := <-c.addAdmin:
			log.Println("adding admin ", a.websocket.RemoteAddr())
			c.admins[a] = true
		case a := <-c.removeAdmin:
			log.Println("removing admin ", a.websocket.RemoteAddr())
			delete(c.admins, a)
			a.websocket.Close()
		}
	}
}

func (c *Controller) getPeerHistory() {
	for p := range c.peers {
		fmt.Printf("\n%s@%s\n", p.info["username"], p.info["hostname"])
		for _, cr := range p.history {
			fmt.Println("---")
			fmt.Println(cr.Command)
			fmt.Println("")
			fmt.Println(cr.Result)
		}
	}
}
