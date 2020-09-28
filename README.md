# MP2

## Usage

In one terminal, start the chatroom server and enter a port number as an argument.

```bash
go run Server.go 6000
```

In other terminals, start the chatroom clients and enter host:port number and username as arguments.
```bash
go run Client.go 127.0.0.1:6000 Jack
```

```bash
go run Client.go 127.0.0.1:6000 Alice
```

At this point, you should be prompted to input information:
```bash
Type EXIT if you want to leave.

Sender: Jack
Receiver: Alice
Message content: Hi
```
If the receiver has not been started yet, the sender client should see the following commands:
```bash
Error: the person you are sending to has not been connected yet.
Sender: 

```

If the receiver is connected successfully, the receiver will see the message displayed in the terminal. For example:
```bash
Received message from Jack

Message: hi

Sender: 
```

The server can exit with the command line input "EXIT", and it will send termination signal to all the connected clients to make them exit.
```bash
Listening on port:6000. Please type 'EXIT' to quit.
EXIT
Server is exiting...
```

The client can also exit the chatroom when they receive "EXIT" from user command. For example:
```bash
Type EXIT if you want to leave. Press anything else to continue.
EXIT
Client "Alice" is exiting...
```
## Structure and Design
#### There are two layers in our design: application layer and the network layer:


### 1) Network Layer 
* Our networking layer is located in `Server.go`. We start two goroutines within main.go. The first one is called `startServer()` and the second is `exit()`

* `startServer()` creates the TCP connection and listen on provided port for requests. It stores the username of the client's in a map as a key, and the channel as the associating value. We use a map so that we are able to know which channel we have to redirect messages to(message's receiver). For each connection, we call a goroutine to handle the communication so that the server can support handling multiple concurrent clients. The goroutine will encode the message on the receiverchannel if the receiver is connected and send error message otherwise.    

* `exit()` is used to wait for an EXIT command from the command line. If this command is detected, then a channel is used to communicate this information with the main thread. Then the main thread sends the signal to all other TCP channels to terminates those, and finally terminates itself.


### 2) Application Layer 
* Our application code is located in `Client.go`.

* `Client.go` starts by dialing via TCP to the ip and port included in the user input. After this is done, it calls a go routine called `listen()`, each go routine representing a client. 

* `listen()` is used as a goroutine so that it can handling incoming requests from the server concurrently. It also has communicates with the main thread using a channel in case the server sends the "EXIT" signal.

* If the clients receive the "EXIT" signal from the server, then all of them will terminate.

* The rest of the application is reading user input from `getUserInput()`, which fetches the "Sender,Receiver, and Content" for each message. If "EXIT" is received from the user input, then this function sends a signal to the main function to terminate this process.

* Once the user input is received for a message, the message is converted into a struct `Message` and sent through the TCP channel through gob in a function called `messaging()`.

* Once the client receives a message, the sender and the message will be displayed on the screen.

### Message
```bash
	Sender      string
	Receiver    string
	Content     string
```

### File Structure and Abstraction
* Our project is broken down in the following structure:

```bash
- server.go
- client.go
- utils
   - Message.go
```

We also abstract helper functions such as `exit()` for the server to exit, `exitAllClients()` to send termination signal to all clients in the map, `messaging()` for encoding the message, `getUserInput()` to get sender, receiver, and message from the client. 

## Resources
* [TCP Concurrent Server](https://www.linode.com/docs/development/go/developing-udp-and-tcp-clients-and-servers-in-go/)
* [Gob](https://golang.org/pkg/encoding/gob/)
* [Select](https://gobyexample.com/select)
* Sean's group for the idea to use map for storing channels and usernames.
## Authors
* Jiahong Li
* Zheng Zhou