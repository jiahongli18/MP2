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
	channel := make(chan string, 1)
	//channel2 := make(chan string)

	//TODO: Sends the Exit signal from exit to startServer
	go startServer(channel)
	go exit(channel)

	if <-channel == "EXIT"{
		return
	}
}

func startServer(channel chan string ) (){
	//Scan and Parse in line argument for the port number
	arguments := os.Args
	if len(arguments) != 2 {
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
		//fmt.Print(username)
		m[username] = c

		signal := "EXIT"
		//fmt.Printf(signal)
		//signal = <-channel

		go handleConnection(c, m, signal)
	}
}

//handleConnection handles communication between the clients, checks if the receivers are connected, and delivers the message
func handleConnection(c net.Conn,m map[string]net.Conn, signal string){
	fmt.Printf("")
	for {
		decoder := gob.NewDecoder(c) //initialize gob decoder
		message := new(utils.Message)
		_ = decoder.Decode(message)

		if receiverChannel, ok := m[message.Receiver]; ok {
			// fmt.Println(receiverChannel, "is in map")
			encoder := gob.NewEncoder(receiverChannel)
			msg := utils.Message{message.Sender, message.Receiver, message.Content}
			encoder.Encode(msg)
		} else {
			encoder := gob.NewEncoder(c)
			msg := utils.Message{"error", "error", "error"}
			encoder.Encode(msg)
		}
		 /*if signal != "EXIT"{
		else{
			//TODO:Sends the termination signal to all the connected clients
			fmt.Print("m")
			for _, receiverChannel:= range m{
				//encoder := gob.NewEncoder(receiverChannel)
				fmt.Print(receiverChannel)
				//msg :=utils.Message{"STOP", "STOP", "STOP"}
				//encoder.Encode(msg)
			}*/
		}

	}


//Exit the server program after getting the user command
func exit(channel chan string) {
	fmt.Println("Waiting for exit command...")
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

