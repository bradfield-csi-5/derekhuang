package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

const (
	FILE_HEADER_LEN     = 24
	PACKET_HEADER_LEN   = 16
	ETHERNET_HEADER_LEN = 14
	IPV4_HEADER_LEN     = 20
	TCP_HEADER_LEN      = 20
	IPV4                = 0x0800
	IPV6                = 0x86dd
	IS_BIG_ENDIAN       = 0xd4c3b2a1
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

type NetworkHeaders struct {
	// Ethernet headers
	EthMACDestination [6]byte
	EthMACSource      [6]byte
	EthEthertype      uint16

	// IP headers
	IPVersionAndHeaderLen byte
	_                     byte
	IPTotalLength         uint16
	_                     [2]byte
	_                     [2]byte
	_                     byte
	IPProtocol            byte
	_                     [2]byte
	IPSourceAddr          [4]byte
	IPDestAddr            [4]byte

	// TCP headers
	TCPSourcePort            uint16
	TCPDestPort              uint16
	TCPSeqNum                uint32
	_                        [4]byte
	TCPDataOffsetAndReserved byte
	_                        byte
	_                        [2]byte
	_                        [2]byte
	_                        [2]byte
}

type HTTPDataMap map[uint32][]byte

func main() {
	data, err := os.ReadFile("net.cap")
	check(err)

	fh := FileHeader{}
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.NativeEndian, &fh)
	check(err)

	fmt.Println("========== File Header ==========")
	fmt.Printf("Magic Number: 0x%x\n", fh.MagicNumber)
	fmt.Printf("Major Version: %d\n", fh.MajorVersion)
	fmt.Printf("Minor Version: %d\n", fh.MinorVersion)
	fmt.Printf("Snapshot Length: %d\n", fh.SnapshotLength)
	fmt.Printf("Link Layer Type: %d\n\n", fh.LinkLayerType)

	count := 0
	is_big_endian := fh.MagicNumber == IS_BIG_ENDIAN
	httpResponses := make(HTTPDataMap)
	httpRequests := make(HTTPDataMap)
	for i := FILE_HEADER_LEN; i < len(data); {
		count++

		packet_buf := bytes.NewBuffer(data[i : i+PACKET_HEADER_LEN])
		ph := PacketHeader{}

		// Use the Magic Number to determine byte ordering for the packet header only
		if is_big_endian {
			err = binary.Read(packet_buf, binary.BigEndian, &ph)
		} else {
			err = binary.Read(packet_buf, binary.LittleEndian, &ph)
		}
		check(err)

		fmt.Printf("==================== Packet %d ====================\n", count)

		fmt.Printf("Captured length: %d bytes\n", ph.PacketLength)
		fmt.Printf("Untruncated length: %d bytes\n\n", ph.FullPacketLength)

		next_packet_start := int(uint32(PACKET_HEADER_LEN) + ph.PacketLength)

		// Always use big endian for network headers
		network_buf := bytes.NewBuffer(data[i+PACKET_HEADER_LEN:])
		nh := NetworkHeaders{}
		err = binary.Read(network_buf, binary.BigEndian, &nh)
		check(err)

		fmt.Println("========== Ethernet Headers ==========")

		fmt.Printf("MAC source: ")
		printAddr(nh.EthMACSource[:], "%x", ":")

		fmt.Printf("MAC destination: ")
		printAddr(nh.EthMACDestination[:], "%x", ":")

		ethertype_str := ""
		if nh.EthEthertype == IPV4 {
			ethertype_str = "IPv4"
		} else if nh.EthEthertype == IPV6 {
			ethertype_str = "IPv6"
		} else {
			log.Fatal("Neither IPv4 nor IPv6")
		}

		fmt.Printf("Ethertype: %s\n", ethertype_str)
		fmt.Println()

		fmt.Println("========== IP Headers ==========")

		ip_ver := nh.IPVersionAndHeaderLen >> 4
		fmt.Printf("Version: %d\n", ip_ver)
		if ip_ver != 4 {
			log.Fatalf("IP version == <%d> want <4>\n", ip_ver)
		}

		ip_header_len := nh.IPVersionAndHeaderLen & 0x0f
		fmt.Printf("Header length: %d words (%d bytes)\n", ip_header_len, ip_header_len*4)

		fmt.Printf("Total length: %d bytes\n", nh.IPTotalLength)
		fmt.Printf("Payload length: %d bytes\n", nh.IPTotalLength-uint16(ip_header_len)*4)
		fmt.Printf("Protocol: %d\n", nh.IPProtocol)
		if nh.IPProtocol != 6 {
			log.Fatalf("IP protocol == <%d> want <6>\n", ip_ver)
		}

		fmt.Printf("Source address: ")
		printAddr(nh.IPSourceAddr[:], "%d", ".")

		fmt.Printf("Destination address: ")
		printAddr(nh.IPDestAddr[:], "%d", ".")

		fmt.Println()

		fmt.Println("========== TCP Headers ==========")

		fmt.Printf("Source port: %d\n", nh.TCPSourcePort)
		if nh.TCPSourcePort != 80 && nh.TCPSourcePort != 59295 {
			log.Fatalf("TCP source port == <%d> want <80 || 5925>\n", nh.TCPSourcePort)
		}

		fmt.Printf("Destination port: %d\n", nh.TCPDestPort)
		if nh.TCPDestPort != 80 && nh.TCPDestPort != 59295 {
			log.Fatalf("TCP destination port == <%d> want <80 || 5925>\n", nh.TCPDestPort)
		}

		fmt.Printf("Sequence number: %d\n", nh.TCPSeqNum)

		tcp_data_offset := nh.TCPDataOffsetAndReserved >> 4
		fmt.Printf("Data offset (header length): %d words (%d bytes)\n", tcp_data_offset, tcp_data_offset*4)

		tcp_options_len := tcp_data_offset*4 - TCP_HEADER_LEN
		fmt.Printf("Options length: %d bytes\n", tcp_options_len)

		http_data_start := int(PACKET_HEADER_LEN + ETHERNET_HEADER_LEN + IPV4_HEADER_LEN + tcp_data_offset*4)
		http_data_len := next_packet_start - http_data_start

		fmt.Printf(
			"HTTP data starts at byte %d and ends at byte %d\n",
			http_data_start,
			http_data_start+http_data_len,
		)

		if http_data_len > 0 {
			http_data := make([]byte, http_data_len)
			copy(http_data, data[i+http_data_start:i+http_data_start+http_data_len])

			// Workaround for the two junk bytes at the beginning of the response
			if nh.TCPSourcePort == 80 && nh.TCPDestPort == 59295 && http_data_len > 2 {
				fmt.Printf("HTTP response data length: %d bytes\n", http_data_len)
				if _, exists := httpResponses[nh.TCPSeqNum]; !exists {
					httpResponses[nh.TCPSeqNum] = http_data
				}
			} else if nh.TCPSourcePort == 59295 && nh.TCPDestPort == 80 {
				fmt.Printf("HTTP request data length: %d bytes\n", http_data_len)
				if _, exists := httpRequests[nh.TCPSeqNum]; !exists {
					httpRequests[nh.TCPSeqNum] = http_data
				}
			}
		}

		fmt.Println()
		fmt.Println()

		i += next_packet_start
	}

	fmt.Printf("%d response packets stored\n", len(httpResponses))
	fmt.Printf("%d request packets stored\n", len(httpRequests))
	fmt.Printf("%d packets counted\n", count)

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(getByteData(httpRequests))))
	check(err)
	fmt.Printf("Request: %v\n", req)

	resp, err := http.ReadResponse(
		bufio.NewReader(bytes.NewReader(getByteData(httpResponses))),
		req,
	)
	check(err)
	fmt.Printf("Response: %v\n", resp)
	fmt.Printf("Resp.body: %v\n", resp.Body)

	body, err := io.ReadAll(resp.Body)
	check(err)
	err = os.WriteFile("exercise.jpg", body, 0666)
	check(err)
}

func getByteData(httpData HTTPDataMap) []byte {
	var ret []byte

	// Collect sequence numbers
	seqNums := make([]uint32, 0, len(httpData))
	for sn := range httpData {
		seqNums = append(seqNums, sn)
	}

	// Sort them
	sort.Slice(seqNums, func(i, j int) bool { return seqNums[i] < seqNums[j] })

	// Combine byte data from the map in order
	for _, sn := range seqNums {
		ret = append(ret, httpData[sn]...)
	}

	return ret
}

func check(e error) {
	if e != nil {
		log.Fatalf("check caught error: %v\n", e)
	}
}

func printAddr(addr []byte, numFmt string, sep string) {
	length := len(addr)
	for i := 0; i < length; i++ {
		fmt.Printf(numFmt, addr[i])
		if i < length-1 {
			fmt.Printf("%s", sep)
		}
	}
	fmt.Println()
}
