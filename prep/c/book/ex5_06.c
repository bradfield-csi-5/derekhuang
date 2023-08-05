/*
Rewrite appropriate programs from earlier chapters and exercises with pointers
instead of array indexing. Good possibilities include getline (Chapters 1 and
4), atoi, itoa, and their variants (Chapters 2, 3, and 4), reverse (Chapter 3),
and strindex and getop (Chapter 4).
*/

#include <assert.h>
#include <ctype.h>
#include <stdio.h>
#include <string.h>

int _getline(char *s, int lim);
int atoi(char *s);
double atof(char *s);
void reverse(char *s);
int strindex(char *s, char *t);

int main(int argc, char **argv) {
  // char s[50];
  // printf("length: %d\n", _getline(s, 50));

  assert(atoi("323") == 323);
  assert(atoi("23") == 23);
  assert(atoi("") == 0);

  char a[] = "abcdefg";
  reverse(&a[0]);
  assert(strcmp(a, "gfedcba") == 0);

  char b[] = "abcdefg";
  assert(strindex(&b[0], "a") == 0);
  assert(strindex(&b[0], "ab") == 0);
  assert(strindex(&b[0], "ac") == -1);
  assert(strindex(&b[0], "def") == 3);
  assert(strindex(&b[0], "dfe") == -1);
  assert(strindex(&b[0], "g") == 6);
  printf("5.6 All Passed!\n");
}

/* getline:  read a line into s, return length */
int _getline(char *s, int lim) {
  int c;
  char *p = s;

  while (--lim > 0 && (c = getchar()) != EOF && c != '\n') {
    *s++ = c;
  }
  if (c == '\n') {
    *s++ = c;
  }
  *s = '\0';
  return (s - 1) - p;
}

/* atoi:  convert string s to integer using atof */
int atoi(char *s) { return (int)atof(s); }

/* atof:  convert string s to double */
double atof(char *s) {
  double val, power;
  int sign;

  while (isspace(*s)) { /* skip white space */
    s++;
  }
  sign = (*s == '-') ? -1 : 1;
  if (*s == '+' || *s == '-')
    s++;
  for (val = 0.0; isdigit(*s); s++)
    val = 10.0 * val + (*s - '0');
  if (*s == '.')
    s++;
  for (power = 1.0; isdigit(*s); s++) {
    val = 10.0 * val + (*s - '0');
    power *= 10.0;
  }
  return sign * val / power;
}

/* reverse:  reverse string s in place */
void reverse(char *s) {
  char c;
  char *p;

  for (p = s + (strlen(s) - 1); s < p; s++, p--) {
    c = *s;
    *s = *p;
    *p = c;
  }
}

/* strindex:  return index of t in s, −1 if none */
int strindex(char *s, char *t) {
  int i = 0;
  char *first;
  char *second;

  while (*s != '\0') {
    if (*s == *t) {
      first = s;
      second = t;
      while (*first++ == *second++) {
        if (*second == '\0') {
          return i;
        }
      }
    }
    i++;
    s++;
  }

  return -1;
}
