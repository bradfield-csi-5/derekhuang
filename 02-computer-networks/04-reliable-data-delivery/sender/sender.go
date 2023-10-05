package main

import (
	"bytes"
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

	// start ack and seqnum at 0
	ack := true
	seqnum := true

	fmt.Printf("Sender connected to proxy port %d\n", port)
	for i := 0; i < 10; i++ {
		pkt := mkpkt(ack, seqnum, []byte("bar"))
		err = unix.Sendto(fd, pkt, 0, nil)
		check(err)
		time.Sleep(3 * time.Second)
		// TODO: receive
	}
}

func mkpkt(ack bool, seqnum bool, data []byte) []byte {
	var pkt bytes.Buffer
	header := mkhdr(ack, seqnum)
	pkt.Write([]byte{header})
	pkt.Write(data)
	return pkt.Bytes()
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
