package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
  "strings"
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
  channel := make(chan string)

  go listen(c, channel)

  for {
    select {
      case signal := <-channel:
        if (signal == "EXIT") {
          return
        }
      default:
        sender, receiver, content := getUserInput()
        if (content == "EXIT") {
          return
        }
        msg := utils.Message{sender, receiver, content}
        messaging(msg, c)
      }
  }
}

func listen(c net.Conn, channel chan string){
  for {
	  decoder := gob.NewDecoder(c) //initialize gob decoder
	  //Decode message struct and print it
	  message := new(utils.Message)
	  _ = decoder.Decode(message)

    if(*message == utils.Message{"error","error","error"}) {
      fmt.Printf("\nError: the person you are sending to has not been connected yet.\n")
      fmt.Printf("Sender: ")
    } else if (*message == utils.Message{"EXIT","EXIT","EXIT"}) {
      c.Close()
      channel <- "EXIT"
    } else if (message.Content != "") {
      fmt.Printf("Received message from %q\nMessage: %s\n", strings.TrimSpace(message.Sender), strings.TrimSpace(message.Content))
      fmt.Printf("Sender: ")
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
    if strings.TrimSpace(sender) == "EXIT" {
		  return "EXIT", "EXIT", "EXIT"
	  }

    fmt.Print("Receiver: ")
    receiver, _ = reader.ReadString('\n')
    if strings.TrimSpace(receiver) == "EXIT" {
		  return "EXIT", "EXIT", "EXIT"
	  }

    fmt.Print("Message content: ")
    content, _ = reader.ReadString('\n')
    if strings.TrimSpace(content) == "EXIT" {
		  return "EXIT", "EXIT", "EXIT"
	  }

    return sender,receiver,content
}


	
