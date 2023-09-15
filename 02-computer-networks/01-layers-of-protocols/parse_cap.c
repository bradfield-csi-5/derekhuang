/*
 * Per-file header: 24 bytes
 * Per-packet header: 16 bytes
 * Ethernet header: 14 bytes
 * IPv4 header: 16 bytes
 * TCP header: 16 bytes
 *
 * The magic number is d4c3b2a1 -- byte-ordering is big-endian
 * Major version: 2
 * Minor version: 4
 * Snapshot length: 1514
 * Link-layer header type: 1 ethernet
 * Ethertype: 0x800 (IPv4)
 * MAC addresses:
 *  - 28 60 87 84 e9 c4
 *  - 1b 2e df 60 5e a4
 */

#include <assert.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define FILE_HEADER_LEN 24

#define PACKLEN_SIZE 4
#define FULL_PACKLEN_SIZE 4

#define MAC_DEST_SIZE 6
#define MAC_SRC_SIZE 6
#define ETHERTYPE_SIZE 2
#define IPV4 0x8
#define IPV6 0xdd86

#define IP_TOTAL_LEN_OFFSET 2
#define IP_TOTAL_LEN_SIZE 2
#define IP_TOTAL_LEN_MIN 20
#define IP_TOTAL_LEN_MAX 65535
#define IP_TOTAL_LEN_TO_PROTOCOL_OFFSET 7
#define IP_TCP_PROTOCOL 6
#define IP_PROTOCOL_TO_SRC_OFFSET 3
#define IP_SRC_SIZE 4
#define IP_DEST_SIZE 4

#define TCP_SRC_PORT_SIZE 2
#define TCP_DEST_PORT_SIZE 2
#define TCP_SEQ_NUM_SIZE 4
#define TCP_SEQ_TO_DATA_OFFSET 8

unsigned int ctoi(unsigned char *buf, int size, bool big_endian) {
  int shift = 0;
  unsigned int ret = 0;
  if (big_endian) {
    for (int i = size - 1; i >= 0; i--) {
      ret |= *(buf + i) << shift;
      shift += 8;
    }
  } else {
    for (int i = 0; i < size; i++) {
      ret |= *(buf + i) << shift;
      shift += 8;
    }
  }
  return ret;
}

void print_mac_addr(unsigned char *buf, int size) {
  while (size > 0) {
    printf("%x ", *(buf + --size));
  }
}

