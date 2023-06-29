#include <assert.h>
#include <stdio.h>

/*
Write a function `setbits(x, p, n, y)` that returns `x` with the `n` bits that
begin at position `p` set to the rightmost `n` bits of `y`, leaving the other
bits unchanged.
*/

unsigned setbits(unsigned x, unsigned p, unsigned n, unsigned y) {
  unsigned window = p - n + 1;

  unsigned n_mask = (1 << n) - 1;

  unsigned x_mask = ~(n_mask << window);

  unsigned x_cleared = x & x_mask;

  unsigned y_extracted = (y & n_mask) << window;

  return x_cleared | y_extracted;
}

int main(int argc, char **argv) {
  assert(setbits(245, 7, 4, 106) == 165);
  assert(setbits(245, 5, 2, 106) == 229);
  assert(setbits(245, 3, 4, 106) == 250);
  assert(setbits(245, 4, 2, 106) == 245);
  assert(setbits(240, 5, 4, 15) == 252);
  printf("All Passed!\n");
  return 0;
}
