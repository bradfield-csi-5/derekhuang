package main

import (
	"fmt"
	"log"
	"syscall"
)

const (
	port = 8000
)

func main() {
	// open tcp socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	check(err)
	defer syscall.Close(fd)

	// bind to port
	err = syscall.Bind(fd, &syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: port})
	check(err)

	// listen
	err = syscall.Listen(fd, 1)
	check(err)

	fmt.Printf("Listening on port %d\n", port)

	// accept
	buf := make([]byte, 2048)
	for {
		connfd, connsa, err := syscall.Accept(fd)
		check(err)
		defer syscall.Close(connfd)

		bytesRead, _, err := syscall.Recvfrom(connfd, buf, 0)
		check(err)

		if bytesRead == 0 {
			fmt.Println("client closed connection")
			return
		}

		err = syscall.Sendto(connfd, buf[:bytesRead], 0, connsa)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		log.Fatalln("check caught error:", e)
	}
}
