package main

import (
	"bufio"
	"github.com/gotterdemarung/cfmt"
	"github.com/gotterdemarung/go-noncense/noncense"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	ConnectionType = "tcp"
)

func main() {
	args := os.Args[1:]
	cpu := runtime.NumCPU()
	cfmt.Println()
	cfmt.Println(cfmt.FHeader("NONCEs holder"))

	if len(args) != 3 {
		help(nil, "")
	}

	addr := args[0]

	mapSize, err := strconv.Atoi(args[1])
	if err != nil || mapSize < 2 {
		help(err, "Not valid map size, minimum 2")
	}
	procs, err := strconv.Atoi(args[2])
	if err != nil || mapSize < 1 {
		help(err, "Not valid threads count")
	}

	table := cfmt.TableBuffer2{}
	table.Add("Number of CPU", cpu)
	table.Add("Number of threads", procs)
	table.Add("Amount of NONCE", mapSize)
	table.Add("HTTP address", addr)
	cfmt.Println(table)

	runtime.GOMAXPROCS(procs)

	l, err := net.Listen(ConnectionType, addr)
	if err != nil {
		help(err, "Error listening")
	}
	defer l.Close()

	adder := noncense.NewNoncesAdder(uint32(mapSize))
	servedConnections := 0

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		servedConnections++
		if err != nil {
			help(err, "Error accepting")
		}
		if (servedConnections % 10000) == 0 {
			cfmt.Println(
				cfmt.FString("Served"),
				" ",
				cfmt.FInt(servedConnections),
			)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, adder)
	}
}

// Prints help message
func help(err error, errDescripion string) {
	cfmt.Println(
		cfmt.FString("app"),
		" ",
		"<hostname:port>",
		" ",
		"<map size>",
		" ",
		"<threads>",
	)

	if err == nil && errDescripion == "" {
		cfmt.Println()
		cfmt.Println()
		cfmt.Println(cfmt.FHeader("Usage example"))
		cfmt.Println("./noncense-server localhost:8888 1000000 4")
		cfmt.Println("Runs server on local port 8888. Holds one million of NONCEs")
		cfmt.Println("and processes all incoming requests in 4 threads")
		cfmt.Println()
	}

	if errDescripion != "" {
		cfmt.Println(cfmt.FError(errDescripion))
	}
	if err != nil {
		cfmt.Println(cfmt.FError(err))
	}

	os.Exit(1)
}

// Handles incoming requests.
func handleRequest(conn net.Conn, adder *noncense.NonceAdder) {
	message, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		cfmt.Println(cfmt.FError("Error reading"))
		cfmt.Println(cfmt.FError(err))
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
