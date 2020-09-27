package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
	"./utils"
)

func main() {
	var m = make(map[string]net.Conn)

	channel := make(chan string)
	go startServer(m)
	go exit(channel)
	signal := <-channel
	if signal == "EXIT" {
		exitAllClients(m)
		return
	}
}

func startServer(m map[string]net.Conn) {
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
		fmt.Println("Listening on port" + PORT + ". Please type 'EXIT' to quit.")
	}
	defer l.Close()

	username := ""

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
	for {
		decoder := gob.NewDecoder(c) //initialize gob decoder
		message := new(utils.Message)
		_ = decoder.Decode(message)

		if receiverChannel, ok := m[message.Receiver]; ok {
			encoder := gob.NewEncoder(receiverChannel)
			msg := utils.Message{message.Sender, message.Receiver, message.Content}
			encoder.Encode(msg)
		} else {
			encoder := gob.NewEncoder(c)
			msg := utils.Message{"error", "error", "error"}
			encoder.Encode(msg)
		}
	}
}

func exitAllClients(m map[string]net.Conn) {
	for _, receiverChannel := range m {
        encoder := gob.NewEncoder(receiverChannel)
		msg := utils.Message{"EXIT", "EXIT", "EXIT"}
		encoder.Encode(msg)
    }
}

func exit(channel chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		var cmd string
		cmd, _ = reader.ReadString('\n')
		if strings.TrimSpace(cmd) == "EXIT" {
			fmt.Println("Server is exiting...")
			//Sends the termination signal to all the connected clients
			channel <- "EXIT"
			return
		}
	}
}

