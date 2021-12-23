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
	// upgrade the http connection to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
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

func handleNewAdminConnection(c *Controller, res http.ResponseWriter, req *http.Request) {
	log.Print("new admin connection from ", req.RemoteAddr)

	// upgrade the http connection to a websocket
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	admin := initAdmin(conn, c)

	c.addAdmin <- admin
	go admin.writer()
}

func emit(c *Controller, res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	body, _ := ioutil.ReadAll(req.Body)
	c.commands <- string(body)
}

func main() {
	addr := flag.String("address", ":8000", "http address to bind server to")
	flag.Parse()

	controller := initController()
	go controller.run()

	log.Println("Starting WebSocket Server on ", *addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handleNewPeerConnection(controller, w, r) })
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) { emit(controller, w, r) })
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) { controller.getPeerHistory() })
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) { handleNewAdminConnection(controller, w, r) })

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
