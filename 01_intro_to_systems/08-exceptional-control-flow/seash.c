#include "csapp.h"
#include <stdio.h>
#include <unistd.h>

#define SHELL "\U0001F41A"
#define WAVE "\U0001F30A"

void eval(char *cmdline, short repl) {
  if (repl) {
    fputs(cmdline, stdout);
    exit(0);
  }
  fputs(cmdline, stdout);
}

int main(int argc, char **argv) {
  char cmdline[MAXLINE];
  int opt;

  char cmd_mode = 0;
  while ((opt = getopt(argc, argv, "c")) != -1) {
    switch (opt) {
    case 'c':
      cmd_mode = 1;
      break;
    }
  }

  while (1) {
    printf("%s ", SHELL);

    if (cmd_mode) {
      eval(argv[optind], 1);
    } else {
      Fgets(cmdline, MAXLINE, stdin);

      if (feof(stdin)) {
        printf("%s %s %s Sea you later %s %s %s\n", WAVE, WAVE, WAVE, WAVE,
               WAVE, WAVE);
        exit(0);
      }

      eval(cmdline, 0);
    }
  }
}
