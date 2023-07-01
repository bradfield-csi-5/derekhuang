#include <assert.h>
#include <stdio.h>

/*
Write a function `rightrot(x, n)` that returns the value of the
integer `x` rotated to the right by `n` bit positions.
*/

unsigned rightrot(unsigned x, unsigned n) {
  // only handle 1 byte for simplicity
  n %= 8;
  if (!n) {
    return x;
  }
  unsigned n_mask = (1 << n) - 1;
  unsigned n_extracted_shifted = (x & n_mask) << (8 - n);
  unsigned x_shifted = x >> n;
  unsigned right_extracted = ~n_extracted_shifted & x_shifted;
  return n_extracted_shifted | right_extracted;
}

int main(int argc, char **argv) {
  unsigned x = 0b11110000;
  assert(rightrot(x, 1) == 0b01111000);
  assert(rightrot(x, 2) == 0b00111100);
  assert(rightrot(x, 3) == 0b00011110);
  assert(rightrot(x, 4) == 0b00001111);
  assert(rightrot(x, 5) == 0b10000111);
  assert(rightrot(x, 6) == 0b11000011);
  assert(rightrot(x, 7) == 0b11100001);
  assert(rightrot(x, 8) == x);
  assert(rightrot(x, 11) == 0b00011110);
  printf("All Passed!\n");
  return 0;
}
