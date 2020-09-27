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

Once the clients started, there will be input instructions for sending a text or exit the chatroom. The sender should be the same with the username. For example:
```bash
Type EXIT if you want to leave. Press anything else to continue. Space

Sender: Jack
Receiver: Alice
Message content: Hi
```
If the receiver has not been started yet, the sender client should see the following commands:
```bash
Error: the person you are sending to has not been connected yet.
Type EXIT or enter Sender: 

```

If the receiver is connected successfully, the receiver will see the message displayed in the terminal. For example:
```bash
Received message from Jack

Message: hi

Type EXIT or enter Sender: 
```

The server can exit with the command line input "EXIT", and it will send termination signal to all the connected clients to make them exit. (TODO)
```bash
Waiting for exit command...
Listening on port:6000
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
The server has a thread for the TCP listening, and inside the thread it uses client thread to handle the communication with each client.
The server has a separate thread that waits for an "EXIT" command.

In the client program, each client uses a goroutine as a thread. 

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

We also abstract helper functions such as userExit() for the clients to exit, exit() for the server to exit, and getUserInput(). 

## Resources
* [TCP Concurrent Server](https://www.linode.com/docs/development/go/developing-udp-and-tcp-clients-and-servers-in-go/)
* [Gob](https://golang.org/pkg/encoding/gob/)
## Authors
* Jiahong Li
* Zheng Zhou