package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type CommanderManager struct {
	listenPort int
	listenHost string
}

func (c *CommanderManager) start() {
	listenAddress := fmt.Sprintf("%s:%d", c.listenHost, c.listenPort)
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("Commander server listening on", listenAddress)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go c.handleCommander(conn)
	}
}

func (c *CommanderManager) handleCommander(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Commander connected:", conn.RemoteAddr())

	writer := bufio.NewWriter(conn)

	// Read the message from the commander containing the destination client and command
	reader := bufio.NewReader(conn)
	command, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading command from commander:", err)
		return
	}

	c.handleCommand(command, conn, writer)
}

func (c *CommanderManager) handleCommand(command string, conn net.Conn, writer *bufio.Writer) {
	command = strings.TrimSpace(command)
	if command == LIST {
		c.listClients(writer)
		writer.Flush()
		return
	}
	parts := strings.Split(command, "|")

	destinationClient := parts[0]
	// Find the destination client
	var destinationConn = findClient(destinationClient)
	if destinationConn == nil {
		writer.Write([]byte("No such client"))
		writer.Flush()
		return
	}

	commandToExecute := parts[1]
	commandToExecute = EXECUTE_COMMAND + DELIMITER + commandToExecute

	// Send the command to the destination client
	output, err := c.sendToClient(destinationConn, commandToExecute)
	if err != nil {
		writer.Write([]byte("Failed Sending Command To Client"))
		writer.Flush()
	}

	// Send the output back to the commander
	_, err = conn.Write([]byte(*output))
	if err != nil {
		fmt.Println("Error sending output to commander:", err)
		return
	}
}

func (c *CommanderManager) listClients(writer *bufio.Writer) {
	if len(clients) == 0 {
		writer.WriteString("No clients connected\n")
		return
	}

	var clientsToRemove []*Client
	for _, client := range clients {
		if !c.isAlive(client.conn) {
			clientsToRemove = append(clientsToRemove, client)
		}
	}

	for _, clientToRemove := range clientsToRemove {
		removeClient(clientToRemove.conn)
	}

	for _, client := range clients {
		_, err := writer.WriteString(client.conn.RemoteAddr().String() + "\n")
		if err != nil {
			fmt.Println("Error writing to commander:", err)
			break
		}
	}
}

func (c *CommanderManager) sendToClient(destinationConn *net.Conn, commandToExecute string) (*string, error) {
	_, err := (*destinationConn).Write([]byte(commandToExecute + "\n"))
	if err != nil {
		fmt.Println("Error sending command to client:", err)
		return nil, err
	}

	// Read the output of the command from the destination client and send it back to the commander
	output, err := bufio.NewReader(*destinationConn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading output from client:", err)
		removeClient(*destinationConn)
		return nil, err
	}
	return &output, nil
}

func (c *CommanderManager) isAlive(destinationConn net.Conn) bool {
	keepAliveResult, err := c.sendToClient(&destinationConn, KEEP_ALIVE)
	if err != nil {
		return false
	}

	return strings.TrimSpace(*keepAliveResult) == "ALIVE"
}
