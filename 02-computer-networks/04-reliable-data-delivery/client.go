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
	check("atoi:", err)

	// open a socket to the client (nc)
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	check("error with opening client socket:", err)
	err = unix.Bind(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: publicPort},
	)
	check("error binding to client port:", err)

	fmt.Printf("Listening on port %d\n", publicPort)
	for {
		clientBuf := make([]byte, 64)
		n, _, err := unix.Recvfrom(fd, clientBuf, 0)
		check("error setting up recvfrom:", err)
		if n == 0 {
			fmt.Printf("No bytes received. Connection Closed.\n")
		} else {
			fmt.Printf("Received %d bytes from client\n", n)
		}

		// connect with the proxy server
		proxyFd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
		check("error with opening proxy socket:", err)
		err = unix.Connect(
			proxyFd,
			&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: proxyPort},
		)
		check("error connecting to proxy:", err)

		err = unix.Sendto(proxyFd, clientBuf, 0, nil)
		check("error sending to proxy:", err)
	}

}

func check(prefix string, e error) {
	if e != nil {
		log.Fatalln(prefix, e)
	}
}
