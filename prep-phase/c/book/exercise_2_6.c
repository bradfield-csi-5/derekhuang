#include <assert.h>
#include <stdio.h>

/*
Write a function `setbits(x, p, n, y)` that returns `x` with the `n` bits that
begin at position `p` set to the rightmost `n` bits of `y`, leaving the other
bits unchanged.
*/

unsigned setbits(unsigned x, int p, int n, int y) {
  // create mask
  unsigned n_mask = (1 << n) - 1;

  // clear n bits starting at p
  unsigned n_at_p_cleared = ~(n_mask << (p + 1 - n));

  // clear n bits in x
  unsigned n_in_x_cleared = x & n_at_p_cleared;

  // get n right bits in y and shift to p
  unsigned n_in_y_at_p = (y & n_mask) << (p + 1 - n);

  // flip n bits in y in x, preserving other bits
  return n_in_x_cleared | n_in_y_at_p;
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
