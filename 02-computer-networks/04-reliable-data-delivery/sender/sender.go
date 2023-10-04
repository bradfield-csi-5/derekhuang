package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: go run sender.go [proxy_port]")
	}

	proxyPort, err := strconv.Atoi(os.Args[1])
	check(err)

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	check(err)
	err = unix.Connect(
		fd,
		&unix.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: proxyPort},
	)
	check(err)

	fmt.Printf("Sender connected to proxy port %d\n", proxyPort)
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
