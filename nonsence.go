package main

import (
    "fmt"
    "net"
    "os"
	"./noncense"
	"runtime"
	"bufio"
	"strings"
)

const (
	CONN_TYPE = "tcp"
)

func main() {

	procs := runtime.NumCPU()
	apiHost := "localhost"
	apiPort := "5542"
	mapSize := 1000000


	fmt.Println("Starting server")
	fmt.Printf("Number of processes: %v\n", procs)
	fmt.Printf("Amount of NONCE-s:   %v\n", mapSize)
	fmt.Printf("HTTP address:        %v:%v\n\n", apiHost, apiPort)

	runtime.GOMAXPROCS(procs)

	l, err := net.Listen(CONN_TYPE, apiHost + ":" + apiPort);
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	adder := noncense.NewNoncesAdder(uint32(mapSize))
	servedConnections := 0;

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		servedConnections++
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		if servedConnections % 10000 == 0 {
			fmt.Printf("Served %v\n", servedConnections)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, adder)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, adder *noncense.NonceAdder) {
	message, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Error reading:", err.Error())
		conn.Write([]byte("0"))
	} else {
		message = strings.TrimSpace(message)

		result := <-adder.Add(message)

		if result {
			conn.Write([]byte("allow"))
		} else {
			conn.Write([]byte("deny"))
		}
	}
	// Close the connection when you're done with it.
	conn.Close()
}