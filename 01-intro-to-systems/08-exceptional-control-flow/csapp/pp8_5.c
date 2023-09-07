#include "pp.h"

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
