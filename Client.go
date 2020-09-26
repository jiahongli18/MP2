package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
  "./utils"
)

func main() {
	//Scan user input for host:port
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port and the username.")
		return
	}

	//connect to provided host:post via the net library
	c := TCPDial(arguments)
  

  for {
    sender, receiver, content := getUserInput()
    msg := utils.Message{sender, receiver, content}

    messaging(msg, c)
    listen(c)
  }
}

func listen(c net.Conn) {
     decoder := gob.NewDecoder(c) //initialize gob decoder
    
	  //Decode message struct and print it
	  message := new(utils.Message)
	  _ = decoder.Decode(message)

    fmt.Printf("Received message from %q\nMessage: %s\n", message.Sender, message.Content)
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

func getUserInput()(sender string, receiver string, content string) {
    reader := bufio.NewReader(os.Stdin)

    //scan user input for message contents
    fmt.Print("Sender: ")
    sender, _ = reader.ReadString('\n')

    fmt.Print("Receiver: ")
    receiver, _ = reader.ReadString('\n')

    fmt.Print("Message content: ")
    content, _ = reader.ReadString('\n')

    return sender,receiver,content
}