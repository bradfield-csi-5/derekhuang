/*
Write a pointer version of the function strcat that we showed in
Chapter 2: strcat(s,t) copies the string t to the end of s.
*/

#include <stdio.h>

void _strcat(char s[], char t[]);

int main(int argc, char **argv) {
  char f[] = "foo";
  char b[] = "bar";
  _strcat(f, b);
  printf("f: %s\n", f);
}

/* strcat:  concatenate t to end of s; s must be big enough */
void _strcat(char *s, char *t) {
  while (*s) {
    s++;
  }

  while ((*s++ = *t++)) {
  }
}
