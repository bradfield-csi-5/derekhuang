package main

import (
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
)

var qtypes = map[string]int{
	"A":     1,
	"NS":    2,
	"CNAME": 5,
	"SOA":   6,
	"MX":    15,
	"TXT":   16,
}

var typeNames = map[uint16]string{
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	6:  "SOA",
	15: "MX",
	16: "TXT",
}

type ResourceRecord struct {
	name     string
	rtype    uint16
	class    uint16
	ttl      uint32
	rdlength uint16
	rdata    string
}

type DNSMessage struct {
	id, flags, qdcount, ancount, nscount, arcount uint16
	questions, answers, authority, additional     []ResourceRecord
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run dns_client.go [domain] [type]")
	}

	query := newQuery(os.Args[1], os.Args[2])

	gPubDNS := syscall.SockaddrInet4{Addr: [4]byte{8, 8, 8, 8}, Port: 53}

	// open socket
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatal("Error opening socket: ", err)
	}

	// bind to any available port
	err = syscall.Bind(sfd, &syscall.SockaddrInet4{Addr: [4]byte{0, 0, 0, 0}, Port: 0})
	if err != nil {
		log.Fatal("Error binding: ", err)
	}

	err = syscall.Sendto(sfd, serialize(query), 0, &gPubDNS)
	if err != nil {
		log.Fatal("Error sendto: ", err)
	}

	// TODO: is 1024 large enough? solution uses 4096
	buf := make([]byte, 1024)
	for {
		_, sockaddr, err := syscall.Recvfrom(sfd, buf, 0)
		if err != nil {
			log.Fatal("Error recvfrom: ", err)
		}

		// expect ipv4
		fromip4, ok := sockaddr.(*syscall.SockaddrInet4)
		if !ok {
			log.Printf("Not ok expecting ipv4\n")
			continue
		}

		// ignore responses from other hosts
		if fromip4.Addr != gPubDNS.Addr || fromip4.Port != gPubDNS.Port {
			log.Printf(
				"Encountered response from an irrelevant host - Addr: %v Port: %d \n",
				fromip4.Addr,
				fromip4.Port,
			)
			continue
		}

		// TODO: parse buf
		log.Printf("Would parse buf...\n")
	}
}

func newQuery(name, qtype string) *DNSMessage {
	return &DNSMessage{
		id:        uint16(rand.Intn(0xffff)),
		flags:     0x100,
		qdcount:   1,
		questions: []ResourceRecord{{name: name, rtype: uint16(qtypes[qtype]), class: 1}},
	}
}

func serialize(msg *DNSMessage) []byte {
	data := make([]byte, 12)

	i := 0
	binary.BigEndian.PutUint16(data[i:], msg.id)
	i += 2
	binary.BigEndian.PutUint16(data[i:], msg.flags)
	i += 2
	binary.BigEndian.PutUint16(data[i:], msg.qdcount)
	i += 2
	binary.BigEndian.PutUint16(data[i:], msg.ancount)
	i += 2
	binary.BigEndian.PutUint16(data[i:], msg.nscount)
	i += 2
	binary.BigEndian.PutUint16(data[i:], msg.arcount)
	i += 2

	for _, q := range msg.questions {
		// append the length of each label before the label
		for _, s := range strings.Split(q.name, ".") {
			data = append(data, byte(len(s)))
			data = append(data, []byte(s)...)
			i += len(s) + 1
		}
		// don't forget the terminating null byte
		data = append(data, 0x00)
		i++
		// what are these four empty bytes for?
		data = append(data, []byte{0, 0, 0, 0}...)
		binary.BigEndian.PutUint16(data[i:], q.rtype)
		i += 2
		binary.BigEndian.PutUint16(data[i:], q.class)
		i += 2
	}

	return data
}
