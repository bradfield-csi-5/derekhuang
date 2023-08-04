/*
Write the function strend(s,t), which returns 1 if the string t occurs at the
end of the string s, and zero otherwise.
*/

#include <assert.h>
#include <stdio.h>
#include <string.h>

int strend(char *s, char *t);

int main(int argc, char **argv) {
  assert(strend("foo", "oo") == 1);
  assert(strend("foo", "ar") == 0);
  assert(strend("hello", "world") == 0);
  assert(strend("hello", "o") == 1);
  assert(strend("hello", "lo") == 1);
  assert(strend("hello", "llo") == 1);
  printf("5.4 All Passed!\n");
}

int strend(char *s, char *t) {
  s += (strlen(s) - strlen(t));
  while (*s && *t) {
    if (*s++ != *t++) {
      return 0;
    }
  }
  return 1;
}
