/*
Write an alternate version of `squeeze(s1, s2)` that deletes each character in
`s1` that matches any character in the string `s2`
*/

#include <stdio.h>

void squeeze(char s1[], char s2[]) {
  int i;
  int j;
  for (i = 0; s1[i] != '\0'; i++) {
    for (j = 0; s2[j] != '\0'; j++) {
      if (s1[i] == s2[j]) {
        s1[i] = '\0';
      }
    }
  }
}

void print_str(char *s) {
  for (; *s != '\0'; s++) {
    printf("print_str: %c\n", *s);
  }
}

int main(int argc, char **argv) {
  char s1[10] = "hello";
  char s2[10] = "world";
  squeeze(s1, s2);
  print_str(s1);
  return 0;
}
