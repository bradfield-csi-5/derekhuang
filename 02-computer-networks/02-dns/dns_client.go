package main

import (
	"log"
	"syscall"
)

type Message struct {
	ID      uint16
	QR      byte
	Opcode  [4]byte
	AA      byte
	TC      byte
	RD      byte
	RA      byte
	Z       [3]byte
	RCODE   [4]byte
	QDCOUNT byte
	ANCOUNT byte
	NSCOUNT byte
	ARCOUNT byte
	QNAME   string
	QTYPE   byte
	QCLASS  byte
}

func main() {
	googlePubDns := syscall.SockaddrInet4{Addr: [4]byte{8, 8, 8, 8}, Port: 53}

	// Open socket
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatal("Error opening socket: ", err)
	}
	// defer syscall.Close(sfd)

	// Bind to any available port
	err = syscall.Bind(sfd, &syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: 0})
	if err != nil {
		log.Fatal("Error binding: ", err)
	}

	mesg := []byte{
		// ID
		0x03,
		0x97,
		// flags
		0x01,
		0x00,
		// QDCOUNT
		0x0,
		0x1,
		// ANCOUNT
		0x0,
		0x0,
		// NSCOUNT
		0x0,
		0x0,
		// ARCOUNT
		0x0,
		0x0,
		// QNAME
		0x09, 'w', 'i', 'k', 'i', 'p', 'e', 'd', 'i', 'a', 0x03, 'o', 'r', 'g', 0x0,
		// QTYPE
		0x0,
		0x1,
		// QCLASS
		0x0,
		0x1,
	}
	log.Printf("squery: %x\n", mesg)
	err = syscall.Sendto(sfd, mesg, 0, &googlePubDns)
	if err != nil {
		log.Fatal("Error sendto: ", err)
	}

	buf := make([]byte, 1024)
	for {
		_, sockaddr, err := syscall.Recvfrom(sfd, buf, 0)
		if err != nil {
			log.Fatal("Error recvfrom: ", err)
		}

		// Expect ipv4
		fromip4, ok := sockaddr.(*syscall.SockaddrInet4)
		if !ok {
			log.Printf("Not ok expecting ipv4\n")
			continue
		}

		// Ignore responses from other hosts
		if fromip4.Addr != googlePubDns.Addr || fromip4.Port != googlePubDns.Port {
			log.Printf(
				"Encountered response from an irrelevant host - Addr: %v Port: %d \n",
				fromip4.Addr,
				fromip4.Port,
			)
			continue
		}

		// Parse buf
		log.Printf("Would parse buf here...\n")
	}
}
