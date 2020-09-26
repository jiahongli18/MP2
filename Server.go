package main

import (
	"./utils"
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)
//Starts the listening side from the server
func StartServer(channel chan string) {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	//Accept multiple client connections concurrently
	for {
		fmt.Println("h")
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// netData, err := bufio.NewReader(c).ReadString('\n')
		go handleConnection(c, channel)
	}

}

//Check if a client is connected
func checkConn(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

//TODO:A thread to handle the communication with each client
func handleConnection(c net.Conn, channel chan string) {
	var clients []string
	fmt.Print("..")
	for {
		decoder := gob.NewDecoder(c)
		message := new(utils.Message)

		//Create a gob decoder to decode message from the client
		_ = decoder.Decode(message)
		clients = append(clients, message.From)
		fmt.Printf("To : %+v \nFrom : %+v \nContent : %+v", message.To, message.From, message.Content)

		result := checkConn(clients, message.To)

		if result == true {
			//Create a gob encoder to send message to the client
			encoderMsg := gob.NewEncoder(c)
			fmt.Println(clients)
			_ = encoderMsg.Encode(message.Content)
		}
		if result == false {
			encoderErr := gob.NewEncoder(c)
			_ = encoderErr.Encode("To-client is not connected.")
		}
		signal := <-channel
		if signal == "Termination" {
			fmt.Print("Received")

			return
		}
		/*netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		temp := strings.TrimSpace(string(netData))
		if temp == "STOP" {
			break
		}
		fmt.Println(temp)
		*/

	}
	c.Close()
}

//Exit thread waits for user's exit command, exits the program, and sends termination signal to all clients
func Exit(channel chan string) {
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

func main() {
	channel := make(chan string)
	StartServer(channel)
	go Exit(channel)
	fmt.Println(<-channel)
	//go StartServer(channel)

}
