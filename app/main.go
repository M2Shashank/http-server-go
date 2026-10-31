package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func do(conn net.Conn) {
	defer conn.Close()
	var buff []byte = make([]byte, 1024)

	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var request string = string(buff[:n])
	var lines []string = strings.Split(request, "\r\n")
	var requestLine string = lines[0]
	var parts []string = strings.Fields(requestLine)
	var requestTarget string = parts[1]

	if requestTarget == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if len(requestTarget) > 5 && requestTarget[:5] == "/echo" {
		var bodyString string = requestTarget[7:]
		var len int = len(bodyString)
		var res string = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + string(len) + "\r\n\r\n" + bodyString
		conn.Write([]byte(res))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
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
