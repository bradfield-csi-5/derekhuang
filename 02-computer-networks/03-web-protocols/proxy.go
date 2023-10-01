package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	publicPort = 8000
	serverPort = 9000
)

type Request struct {
	method   string
	uri      string
	protocol string
	headers  map[string]string
	body     []byte
}

func main() {
	cache := make(map[string][]byte)

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	check("error with opening socket to client:", err)
	err = unix.Bind(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: publicPort},
	)
	check("error with Bind:", err)
	err = unix.Listen(fd, 5)
	check("error with Listen:", err)

	fmt.Printf("Listening on public port %d\n", publicPort)
	for {
		// accept from socket
		clientFd, _, err := unix.Accept(fd)
		check("error with Accept:", err)
		defer unix.Close(fd)

		// receive packet
		clientBuf := make([]byte, 4096)
		n, _, err := unix.Recvfrom(clientFd, clientBuf, 0)
		check("error with Recvfrom:", err)
		if n == 0 {
			fmt.Printf("No bytes received from the client. Connection closed.\n")
		} else {
			fmt.Printf("Received %d bytes from client\n", n)
		}

		// parse assuming the request is http
		req := string(clientBuf)
		parts := strings.Split(req, " ")
		path := parts[1]

		cached, exists := cache[path]
		if exists {
			fmt.Printf("Cache hit! Skipping server request and responding...\n\n")

			// send back
			err = unix.Sendto(clientFd, cached, 0, nil)
			check("error responding to client:", err)
		} else {
			fmt.Println("Cache miss. Fetching from server...")

			// create a new socket for the server connection
			serverSock, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
			check("error with opening socket to server:", err)
			defer unix.Close(serverSock)

			// connect to the server port
			err = unix.Connect(
				serverSock,
				&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: serverPort},
			)
			check("error connecting with server:", err)

			// send request to the server
			err = unix.Sendto(serverSock, clientBuf, 0, nil)
			check("error sending to server:", err)

			// receive packets from the server
			serverBuf := make([]byte, 4096)

			// buffer for chunked data
			buf := make([]byte, 0)
			for {
				n, _, err = unix.Recvfrom(serverSock, serverBuf, 0)
				check("error recvfrom server:", err)
				if n > 0 {
					fmt.Printf("  -> received %d bytes from server\n", n)
					buf = append(buf, serverBuf...)
					err = unix.Sendto(clientFd, serverBuf, 0, nil)
					check("error sending to client:", err)
				} else {
					fmt.Printf(
						"    -> server finished sending. Inserting path '%s' into cache...\n\n",
						path,
					)
					cache[path] = buf
					break
				}
			}
		}
	}
}

func check(prefix string, e error) {
	if e != nil {
		log.Fatalln(prefix, e)
	}
}
