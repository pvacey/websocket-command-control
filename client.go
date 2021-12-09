package main

import (
	"log"
	"github.com/gorilla/websocket"
	"time"
)

func send(c *websocket.Conn, msg string) {
	err := c.WriteMessage(1, []byte(msg))
	if err != nil {
		log.Println("send error: ", err)
	}
}

func recieve(c *websocket.Conn) string {
	_, payload, err := c.ReadMessage()
	if err != nil {
		log.Println("recieve error: ", err)
	}
	return string(payload)
}

func main() {
	host := "ws://localhost:8000"
	log.Println("Connecting to host: ", host)
	
	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	
	send(conn, "hey")
	msg := recieve(conn)
	log.Println("recieved message: ", msg)

	// send the proper disconnect signal to the other end
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(1)
}
