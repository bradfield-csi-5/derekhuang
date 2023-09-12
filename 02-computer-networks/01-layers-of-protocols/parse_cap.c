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

#define MAXLINE 1024

void print_byte(unsigned char *buf, bool big_endian, bool decimal) {
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

int main(int argc, char **argv) {
  int buf[MAXLINE];
  FILE *fp;

  if ((fp = fopen("net.cap", "r")) == NULL) {
    perror("Error opening file");
    return 1;
  }

  int count = 0;
  int nitems = 2;
  int pack_len;
  int full_pack_len;

  // Packet length starts at byte 33
  fseek(fp, 32, SEEK_CUR);
  while (fread(buf, sizeof(int), nitems, fp)) {
    pack_len = buf[0];
    full_pack_len = buf[1];

    if (pack_len != full_pack_len) {
      printf("Partial packet: %d of %d\n", pack_len, full_pack_len);
    }

    // Jump to the next packet length byte in the next header
    if (fseek(fp, pack_len + 8, SEEK_CUR) != 0) {
      perror("Error seeking");
      return 1;
    }

    count++;
  }
  fclose(fp);

  printf("%d packets counted\n", count);
  assert(count == 99);
}
