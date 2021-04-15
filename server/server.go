package main

import (
	"fmt"
	"net"
	"bufio"
)

const (
	host = "localhost"
	port = 8080
	packetType = "tcp"
)

var	usernames map[string]string = make(map[string]string)
var connections []net.Conn = make([]net.Conn, 0, 10)

func main() {
	fmt.Printf("Starting server. Details:\n- Host: %s\n- Port: %d\n- Type: %s\n", host, port, packetType)

	// create the tcp socket on our host/port
	socket, err := net.Listen(packetType, fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("Error starting %s server (%s)\n", packetType, fmt.Sprintf("%s:%d", host, port))
		return
	}

	defer socket.Close() // ensure we close the socket

	for {
		client, err := socket.Accept() // accept a new client
		if err != nil {
			fmt.Printf("Error accepting client (%s)\n", err.Error())
			return
		}
		
		fmt.Printf("%s connected to server\n", client.RemoteAddr().String())
		
		go handleConnection(client)	
	}
}

func notifyAllClients(msg string) {
	for _, client := range connections {
		client.Write([]byte(msg))
	}
}

func handleConnection(client net.Conn) {
	client.Write([]byte("[SERVER]: Please enter a username: ")) // prompt the user for a username

	buffer, err := bufio.NewReader(client).ReadBytes('\n')

	if err != nil { // client disconnected
		return
	}

	usernames[client.RemoteAddr().String()] = string(buffer[:len(buffer)-1])
	connections = append(connections, client)	

	notifyAllClients(fmt.Sprintf("[+] %s (%s) connected.\n", usernames[client.RemoteAddr().String()], client.RemoteAddr().String()))

	go handleClient(client)
}

func handleClient(client net.Conn) {
	buffer, err := bufio.NewReader(client).ReadBytes('\n')
	
	if err != nil { // client disconnected
		fmt.Printf("%s disconnected from server\n", client.RemoteAddr().String())
		notifyAllClients(fmt.Sprintf("[-] %s (%s) disconnected.\n", usernames[client.RemoteAddr().String()], client.RemoteAddr().String()))
		client.Close()		

		return
	}

	fmt.Printf("[%s]: %s\n", usernames[client.RemoteAddr().String()], string(buffer[:len(buffer)-1]))
	notifyAllClients(fmt.Sprintf("[%s]: %s\n", usernames[client.RemoteAddr().String()], string(buffer[:len(buffer)-1])))
		
	handleClient(client)
}
