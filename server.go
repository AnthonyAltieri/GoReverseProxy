package main

import (
	"fmt"
	"net"
	"os"
	_ "./http"
	_ "./routing"
	"reverseproxy/http"
	"reverseproxy/routing"
	"time"
	"math/rand"
)

const (
	SERVER_HOST 		 = "localhost"
	SERVER_PORT 		 = "8080"
	CONNECTION_TYPE  = "tcp"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	routingTable := routing.ParseRoutingTable()
	fmt.Println("Routing Table:")
	routing.PrintRoutingTable(routingTable)

	listener, err := net.Listen(CONNECTION_TYPE, fmt.Sprintf("%s:%s", SERVER_HOST, SERVER_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println(fmt.Sprintf("\nListening on %s:%s", SERVER_HOST, SERVER_PORT))

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}

		fmt.Println("Received connection:", connection)
		go handleRoutine(connection, routingTable)
	}
}

func handleRoutine(connection net.Conn, routingTable routing.RoutingTable) {
	buffer := make([]byte, 1024)
  bufferLength, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	var request http.Request = http.FormatRequest(buffer, bufferLength)
	http.PrintRequest(request)
	fmt.Println()

	var routes *[]routing.Route = routing.FindPath(request.Path, routingTable)

	if routes == nil {
		fmt.Println("Path not found in routing table")
		connection.Close()
	}



	connection.Write([]byte("Message Received"))
	connection.Close()
}
