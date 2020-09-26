package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
	// "log"
	"./utils"
)

func main() {
	channel := make(chan string)
	go startServer()
	go exit(channel)
	signal := <-channel
	if signal == "Termination" {
		return
	}
}

func startServer() {
	//Scan and Parse in line argument for the port number
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	//Create TCP connection and listen on provided port for requests
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Listening on port" + PORT)
	}
	defer l.Close()

	username := ""
	m := make(map[string]net.Conn)

	for {
		//The server accepts and begins to interact with TCP client
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		netData, _ := bufio.NewReader(c).ReadString('\n')
		username = netData
		m[username] = c

		go handleConnection(c, m)
	}
}

func handleConnection(c net.Conn,m map[string]net.Conn) {
	decoder := gob.NewDecoder(c) //initialize gob decoder
	encoder := gob.NewEncoder(c)

	message := new(utils.Message)
	_ = decoder.Decode(message)
	
	fmt.Println(c)
	fmt.Println(m[message.Sender])

	if val, ok := m[message.Receiver]; ok {
    fmt.Println(val, "is in map")
	encoder2 := gob.NewEncoder(val)
	// msg := utils.Message{message.Sender, message.Receiver, message.Content}

	// encoder2.Encode(msg)
	encoder2.Encode(message.Content)
	} else {
    	encoder.Encode("error")
	}
}

func exit(channel chan string) {
	fmt.Println("Waiting for exit command...")
	for {
		reader := bufio.NewReader(os.Stdin)
		var cmd string
		cmd, _ = reader.ReadString('\n')
		if strings.TrimSpace(cmd) == "EXIT" {
			fmt.Println("Server is exiting...")
			//Sends the termination signal to all the connected clients
			channel <- "Termination"
			return
		}
	}
}

