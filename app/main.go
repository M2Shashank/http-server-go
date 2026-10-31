package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func do(conn net.Conn) {
	var buff []byte = make([]byte, 1024)

	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		conn.Close()
		return
	}

	var request string = string(buff[:])
	var lines []string = strings.Split(request, "\r\n")
	var requestLine string = lines[0]
	var parts []string = strings.Fields(requestLine)
	var requestTarget string = parts[1]

	var response []byte
	if requestTarget == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	conn.Write(response)
	conn.Close()
}

func main() {
	ls, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer ls.Close()
	fmt.Println("Server is listening on port 4221...")

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		do(conn)
	}
}
