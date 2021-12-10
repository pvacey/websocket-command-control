package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

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
	peer := initPeer(conn, c)

	c.addPeer <- peer
	go peer.reader()
	go peer.writer()
}

func emit(c *Controller, res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	c.outbox <- string(body)
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
