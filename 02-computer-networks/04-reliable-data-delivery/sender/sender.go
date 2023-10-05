package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/unix"
)

const port = 12345

func main() {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	check(err)
	err = unix.Connect(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: port},
	)
	check(err)

	fmt.Printf("Sender connected to proxy port %d\n", port)
	for i := 0; i < 10; i++ {
		err = unix.Sendto(fd, []byte("bar"), 0, nil)
		check(err)
		time.Sleep(3 * time.Second)
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
