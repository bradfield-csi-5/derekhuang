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

unsigned int ctoi(unsigned char *buf, int size) {
  int shift = 0;
  unsigned int ret = 0;
  for (int i = 0; i < size; i++) {
    ret |= *(buf++) << shift;
    shift += 8;
  }
  return ret;
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
  unsigned int packlen = 0;
  unsigned int full_packlen = 0;

  // Skip the per-file header
  buf += FILE_HEADER_LEN;

  while (*buf) {
    // Skip the packet timestamp bytes
    buf += 8;

    packlen = ctoi(buf, PACKLEN_SIZE);
    buf += PACKLEN_SIZE;

    full_packlen = ctoi(buf, FULL_PACKLEN_SIZE);
    buf += FULL_PACKLEN_SIZE;

    if (packlen != full_packlen) {
      printf("Partial packet captured: %d of %d\n", packlen, full_packlen);
    }

    buf += packlen;

    count++;
  }

  fclose(fp);

  printf("%d packets counted\n", count);
  assert(count == 99);
}
