package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
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

	fmt.Println(lines[0])
	fmt.Println(lines[1])
	fmt.Println(lines[2])

	if requestTarget == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	if len(requestTarget) > 5 && requestTarget[:5] == "/echo" {
		var bodyString string = requestTarget[6:]
		var len int = len(bodyString)
		var res string = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len, bodyString)
		conn.Write([]byte(res))
		return
	}

	if requestTarget == "/user-agent" {
		var userAgentInfo string = ""

		for _, line := range lines {
			if strings.HasPrefix(line, "User-Agent: ") {
				userAgentInfo = line[12:]
				break
			}
		}

		var len int = len(userAgentInfo)
		var res string = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len, userAgentInfo)

		conn.Write([]byte(res))
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

		handleConnection(conn)
	}
}
