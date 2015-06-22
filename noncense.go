package main

import (
    "fmt"
    "net"
    "os"
	"./noncense"
	"runtime"
	"bufio"
	"strings"
	"strconv"
)

const (
	CONN_TYPE = "tcp"
)

func main() {
	args  	:= os.Args[1:]
	cpu 	:= runtime.NumCPU()

	if len(args) != 3 {
		fmt.Println("app <hostname:port> <map size> <threads>\n");
		os.Exit(1)
	}

	addr := args[0]

	mapSize, err := strconv.Atoi(args[1])
	if err != nil || mapSize < 2 {
		fmt.Printf("Not valid map size, minimum 2\n\n")
		os.Exit(2)
	}
	procs, err := strconv.Atoi(args[2])
	if err != nil || mapSize < 1 {
		fmt.Printf("Not valid threads count\n\n")
		os.Exit(2)
	}

	fmt.Println("Starting server")
	fmt.Printf("Number of CPU:       %v\n", cpu)
	fmt.Printf("Number of processes: %v\n", procs)
	fmt.Printf("Amount of NONCE-s:   %v\n", mapSize)
	fmt.Printf("HTTP address:        %v\n\n", addr)

	runtime.GOMAXPROCS(procs)

	l, err := net.Listen(CONN_TYPE, addr);
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
	conn.Close()
}