/*
Rewrite the function `lower`, which converts upper case letters to lower case,
with a conditional expression instead of `if-else`.
*/

#include <stdio.h>

int lower(int c) { return c >= 'A' && c <= 'Z' ? c + 'a' - 'A' : c; }

int main(int argc, char **argv) {
  printf("lower('A') == 'a': %d\n", lower('A') == 'a');
  printf("lower('A') == 'a': %d\n", lower('a') == 'a');
  printf("lower('A') == 'a': %d\n", lower('A') == 'A');
}
