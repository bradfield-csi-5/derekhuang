package main

import (
	"log"
	"math/rand"
	"syscall"
)

type Header struct {
	ID     int
	QR     byte
	Opcode [4]byte
	AA     byte
	TC     byte
	RD     byte
	// RA     byte
	// Z      [3]byte
	// RCODE  [4]byte
}

type Question struct {
	QNAME
}

func main() {
	// Open socket
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		log.Fatal("Error opening socket: ", err)
	}
	defer syscall.Close(sfd)

	// Bind to any available port
	err = syscall.Bind(sfd, &syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}})
	if err != nil {
		log.Fatal("Error binding: ", err)
	}

	// Encode the query message
	googlePubDns := syscall.SockaddrInet4{Addr: [4]byte{8, 8, 8, 8}, Port: 53}
	dnsID := rand.Intn(0xffff)
	header := Header{
        ID: dnsID,               // 16 bit id
        QR: 0,                   // 0 query 1 response
        Opcode: [4]byte{0, 0, 0, 0}, // opcode
        AA: 0,
        TC: 0,
        RD: 1, // recursion desired
        // RA: 0,
        // Z: [3]byte{0, 0, 0},
        // RCODE: [4]byte{0, 0, 0, 0},
	}
	err = syscall.Sendto(sfd, msg, 0, &googlePubDns)
	if err != nil {
		log.Fatal("Error sendto: ", err)
	}
}
