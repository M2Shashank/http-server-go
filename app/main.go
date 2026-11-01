package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"io"
	"strconv"
	"strings"
	"bufio"
)

type Request struct {
	method, target, version string
	headers map[string]string
	body []byte
}

func parseRequest(reader *bufio.Reader) (Request, error) {
	var request Request = Request{
		headers: make(map[string]string),
	}

	requestStatus, err := reader.ReadString('\n')
	if err != nil {
		return request, fmt.Errorf("Error reading request: %s", err.Error())
	}
	var parts []string = strings.Fields(requestStatus)
	request.method = parts[0]
	request.target = parts[1]
	request.version = parts[2]

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return request, fmt.Errorf("Error reading request: %s", err.Error())
		}


		if line == "\r\n" {
			break
		}

		key, value, ok := strings.Cut(line, ":")
        if !ok { continue }

        request.headers[strings.TrimSpace(key)] =
            strings.TrimSpace(value)
	}

	var n, _ = strconv.Atoi(request.headers["Content-Length"])
	request.body = make([]byte, n)

	_, err = io.ReadFull(reader, request.body)
	if err != nil {
		return request, fmt.Errorf("Error reading request body: %s", err.Error())
	}

	return request, nil
}

type Response struct {
	statusCode int
	version, status string
	headers map[string]string
	body []byte
}

func (r Response) BuildResponse() string {
	var response strings.Builder = strings.Builder{}

	var firstline string = fmt.Sprintf("%s %d %s\r\n", r.version, r.statusCode, r.status)
	response.WriteString(firstline)

	for key, value := range r.headers {
		var headerLine string = fmt.Sprintf("%s: %s\r\n", key, value)
		response.WriteString(headerLine)
	}
	response.WriteString("\r\n")
	response.Write(r.body)

	return response.String()
}

var directory = flag.String(
	"directory",
	"./",
	"Directory to serve files from (default: current directory)",
)

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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var reader *bufio.Reader = bufio.NewReader(conn)

	for {

		request, err := parseRequest(reader)
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			return
		}

		var response Response = Response{
			version: "HTTP/1.1",
			headers: make(map[string]string),
		}

		var closeConn bool = request.headers["Connection"] == "close"
		if closeConn {
			response.headers["Connection"] = "close"
		}

		if request.target == "/" {
			response.statusCode = 200
			response.status = "OK"

			conn.Write([]byte(response.BuildResponse()))

			if closeConn { return }
			continue
		}

		if strings.HasPrefix(request.target, "/echo/") {
			var query string = request.target[6:]

			response.statusCode = 200
			response.status = "OK"
			response.body = []byte(query)
			response.headers["Content-Type"] = "text/plain"
			response.headers["Content-Length"] = strconv.Itoa(len(query))
			var encoding string = request.headers["Accept-Encoding"]
			if strings.Contains(encoding, "gzip") {
				response.headers["Content-Encoding"] = "gzip"
			}

			conn.Write([]byte(response.BuildResponse()))
			
			if closeConn { return }
			continue
		}

		if request.target == "/user-agent" {
			var userAgentInfo string = request.headers["User-Agent"]

			response.statusCode = 200
			response.status = "OK"
			response.body = []byte(userAgentInfo)
			response.headers["Content-Type"] = "text/plain"
			response.headers["Content-Length"] = strconv.Itoa(len(userAgentInfo))

			conn.Write([]byte(response.BuildResponse()))

			if closeConn { return }
			continue
		}

		if strings.HasPrefix(request.target, "/files/") {
			var filePath string = (*directory) + request.target[7:]

			switch request.method {
				case "GET":
				{
					fmt.Println("Serving file:", filePath)

					data, err := os.ReadFile(filePath)
					if err != nil {
						response.statusCode = 404
						response.status = "Not Found"
						conn.Write([]byte(response.BuildResponse()))

						if closeConn { return }
						continue
					}

					response.statusCode = 200
					response.status = "OK"
					response.body = data
					response.headers["Content-Type"] = "application/octet-stream"
					response.headers["Content-Length"] = strconv.Itoa(len(data))

					conn.Write([]byte(response.BuildResponse()))
				}
				case "POST":
				{
					fmt.Println("Writing to file:", filePath)

					var data []byte = request.body

					var err1 error = os.MkdirAll(*directory, 0755)
					var err2 error = os.WriteFile(filePath, data, 0644)
					if err1 != nil || err2 != nil {
						response.statusCode = 500
						response.status = "Internal Server Error"
						conn.Write([]byte(response.BuildResponse()))
						
						if closeConn { return }
						continue
					}

					response.statusCode = 201
					response.status = "Created"

					conn.Write([]byte(response.BuildResponse()))
				}
			}

			if closeConn { return }
			continue
		}

		response.statusCode = 404
		response.status = "Not Found"

		conn.Write([]byte(response.BuildResponse()))

		if closeConn { return }

	}
}