/*
In a two’s complement number system, `x &= (x − 1)` deletes the rightmost 1-bit
in `x`. Explain why. Use this observation to write a faster version of bitcount.
*/

#include <assert.h>
#include <stdio.h>

/*
Subtracting 1 always clears the rightmost 1-bit

1111 - 1 == 1110
1110 - 1 == 1101
1101 - 1 == 1100
11110000 - 1 == 11101111

so bitwise AND with the result clears all bits starting from the rightmost 1-bit

1111 & 1110 == 1110
1110 & 1101 == 1100
1101 & 1100 == 1100
11110000 & 11101111 == 11100000
*/

int bitcount(unsigned x) {
  int b;
  for (b = 0; x != 0; x >>= 1) {
    if (x & 01) {
      b++;
    }
  }
  return b;
}

int better_bitcount(unsigned x) {
  int b = 0;
  for (; x != 0; x &= (x - 1)) {
    b++;
  }
  return b;
}

int main(int argc, char **argv) {
  assert(better_bitcount(0b00000000) == 0);
  assert(better_bitcount(0b00000001) == 1);
  assert(better_bitcount(0b00000011) == 2);
  assert(better_bitcount(0b00000111) == 3);
  assert(better_bitcount(0b00001111) == 4);
  assert(better_bitcount(0b00011111) == 5);
  assert(better_bitcount(0b00111111) == 6);
  assert(better_bitcount(0b01111111) == 7);
  assert(better_bitcount(0b11111111) == 8);
  printf("All Passed!\n");
  return 0;
}
