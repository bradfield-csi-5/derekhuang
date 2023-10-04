package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	publicPort = 8000
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run client.go [proxy_port]")
	}

	proxyPort, err := strconv.Atoi(os.Args[1])
	check(err)

	// open a socket to the proxy
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	check(err)
	// err = unix.Bind(
	// 	fd,
	// 	&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: publicPort},
	// )
	// check("error binding to client port:", err)
	err = unix.Connect(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: proxyPort},
	)

	fmt.Printf("Listening on port %d\n", publicPort)
	for {
		clientBuf := make([]byte, 64)
		n, _, err := unix.Recvfrom(fd, clientBuf, 0)
		check(err)
		if n == 0 {
			fmt.Printf("No bytes received. Connection Closed.\n")
		} else {
			fmt.Printf("Received %d bytes from client\n", n)
		}

		// connect with the proxy server
		proxyFd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
		check(err)
		err = unix.Connect(
			proxyFd,
			&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: proxyPort},
		)
		check(err)

		err = unix.Sendto(proxyFd, clientBuf, 0, nil)
		check(err)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
