package main

const (
	EXECUTE_COMMAND = "EXEC"
	KEEP_ALIVE      = "KEEP_ALIVE"
	LIST            = "LIST"
	DELIMITER       = "|@|"
)

func main() {
	// Start the server on port 8888 for clients
	clientManager := ClientManager{listenHost: "", listenPort: 8888}
	go clientManager.start()

	// Start the server on port 9999 for commanders
	commanderManager := CommanderManager{listenHost: "", listenPort: 9999}
	commanderManager.start()
}
