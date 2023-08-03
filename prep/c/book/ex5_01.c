#include <ctype.h>
#include <stdio.h>

#define BUFSIZE 100

int getch(void);
void ungetch(int);
int getint(int *);

char buf[BUFSIZE]; /* buffer for ungetch */
int bufp = 0;      /* next free position in buf */

int main(int argc, char **argv) {
  int array[100] = {0};
  for (int n = 0; n < 100 && getint(&array[n]) != EOF; n++) {
    printf("just wrote %d: %d\n", n, array[n]);
  }
  for (int i = 0; i < 100; i++) {
    if (array[i] > 0) {
      printf("%d: %d\n", i, array[i]);
    }
  }
}

/* getint:  get next integer from input into *pn */
int getint(int *pn) {
  int c, sign;

  /* skip white space */
  while (isspace(c = getch())) {
  }
  if (!isdigit(c) && c != EOF && c != '+' && c != '-') {
    ungetch(c); /* it's not a number */
    return 0;
  }
  sign = (c == '-') ? -1 : 1;
  if (c == '+' || c == '-') c = getch();
  while (!isdigit(c = getch())) {
  }
  for (*pn = 0; isdigit(c); c = getch()) *pn = 10 * *pn + (c - '0');
  *pn *= sign;
  if (c != EOF) ungetch(c);
  return c;
}

/* get a (possibly pushed back) character */
int getch(void) { return (bufp > 0) ? buf[--bufp] : getchar(); }

/* push character back on input */
void ungetch(int c) {
  if (bufp >= BUFSIZE)
    printf("ungetch: too many characters\n");
  else
    buf[bufp++] = c;
}
