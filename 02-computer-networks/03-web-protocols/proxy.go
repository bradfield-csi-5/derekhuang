package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"syscall"
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
	// open public tcp socket
	pubSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err)
	defer syscall.Close(pubSock)

	// bind to public port
	err = syscall.Bind(
		pubSock,
		&syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: publicPort},
	)
	check(err)

	// listen to public socket
	err = syscall.Listen(pubSock, 1)
	check(err)

	fmt.Printf("Listening on public port %d\n", publicPort)

	for {
		// accept from public socket
		fd, _, err := syscall.Accept(pubSock)
		check(err)
		go forward(fd)
	}
}

func forward(fd int) {
	defer syscall.Close(fd)

	// decode what was received from the public socket
	for {
		buf := make([]byte, 4096)

		// receive from public socket
		bytesRead, _, err := syscall.Recvfrom(fd, buf, 0)
		check(err)

		req := decode(buf[:bytesRead])
		fmt.Printf("req.method: %s\n", req.method)
		fmt.Printf("req.uri: %s\n", req.uri)
		fmt.Printf("req.protocol: %s\n", req.protocol)
		for k, v := range req.headers {
			fmt.Printf("headers[%s]: %s\n", k, v)
		}

		// TODO: make work
		// // open server tcp socket
		// serverSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		// check(err)
		// defer syscall.Close(serverSock)
		//
		// // bind to server port
		// err = syscall.Bind(serverSock, &syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: serverPort})
		// check(err)
		//
		// // listen to server socket
		// err = syscall.Listen(serverSock, 1)
		// check(err)
		//
		// fmt.Printf("Listening on server port %d\n", serverPort)
		//
		// // accept from server socket
		// _, servSa, err := syscall.Accept(serverSock)
		// check(err)
		//
		// fmt.Printf("Sending %v to server\n", buf)
		// // forward request to server
		// resp := "HTTP/1.1 200 OK\n"
		// resp += "Content-Length: %d\n"
		// resp += "Content-Type: text/plain\n\n%s"
		//
		// err = syscall.Sendto(
		// 	serverSock,
		// 	[]byte(fmt.Sprintf(resp, pubBytesRead, pubBuf[:pubBytesRead])),
		// 	0,
		// 	servSa,
		// )
		// check(err)
	}
}

func decode(received []byte) *Request {
	parts := bytes.SplitN(received, []byte{'\r', '\n', '\r', '\n'}, 2)
	lines := bytes.Split(parts[0], []byte{'\n'})
	first := bytes.Split(lines[0], []byte{' '})
	var req Request
	req.method = string(first[0])
	req.uri = string(first[1])
	req.protocol = string(first[2])
	req.headers = make(map[string]string)
	for i := 1; i < len(lines); i++ {
		lineParts := strings.Split(string(lines[i]), ": ")
		req.headers[strings.ToLower(string(lineParts[0]))] = strings.ToLower(string(lineParts[1]))
	}
	return &req
}

func check(e error) {
	if e != nil {
		log.Fatalln("check caught error:", e)
	}
}
