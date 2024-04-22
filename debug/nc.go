package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var host string
	var port int
	var listen bool

	flag.StringVar(&host, "host", "", "Hostname or IP address to connect to")
	flag.IntVar(&port, "port", 0, "Port number to connect to or listen on")
	flag.BoolVar(&listen, "listen", false, "Listen for incoming connections")
	flag.Parse()

	if listen {
		listenAndServe(port)
	} else {
		connectAndInteract(host, port)
	}
}

func listenAndServe(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on port %d...\n", port)

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connection from %s established\n", conn.RemoteAddr())

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := conn.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading data:", err)
		return
	}
}

func connectAndInteract(host string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := conn.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading data:", err)
		return
	}
}