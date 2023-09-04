/*
 * A few functions from csapp.c
 * https://csapp.cs.cmu.edu/3e/ics3/code/src/csapp.c
 */

#include "csapp.h"

void app_error(char *msg) {
  fprintf(stderr, "%s\n", msg);
  exit(0);
}

void unix_error(char *msg) {
  fprintf(stderr, "%s: %s\n", msg, strerror(errno));
  exit(0);
}

char *Fgets(char *ptr, int n, FILE *stream) {
  char *rptr;

  if (((rptr = fgets(ptr, n, stream)) == NULL) && ferror(stream))
    app_error("Fgets error");

  return rptr;
}
