package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
  "./utils"
	"strings"
)

func main() {
	//Scan user input for host:port
	arguments := os.Args
	if len(arguments) != 3 {
		fmt.Println("Please provide host:port and the username.")
		return
	}
	fmt.Print("Type EXIT if you want to leave. Press anything else to continue.\n")
	//connect to provided host:post via the net library
	c := TCPDial(arguments)

	go listen(c)

	for {
		if "EXIT" == userExit() {
			return
		}
		sender, receiver, content := getUserInput()

		msg := utils.Message{sender, receiver, content}

		messaging(msg, c)
	}

}

func listen(c net.Conn){
  for {
	  decoder := gob.NewDecoder(c) //initialize gob decoder
	  //Decode message struct and print it
	  message := new(utils.Message)
	  _ = decoder.Decode(message)

	  //TODO:Receive the termination signal and stop listening
	  /*if (*message == utils.Message{"STOP", "STOP", "STOP"}){
	  	fmt.Print("hi")
	  	return
	  }*/
	  if(*message == utils.Message{"error", "error", "error"}) {
		  fmt.Printf("\nError: the person you are sending to has not been connected yet.\n")
		  fmt.Printf("Type EXIT or enter Sender: ")
	  } else {
		  fmt.Printf("Received message from %s\nMessage: %s\n", message.Sender, message.Content)
		  fmt.Printf("Type EXIT or enter Sender: ")
	  }
  }
}

func messaging(msg utils.Message, c net.Conn) {
  //create a gob encoder and code the message struct
    encoder := gob.NewEncoder(c)
    _ = encoder.Encode(msg)
}

func TCPDial(arguments []string)(c net.Conn) {
  CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

  username := arguments[2]
  fmt.Fprintf(c, username + "\n")

  return c
}

//Gets the sender, receiver and message content from the user input
func getUserInput()(sender string, receiver string, content string) {
	reader := bufio.NewReader(os.Stdin)

	//scan user input for message contents
	fmt.Print("Sender: ")
	sender, _ = reader.ReadString('\n')

	fmt.Print("Receiver: ")
	receiver, _ = reader.ReadString('\n')

	fmt.Print("Message content: ")
	content, _ = reader.ReadString('\n')

	return sender, receiver, content
}

//Exit the client program after getting the user command
func userExit()(exit string){
	arguments := os.Args

	reader := bufio.NewReader(os.Stdin)
	var cmd string
	cmd, _ = reader.ReadString('\n')
	if strings.TrimSpace(cmd) == "EXIT" {
		fmt.Printf("Client %q is exiting...\n", arguments[2])
		return "EXIT"
	}
	return " "
}