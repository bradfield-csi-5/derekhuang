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

	// start ack and seqnum at 0
	ack := false
	// seqnum := false

	fmt.Printf("Receiver listening on port %d\n", port)
	for {
		buf := make([]byte, 4)
		_, _, err := unix.Recvfrom(fd, buf, 0)
		check(err)
		fmt.Printf("Received: %b %s\n", buf, buf)
		if isACK(buf[0], ack) {
			// send ack
		} else {
			// send nak
		}
	}
}

func isACK(header byte, expected bool) bool {
	ack := header & 0x80
	if expected {
		return ack == 1
	}
	return ack == 0
}

func mkhdr(ack bool, seqnum bool) byte {
	header := byte(0)
	if ack {
		header |= 0x80
	}
	if seqnum {
		header |= 0x40
	}
	return header
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
