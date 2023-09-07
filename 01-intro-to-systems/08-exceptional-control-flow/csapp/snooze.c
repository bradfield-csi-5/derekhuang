#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void sigint_handler(int sig) { printf("Caught SIGINT! Ignoring...\n"); }

unsigned int snooze(unsigned int secs) {
  unsigned int rc = sleep(secs);
  printf("Slept for %d of %d secs.\n", secs - rc, secs);
  return rc;
}
