/*
Write a function `setbits(x, p, n, y)` that returns `x` with the `n` bits that
begin at position `p` set to the rightmost `n` bits of `y`, leaving the other
bits unchanged.
*/

#include <assert.h>
#include <stdio.h>

unsigned setbits(unsigned x, unsigned p, unsigned n, unsigned y) {
  unsigned window = p - n + 1;
  unsigned n_mask = (1 << n) - 1;
  unsigned x_mask = ~(n_mask << window);
  unsigned x_cleared = x & x_mask;
  unsigned y_extracted = (y & n_mask) << window;
  return x_cleared | y_extracted;
}

int main(int argc, char **argv) {
  unsigned x = 0b11110101;
  unsigned y = 0b01101010;
  assert(setbits(x, 7, 4, y) == 0b10100101);
  assert(setbits(x, 5, 2, y) == 0b11100101);
  assert(setbits(x, 3, 4, y) == 0b11111010);
  assert(setbits(x, 4, 2, y) == x);
  assert(setbits(0b11110000, 5, 4, 0b00001111) == 0b11111100);
  printf("All Passed!\n");
  return 0;
}
