#include "csapp.h"
#include <stdio.h>

void eval(char *cmdline) { fputs(cmdline, stdout); }

int main(int argc, char **argv) {
  char cmdline[MAXLINE];

  while (1) {
    printf("\U0001F30A ");

    Fgets(cmdline, MAXLINE, stdin);

    if (feof(stdin))
      exit(0);

    eval(cmdline);
  }
}
