package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	publicPort = 8000
	serverPort = 9000
)

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
		defer unix.Close(clientFd)

		// receive packet
		clientBuf := make([]byte, 4096)
		n, _, err := unix.Recvfrom(clientFd, clientBuf, 0)
		check("error with Recvfrom:", err)
		if n == 0 {
			fmt.Printf("No bytes received from the client. Connection closed.\n")
		} else {
			fmt.Printf("Received %d bytes from client\n", n)
		}

		// parse with the assumption the request is http
		// change connection to keep-alive before creating the http.Request
		// since headers can't be updated on an existing request
		// NOTE: the assumption this is a complete http request doesn't work
		// for the concurrent test
		clientReqStr := string(clientBuf[:n])
		clientReqStr = strings.Replace(
			clientReqStr,
			"Connection: close",
			"Connection: Keep-Alive",
			1,
		)
		httpReq, err := parse([]byte(clientReqStr))
		check("error parsing client req:", err)

		cached, exists := cache[httpReq.URL.String()]
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
			clientReq, err := encode(httpReq)
			check("error encoding client request:", err)
			err = unix.Sendto(serverSock, clientReq, 0, nil)
			check("error sending to server:", err)

			// receive packets from the server
			serverBuf := make([]byte, 4096)

			buf := make([]byte, 0)
			for {
				n, _, err = unix.Recvfrom(serverSock, serverBuf, 0)
				check("error recvfrom server:", err)
				if n > 0 {
					fmt.Printf("  -> received %d bytes from server\n", n)
					buf = append(buf, serverBuf[:n]...)
					err = unix.Sendto(clientFd, serverBuf[:n], 0, nil)
					check("error sending to client:", err)
				} else {
					fmt.Printf(
						"    -> server finished sending. Inserting path '%s' into cache...\n\n",
						httpReq.URL.String(),
					)
					cache[httpReq.URL.String()] = buf
					break
				}
			}
		}
	}
}

func parse(reqBytes []byte) (*http.Request, error) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(reqBytes)))
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encode(req *http.Request) ([]byte, error) {
	var buffer bytes.Buffer
	err := req.Write(&buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func check(prefix string, e error) {
	if e != nil {
		log.Fatalln(prefix, e)
	}
}
