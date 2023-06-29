#include <assert.h>
#include <stdio.h>

/*
Write a function `invert(x, p, n)` that returns `x` with the `n` bits that begin
at position `p` inverted (i.e., 1 changed into 0 and vice versa), leaving the
others unchanged.
*/

unsigned invert(unsigned x, unsigned p, unsigned n) {
  unsigned window = p + 1 - n;
  unsigned n_mask = (1 << n) - 1;
  unsigned x_mask = n_mask << window;
  return x ^ x_mask;
}

int main(int argc, char **argv) {
  unsigned x = 0b11110101;
  assert(invert(x, 4, 2) == 0b11101101);
  assert(invert(x, 5, 2) == 0b11000101);
  assert(invert(x, 3, 4) == 0b11111010);
  assert(invert(x, 7, 8) == 0b00001010);
  printf("2-7 All Passed!\n");
  return 0;
}
