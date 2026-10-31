package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type Request struct {
	method, target, version string
	headers                 map[string]string
	body                    []byte
}

func parseRequest(requestStr string) Request {
	var request Request = Request{}
	var lines []string = strings.Split(requestStr, "\r\n")

	var requestLine string = lines[0]
	var parts []string = strings.Fields(requestLine)
	request.method = parts[0]
	request.target = parts[1]
	request.version = parts[2]
	request.headers = make(map[string]string)

	var bodyString strings.Builder = strings.Builder{}
	var isBody bool = false
	for _, line := range lines[1:] {
		if line == "" {
			isBody = true
			continue
		}

		if isBody {
			bodyString.WriteString(line)
			bodyString.WriteString("\r\n")
			continue
		}

		var parts []string = strings.Split(line, ": ")
		if len(parts) == 2 {
			request.headers[parts[0]] = parts[1]
		}
	}

	request.body = []byte(bodyString.String())

	return request
}

var directory = flag.String(
	"directory",
	"",
	"Directory to serve files from (default: current directory)",
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buff []byte = make([]byte, 1024)

	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}

	var requestStr string = string(buff[:n])
	var request Request = parseRequest(requestStr)
	var requestTarget string = request.target

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
		var userAgentInfo string = request.headers["User-Agent"]
		var len int = len(userAgentInfo)

		var res string = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len, userAgentInfo)

		conn.Write([]byte(res))
		return
	}

	if strings.HasPrefix(requestTarget, "/file/") {
		var filePath string = *directory + requestTarget[6:]

		data, err := os.ReadFile(filePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			return
		}

		var responseText string = fmt.Sprintf(
			"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
			len(data), data,
		)

		conn.Write([]byte(responseText))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
}

func main() {
	flag.Parse()

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

		go handleConnection(conn)
	}
}
