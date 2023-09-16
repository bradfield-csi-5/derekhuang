package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	FILE_HEADER_LEN     = 24
	PACKET_HEADER_LEN   = 16
	ETHERNET_HEADER_LEN = 14
	IP_HEADER_LEN       = 20
	TCP_HEADER_LEN      = 20

	PACKLEN_SIZE                    = 4
	FULL_PACKLEN_SIZE               = 4
	MAC_DEST_SIZE                   = 6
	MAC_SRC_SIZE                    = 6
	ETHERTYPE_SIZE                  = 2
	IPV4                            = 0x8
	IPV6                            = 0xdd86
	IP_TOTAL_LEN_OFFSET             = 2
	IP_TOTAL_LEN_SIZE               = 2
	IP_TOTAL_LEN_MIN                = 20
	IP_TOTAL_LEN_MAX                = 65535
	IP_TOTAL_LEN_TO_PROTOCOL_OFFSET = 7
	IP_TCP_PROTOCOL                 = 6
	IP_PROTOCOL_TO_SRC_OFFSET       = 3
	IP_SRC_SIZE                     = 4
	IP_DEST_SIZE                    = 4
	TCP_SRC_PORT_SIZE               = 2
	TCP_DEST_PORT_SIZE              = 2
	TCP_SEQ_NUM_SIZE                = 4
	TCP_SEQ_TO_DATA_OFFSET          = 8
)

type FileHeader struct {
	MagicNumber    uint32
	MajorVersion   uint16
	MinorVersion   uint16
	_              [4]byte
	_              [4]byte
	SnapshotLength uint32
	LinkLayerType  uint32
}

type PacketHeader struct {
	_                [4]byte
	_                [4]byte
	PacketLength     uint32
	FullPacketLength uint32
}

type EthernetHeader struct {
	MACDestination [6]byte
	MACSource      [6]byte
	Ethertype      uint16
}

type IPHeader struct {
	VersionAndHeaderLen byte
	_                   byte
	TotalLength         uint16
	_                   [2]byte
	_                   [2]byte
	_                   byte
	Protocol            byte
	_                   [2]byte
	SourceAddr          uint32
	DestinationAddr     uint32
}

type TCPHeader struct {
	SourcePort            uint16
	DestinationPort       uint16
	SequenceNum           uint32
	AckNum                uint32
	DataOffsetAndReserved byte
	_                     byte
	_                     [2]byte
	_                     [2]byte
	_                     [2]byte
}

type CombinedHeader struct {
	// Packet headers
	_                [4]byte
	_                [4]byte
	PacketLength     uint32
	FullPacketLength uint32

	// Ethernet headers
	MACDestination [6]byte
	MACSource      [6]byte
	Ethertype      uint16

	// IP headers
	VersionAndHeaderLen byte
	_                   byte
	TotalLength         [2]byte
	_                   [2]byte
	_                   [2]byte
	_                   byte
	Protocol            byte
	_                   [2]byte
	SourceAddr          uint32
	DestinationAddr     uint32

	// TCP headers
	SourcePort            uint16
	DestinationPort       uint16
	SequenceNum           uint32
	AckNum                uint32
	DataOffsetAndReserved byte
	_                     byte
	_                     [2]byte
	_                     [2]byte
	_                     [2]byte
}

var headers = []struct {
	header interface{}
	size   uint8
}{
	// {FileHeader{}, FILE_HEADER_LEN},
	{PacketHeader{}, PACKET_HEADER_LEN},
	{EthernetHeader{}, ETHERNET_HEADER_LEN},
	{IPHeader{}, IP_HEADER_LEN},
	{TCPHeader{}, TCP_HEADER_LEN},
}

func main() {
	data, err := os.ReadFile("../net.cap")
	check(err)

	fh := FileHeader{}
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.LittleEndian, &fh)
	check(err)

	fmt.Println("========== File Header ==========")
	fmt.Printf("Magic Number: 0x%x\n", fh.MagicNumber)
	fmt.Printf("Major Version: %d\n", fh.MajorVersion)
	fmt.Printf("Minor Version: %d\n", fh.MinorVersion)
	fmt.Printf("Snapshot Length: %d\n", fh.SnapshotLength)
	fmt.Printf("Link Layer Type: %d\n\n", fh.LinkLayerType)

	count := 0
	length := len(data)
	for i := FILE_HEADER_LEN; i < length; {
		count++

		buf = bytes.NewBuffer(data[i:])

		ch := CombinedHeader{}
		err = binary.Read(buf, binary.LittleEndian, &ch)
		check(err)

		fmt.Printf("========== Packet %d ==========\n", count)
		fmt.Printf("Captured length: %d bytes\n", ch.PacketLength)
		fmt.Printf("Untruncated length: %d bytes\n\n", ch.FullPacketLength)

		fmt.Println("========== Ethernet Headers ==========")
		fmt.Printf("MAC destination: ")
		printMACAddr(ch.MACDestination)
		fmt.Printf("MAC source: ")
		printMACAddr(ch.MACSource)
		ethertype_str := ""
		if ch.Ethertype == IPV4 {
			ethertype_str = "IPv4"
		} else if ch.Ethertype == IPV6 {
			ethertype_str = "IPv6"
		} else {
			log.Fatal("Neither IPv4 nor IPv6")
		}
		fmt.Printf("Ethertype: %s\n", ethertype_str)
		fmt.Println()

		fmt.Println("========== IP Headers ==========")
		ip_ver := ch.VersionAndHeaderLen >> 4
		ip_header_len := ch.VersionAndHeaderLen & 0x0f
		ip_total_len := getTotalLen(ch.TotalLength)
		fmt.Printf("Version: %d\n", ip_ver)
		fmt.Printf("Header length: %d words (%d bytes)\n", ip_header_len, ip_header_len*4)
		fmt.Printf("Total length: %d bytes\n", ip_total_len)
		fmt.Printf("Payload length: %d bytes\n", ip_total_len-uint16(ip_header_len)*4)
		fmt.Println()
		fmt.Println()
		// fmt.Printf("Entire Header: %+v\n", ch)

		i += int(uint32(PACKET_HEADER_LEN) + ch.PacketLength)
	}
	fmt.Printf("%d packets counted\n", count)
}

func check(e error) {
	if e != nil {
		log.Fatal("caught error")
	}
}

func getTotalLen(tl [2]byte) uint16 {
	var shift uint8 = 0
	var ret uint16 = 0
	for i := 1; i >= 0; i-- {
		ret |= uint16(tl[i] << shift)
		shift += 8
	}
	return ret
}

func printMACAddr(addr [6]byte) {
	for i := 0; i < 6; i++ {
		fmt.Printf("%x ", addr[i])
	}
	fmt.Println()
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal("error reading bytes from file")
	}
	return bytes
}
