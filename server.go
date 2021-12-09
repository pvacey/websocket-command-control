package main

import (
    "flag"
	"log"
	"net/http"
	//"fmt"
	//"time"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func serveWS(res http.ResponseWriter, req *http.Request) {
	log.Print("Serve Websocket to ", req.RemoteAddr)
	
	// upgrade the http connection to a websocket
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	
	// i think i want one goroutine that writes messages to clients 
	// which maybe keeps the connection alive with pings
	// and a seperate goroutine that reads the response messages from clients
	// and stores them per client
	
	// enter a loop, read the message, print it, reply
	for {
		_, payload, err := conn.ReadMessage()
	    // err is returned when the connection is closed by the client
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("incoming message: ", string(payload))
		err = conn.WriteMessage(1, []byte("hello"))
		if err != nil {
			log.Println(err)
			return
		}
	}
	
}

func main() {
	addr := flag.String("address", ":8000", "http address to bind server to")
	flag.Parse()
	log.Println("Starting WebSocker Server on ",*addr)
	http.HandleFunc("/", serveWS)
	err := http.ListenAndServe(*addr, nil) 
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
