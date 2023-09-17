package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	FILE_HEADER_LEN   = 24
	PACKET_HEADER_LEN = 16
	IPV4              = 0x8
	IPV6              = 0xdd86
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
	SourceAddr          [4]byte
	DestinationAddr     [4]byte

	// TCP headers
	SourcePort            [2]byte
	DestinationPort       [2]byte
	SequenceNum           [4]byte
	AckNum                [4]byte
	DataOffsetAndReserved byte
	_                     byte
	_                     [2]byte
	_                     [2]byte
	_                     [2]byte
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
		fmt.Printf("Version: %d\n", ip_ver)

		ip_header_len := ch.VersionAndHeaderLen & 0x0f
		fmt.Printf("Header length: %d words (%d bytes)\n", ip_header_len, ip_header_len*4)

		ip_total_len := ctoi(ch.TotalLength[:])
		fmt.Printf("Total length: %d bytes\n", ip_total_len)
		fmt.Printf("Payload length: %d bytes\n", ip_total_len-uint64(ip_header_len)*4)
		fmt.Printf("Protocol: %d\n", ch.Protocol)

		ip_src_addr := ctoi(ch.SourceAddr[:])
		fmt.Printf("Source address: 0x%x\n", ip_src_addr)

		ip_dest_addr := ctoi(ch.DestinationAddr[:])
		fmt.Printf("Destination address: 0x%x\n", ip_dest_addr)
		fmt.Println()

		fmt.Println("========== TCP Headers ==========")
		tcp_src_port := ctoi(ch.SourcePort[:])
		fmt.Printf("Source port: %d\n", tcp_src_port)

		tcp_dest_port := ctoi(ch.DestinationPort[:])
		fmt.Printf("Destination port: %d\n", tcp_dest_port)

		tcp_seq_num := ctoi(ch.SequenceNum[:])
		fmt.Printf("Sequence number: %d\n", tcp_seq_num)

		tcp_data_offset := ch.DataOffsetAndReserved >> 4
		fmt.Printf("Data offset (header length): %d words (%d bytes)\n", tcp_data_offset, tcp_data_offset*4)
		fmt.Println()
		fmt.Println()

		i += int(uint32(PACKET_HEADER_LEN) + ch.PacketLength)
	}
	fmt.Printf("%d packets counted\n", count)
}

func check(e error) {
	if e != nil {
		log.Fatal("caught error")
	}
}

func ctoi(buf []byte) uint64 {
	var ret uint64 = 0
	shift := 0
	size := len(buf)
	for i := size - 1; i >= 0; i-- {
		ret |= uint64(buf[i]) << shift
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
