#include "pp.h"

int main(int argc, char **argv) {
  unsigned int slp = 10;

  if (argc > 1) {
    slp = atoi(argv[1]);
  }

  if (signal(SIGINT, sigint_handler) == SIG_ERR) {
    fprintf(stderr, "Something went wrong with catching SIGINT.");
    exit(0);
  }

  printf("Sleeping for %d seconds...\n", slp);

  snooze(slp);
}
