/*
 * The magic number is d4c3b2a1 -- byte-ordering is big-endian
 * Major version: 4
 * Minor version: 2
 * Snapshot length: 1514
 * Link-layer header type: 1 ethernet
 */

#include <assert.h>
#include <stdbool.h>
#include <stdio.h>
#include <string.h>

#define MAXLINE 1024
#define PACKLEN_SIZE 4
#define FULL_PACKLEN_SIZE 4
#define MAC_DEST_SIZE 6
#define MAC_SRC_SIZE 6
#define ETHERTYPE_SIZE 2

void print_byte(unsigned char buf[4], bool big_endian, bool decimal) {
  if (big_endian) {
    for (int i = 0; i < 4; i++) {
      if (buf[i] == 0) {
        if (decimal) {
          printf("%03d ", 0);
        } else {
          printf("%02x ", 0);
        }
      } else {
        if (decimal) {
          printf("%03d ", buf[i]);
        } else {
          printf("%02x ", buf[i]);
        }
      }
    }
  } else {
    for (int i = 3; i >= 0; i--) {
      if (buf[i] == 0) {
        if (decimal) {
          printf("%03d ", 0);
        } else {
          printf("%02x ", 0);
        }
      } else {
        if (decimal) {
          printf("%03d ", buf[i]);
        } else {
          printf("%02x ", buf[i]);
        }
      }
    }
  }
}

unsigned int cbtoi(unsigned char buf[4]) {
  int left_shift = 0;
  int ret = 0;
  for (int i = 0; i < 4; i++) {
    ret |= buf[i] << left_shift;
    left_shift += 8;
  }
  return ret;
}

int main(int argc, char **argv) {
  FILE *fp;

  if ((fp = fopen("net.cap", "r")) == NULL) {
    perror("Error opening file");
    return 1;
  }

  char buf[MAXLINE];

  unsigned char packlen_buf[PACKLEN_SIZE];
  unsigned int packlen;

  unsigned char full_packlen_buf[FULL_PACKLEN_SIZE];
  unsigned int full_packlen;

  char mac_dest[MAC_DEST_SIZE];
  char mac_src[MAC_SRC_SIZE];

  char ethertype[ETHERTYPE_SIZE];

  int nitems = 22;
  int count = 0;

  // Packet length starts at byte 33
  fseek(fp, 32, SEEK_CUR);
  while (fread(buf, sizeof(char), nitems, fp)) {
    // TODO: there has to be a better way...
    memcpy(packlen_buf, buf, PACKLEN_SIZE);
    memcpy(full_packlen_buf, buf + PACKLEN_SIZE, FULL_PACKLEN_SIZE);
    memcpy(mac_dest, buf + PACKLEN_SIZE + FULL_PACKLEN_SIZE, MAC_DEST_SIZE);
    memcpy(mac_src, buf + PACKLEN_SIZE + FULL_PACKLEN_SIZE + MAC_DEST_SIZE,
           MAC_SRC_SIZE);
    memcpy(ethertype,
           buf + PACKLEN_SIZE + FULL_PACKLEN_SIZE + MAC_DEST_SIZE +
               MAC_SRC_SIZE,
           ETHERTYPE_SIZE);

    packlen = cbtoi(packlen_buf);
    full_packlen = cbtoi(full_packlen_buf);

    if (packlen != full_packlen) {
      printf("Partial packet: %d of %d\n", packlen, full_packlen);
    }

    // Jump to the next packet length byte in the next header.
    // Subtract 6 from packlen to offset the additional bytes being loaded into
    // the buffer for the parts of the ethernet header we care about.
    if (fseek(fp, packlen - 6, SEEK_CUR) != 0) {
      perror("Error seeking");
      return 1;
    }

    count++;
  }
  fclose(fp);

  printf("%d packets counted\n", count);
  assert(count == 99);
}
