package main

import (
	"fmt"
	"net"
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

type ClientManager struct {
	listenPort int
	listenHost string
}

func (cm *ClientManager) start() {
	listenAddress := fmt.Sprintf("%s:%d", cm.listenHost, cm.listenPort)
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("Client server listening on", listenAddress)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go cm.handleClient(conn)
	}
}

func (cm *ClientManager) handleClient(conn net.Conn) {
	clientsMu.Lock()
	clients = append(clients, &Client{conn: conn})
	clientsMu.Unlock()

	fmt.Println("Client connected:", conn.RemoteAddr())
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

func findClient(destinationClient string) *net.Conn {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	var destinationConn net.Conn
	for _, client := range clients {
		if client.conn.RemoteAddr().String() == destinationClient {
			destinationConn = client.conn
			break
		}
	}

	if destinationConn == nil {
		fmt.Println("Destination client not found")
		return nil
	}

	return &destinationConn
}
