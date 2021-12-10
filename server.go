package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	//"time"
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

func (c *Controller) run () {
	for {
		select{
		case msg := <-c.outbox:
			for p := range c.peers {
				p.inbox <- msg	
			}
		case p := <-c.addPeer:
			log.Println("adding peer ", p.websocket.RemoteAddr())
			c.peers[p] = true
		case p := <-c.removePeer:
			log.Println("removing peer ", p.websocket.RemoteAddr())
			delete(c.peers, p)
		}
	}
}

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
	//log.Println("reader closed")
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
	//log.Println("writer closed")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleNewPeerConnection(c *Controller, res http.ResponseWriter, req *http.Request) {
	//log.Print("Serve Websocket to ", req.RemoteAddr)

	// upgrade the http connection to a websocket
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	//defer conn.Close()
	peer := initPeer(conn, c)

	c.addPeer <-peer
	go peer.reader()
	go peer.writer()
}

func emit(c *Controller, res http.ResponseWriter, req *http.Request) {
	c.outbox <-	"hello there!!!!!"
}

func main() {
	addr := flag.String("address", ":8000", "http address to bind server to")
	flag.Parse()

	controller := initController()
	go controller.run()

	log.Println("Starting WebSocker Server on ", *addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handleNewPeerConnection(controller, w, r) })
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) { emit(controller, w, r) })

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
