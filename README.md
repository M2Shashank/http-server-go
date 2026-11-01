# HTTP Server in Go

A small HTTP/1.1 server built from scratch in Go as part of the
CodeCrafters challenge.

The server uses TCP sockets directly, parses HTTP requests, supports multiple
clients concurrently, and serves files from a configurable directory.

## Features

- HTTP/1.1 request parsing
- Concurrent client handling with goroutines
- Persistent connections and `Connection: close`
- `GET /`
- `GET /echo/{text}`
- `GET /user-agent`
- `GET /files/{filename}`
- `POST /files/{filename}`
- Optional gzip compression for echo responses
- Configurable file directory

## Project structure

```text
.
├── app/
│   └── main.go       # HTTP server implementation and entry point
├── hello.txt         # Sample file
├── go.mod
├── go.sum
├── your_program.sh   # Build and run locally
└── README.md
```

## Requirements

- Go 1.27 or newer

## Run locally

Clone the repository and start the server:

```sh
./your_program.sh
```

The server listens on `0.0.0.0:4221`.

On Windows, run the script through a POSIX-compatible shell, or build and run
the application directly:

```sh
go build -o http-server.exe ./app
.\http-server.exe
```

To serve files from a specific directory:

```sh
./your_program.sh --directory ./data
```

## Endpoints

### Root

```sh
curl http://localhost:4221/
```

### Echo

```sh
curl http://localhost:4221/echo/hello
```

The server returns `hello`. If the client advertises gzip support, the
response is gzip-compressed:

```sh
curl --compressed http://localhost:4221/echo/hello
```

### User-Agent

```sh
curl -A "ExampleClient/1.0" http://localhost:4221/user-agent
```

### Files

Read a file:

```sh
curl http://localhost:4221/files/hello.txt
```

Create or overwrite a file:

```sh
curl -X POST --data "Hello from the HTTP server" \
  http://localhost:4221/files/example.txt
```

The file endpoint uses the directory provided through `--directory`. By
default, it uses the current working directory.

## CodeCrafters

Run the CodeCrafters test suite with:

```sh
codecrafters submit
```

This project follows the CodeCrafters program interface, which builds the
application from `app/main.go`.
