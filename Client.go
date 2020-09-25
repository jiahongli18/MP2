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
	var sender string
	var receiver string
	var content string

	reader := bufio.NewReader(os.Stdin)

	//Scan user input for host:port
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	//connect to provided host:post via the net library
	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

  for {
    //scan user input for message contents
    fmt.Print("Sender: ")
    sender, _ = reader.ReadString('\n')
    if(sender == "EXIT") {
      break
    };

    fmt.Print("Receiver: ")
    receiver, _ = reader.ReadString('\n')
    if(receiver == "EXIT") {
      break
    };
    fmt.Print("Message content: ")
    content, _ = reader.ReadString('\n')
    if(content == "EXIT") {
      break
    };

    //create a gob encoder and code the message struct
    encoder := gob.NewEncoder(c)
    msg := utils.Message{sender, receiver, content}
    _ = encoder.Encode(msg)
  }
}