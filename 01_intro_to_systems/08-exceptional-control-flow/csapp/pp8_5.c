#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

unsigned int snooze(unsigned int secs) {
  unsigned int rc = sleep(secs);
  printf("Slept for %d of %d secs.\n", secs - rc, secs);
  return rc;
}

void sigint_handler(int sig) { printf("Caught SIGINT! Ignoring...\n"); }

int main() {
  if (signal(SIGINT, sigint_handler) == SIG_ERR) {
    fprintf(stderr, "Something went wrong with catching SIGINT.");
    exit(0);
  }

  unsigned int remainder = snooze(10);

  if (remainder) {
    printf("Had %d seconds remaining.\n", remainder);
  } else {
    printf("Finished sleeping!\n");
  }
}
