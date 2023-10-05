package main

import (
	"fmt"

	"golang.org/x/sys/unix"
)

const port = 9999

func main() {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	check(err)
	err = unix.Bind(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}, Port: port},
	)
	check(err)

	fmt.Printf("Receiver listening on port %d\n", port)
	for {
		buf := make([]byte, 4)
		_, _, err := unix.Recvfrom(fd, buf, 0)
		check(err)
		fmt.Printf("Received: %v %s\n", buf, buf)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}