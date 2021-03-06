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

func main() {
	//Scan user input for host:port
	arguments := os.Args
	if len(arguments) != 3 {
		fmt.Println("Please provide host:port and the username.")
		return
	}
	fmt.Print("Type EXIT if you want to leave.\n")

	//connect to provided host:post via the net library
	c := TCPDial(arguments)
	channel := make(chan string)

	go listen(c, channel)

	//Wait for input from channel or user input, if "EXIT" command is received, then terminate.
	//If not signal is received, then get user input.
	for {
		select {
		case signal := <-channel:
			if signal == "EXIT" {
				return
			}
		default:
			sender, receiver, content := getUserInput()
			if content == "EXIT" {
				return
			}
			msg := utils.Message{sender, receiver, content}
			messaging(msg, c)
		}
	}
}

//The listen function waits on messages from the TCP server.
func listen(c net.Conn, channel chan string) {
	for {
		decoder := gob.NewDecoder(c) //initialize gob decoder
		//Decode message struct and print it
		message := new(utils.Message)
		_ = decoder.Decode(message)

		if (*message == utils.Message{"error", "error", "error"}) {
			//If the server discovers that the receiver is not connected, print error
			fmt.Printf("\nError: the person you are sending to has not been connected yet.\n")
			fmt.Printf("Sender: ")
		} else if (*message == utils.Message{"EXIT", "EXIT", "EXIT"}) {
			//If the server terminates, the functions sends the exit signal through the channel to the main thread
			c.Close()
			os.Exit(0)
			channel <- "EXIT"
		} else if message.Content != "" {
			//If everything works properly, the message will be displayed on the screen
			fmt.Printf("Received message from %q\nMessage: %s\n", strings.TrimSpace(message.Sender), strings.TrimSpace(message.Content))
			fmt.Printf("Sender: ")
		}
	}
}

//Function is used to encode and send messages across TCP using gob.
func messaging(msg utils.Message, c net.Conn) {
	//create a gob encoder and code the message struct
	encoder := gob.NewEncoder(c)
	_ = encoder.Encode(msg)
}

//Function is used to connect to the TCP channel using the net library.
func TCPDial(arguments []string) (c net.Conn) {
	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	username := arguments[2]
	fmt.Fprintf(c, username+"\n")

	return c
}

//Gets the sender, receiver and message content from the user input
func getUserInput() (sender string, receiver string, content string) {
	reader := bufio.NewReader(os.Stdin)

	//scan user input for message contents
	//If the user input "exit", the function returns
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

	return sender, receiver, content
}