int main(int argc, char **argv) {
  FILE *fp;
  unsigned char *buf;
  long filelen;

  if ((fp = fopen("net.cap", "r")) == NULL) {
    perror("Error opening file");
    return 1;
  }

  // Stolen from https://stackoverflow.com/a/22059317/5374314
  // Move file pointer to the end of the file
  if (fseek(fp, 0, SEEK_END) != 0) {
    perror("Error jumping to end of file");
    return 1;
  }

  // Get the length of the file
  filelen = ftell(fp);

  // Move the file pointer back to the beginning of the file
  rewind(fp);

  // Allocate the size of the file to the buffer
  buf = (unsigned char *)malloc(filelen * sizeof(unsigned char));

  // Read the entire file into the buffer
  fread(buf, filelen, 1, fp);

  int count = 0;

  unsigned int packlen;
  unsigned int full_packlen;

  unsigned int ethertype;
  char *ethertype_str;

  unsigned char ip_ver;
  unsigned char ip_header_len;
  unsigned short ip_total_len;
  unsigned int ip_src_addr;
  unsigned int ip_dest_addr;

  unsigned short tcp_src_port;
  unsigned short tcp_dest_port;
  unsigned int tcp_seq_num;
  unsigned char tcp_data_offset;

  // TODO: parse file header and add assertions
  // Skip the per-file header
  buf += FILE_HEADER_LEN;

  while (*buf) {
    count++;

    // Skip the packet timestamp bytes
    buf += 8;

    printf("========== Packet %d ==========\n", count);
    // Packet length which will be altered as different headers are parsed
    // and used to jump to the next packet
    packlen = ctoi(buf, PACKLEN_SIZE, false);
    printf("Captured length: %d bytes\n", packlen);
    buf += PACKLEN_SIZE;

    // Un-truncated packet length
    full_packlen = ctoi(buf, FULL_PACKLEN_SIZE, false);
    printf("Untruncated length: %d bytes\n", full_packlen);
    assert(packlen == full_packlen);

    // This aligns each jump to the start of the next packet
    buf += FULL_PACKLEN_SIZE;
    printf("\n");

    printf("========== Ethernet Headers ==========\n");
    printf("MAC destination: ");
    print_mac_addr(buf, MAC_DEST_SIZE);
    printf("\n");
    buf += MAC_DEST_SIZE;
    packlen -= MAC_DEST_SIZE;

    printf("MAC source: ");
    print_mac_addr(buf, MAC_SRC_SIZE);
    printf("\n");
    buf += MAC_SRC_SIZE;
    packlen -= MAC_SRC_SIZE;

    ethertype = ctoi(buf, ETHERTYPE_SIZE, false);
    if (ethertype == IPV4) {
      ethertype_str = "IPv4";
    } else if (ethertype == IPV6) {
      ethertype_str = "IPv6";
    } else {
      printf("Neither IPv4 nor IPv6: %d 0x%x\n", ethertype, ethertype);
      perror("Error parsing ethertype");
      return 1;
    }
    printf("Ethertype: %s\n", ethertype_str);
    buf += ETHERTYPE_SIZE;
    packlen -= ETHERTYPE_SIZE;
    printf("\n");

    printf("========== IP Headers ==========\n");
    // First 4 bits are version and last 4 bits are header length
    ip_ver = *buf >> 4;
    ip_header_len = *buf & 0x0f;
    printf("Version: %d\n", ip_ver);
    printf("Header length: %d words (%d bytes)\n", ip_header_len,
           ip_header_len * 4);
    assert(ip_ver == 4);
    assert(ip_header_len == 5);
    buf += IP_TOTAL_LEN_OFFSET;
    packlen -= IP_TOTAL_LEN_OFFSET;

    // IP total length and payload length
    ip_total_len = ctoi(buf, IP_TOTAL_LEN_SIZE, true);
    printf("Total length: %d bytes\n", ip_total_len);
    printf("Payload length: %d bytes\n", ip_total_len - ip_header_len * 4);
    assert((IP_TOTAL_LEN_MIN <= ip_total_len) &&
           (ip_total_len <= IP_TOTAL_LEN_MAX));
    buf += IP_TOTAL_LEN_TO_PROTOCOL_OFFSET;
    packlen -= IP_TOTAL_LEN_TO_PROTOCOL_OFFSET;

    // IP protocol
    printf("Protocol: %d\n", *buf);
    assert(*buf == IP_TCP_PROTOCOL);
    buf += IP_PROTOCOL_TO_SRC_OFFSET;
    packlen -= IP_PROTOCOL_TO_SRC_OFFSET;

    // IP source address
    ip_src_addr = ctoi(buf, IP_SRC_SIZE, true);
    printf("Source address: 0x%x\n", ip_src_addr);
    assert(ip_src_addr == 0xc0a80065 || ip_src_addr == 0xc01efc9a);
    buf += IP_SRC_SIZE;
    packlen -= IP_SRC_SIZE;

    // IP destination address
    ip_dest_addr = ctoi(buf, IP_DEST_SIZE, true);
    printf("Destination address: 0x%x\n", ip_dest_addr);
    assert(ip_dest_addr == 0xc0a80065 || ip_dest_addr == 0xc01efc9a);
    buf += IP_DEST_SIZE;
    packlen -= IP_DEST_SIZE;
    printf("\n");

    printf("========== TCP Headers ==========\n");
    tcp_src_port = ctoi(buf, TCP_SRC_PORT_SIZE, true);
    printf("Source port: %d\n", tcp_src_port);
    assert(tcp_src_port == 80 || tcp_src_port == 59295);
    buf += TCP_SRC_PORT_SIZE;
    packlen -= TCP_SRC_PORT_SIZE;

    tcp_dest_port = ctoi(buf, TCP_DEST_PORT_SIZE, true);
    printf("Destination port: %d\n", tcp_dest_port);
    assert(tcp_dest_port == 80 || tcp_dest_port == 59295);
    buf += TCP_DEST_PORT_SIZE;
    packlen -= TCP_DEST_PORT_SIZE;

    tcp_seq_num = ctoi(buf, TCP_SEQ_NUM_SIZE, true);
    printf("Sequence number: %d\n", tcp_seq_num);
    buf += TCP_SEQ_TO_DATA_OFFSET;
    packlen -= TCP_SEQ_TO_DATA_OFFSET;

    tcp_data_offset = *buf >> 4;
    printf("Header length: %d words (%d bytes)\n", tcp_data_offset,
           tcp_data_offset * 4);
    printf("\n");

    // Jump to the next packet
    buf += packlen;
    printf("\n");
  }

  fclose(fp);

  printf("%d packets counted\n", count);
  assert(count == 99);
}
