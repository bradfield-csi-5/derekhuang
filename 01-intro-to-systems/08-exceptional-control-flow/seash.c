#include "csapp.h"
#include <signal.h>
#include <stdbool.h>

#define SHELL "\U0001F41A"
#define WAVE "\U0001F30A"

extern char **environ;

void graceful_exit() {
  printf("\n%s %s %s Sea you later %s %s %s\n", WAVE, WAVE, WAVE, WAVE, WAVE,
         WAVE);
  exit(0);
}

void intr_handler(int sig) { graceful_exit(); }

void parseline(char *buf, char **argv, bool repl) {
  int argc = 0;
  char *delim;

  size_t buflen = strlen(buf);

  if (!repl) {
    // The last character will always be \n when reading from stdin;
    // replace it with a space
    buf[buflen - 1] = ' ';
  }

  while ((delim = strchr(buf, ' '))) {
    argv[argc++] = buf;

    // Replace each space with a null char to prevent execve
    // from failing to find the command e.g. '/bin/ls '
    *delim = '\0';

    // Update the buf pointer to the char after the space
    buf = delim + 1;
  }

  if (!repl) {
    argv[argc] = NULL;
  }
}

void eval(char *cmdline, bool repl) {
  char buf[MAXLINE];
  char *argv[MAXARGS];
  strcpy(buf, cmdline);
  parseline(buf, argv, repl);

  pid_t pid;
  if ((pid = Fork()) == 0) {
    // Have the child process run the command
    int execrc;
    if ((execrc = execve(argv[0], argv, environ)) < 0) {
      printf("'%s': command not found. rc: %d\n", argv[0], execrc);
      exit(0);
    }
  }

  int status;
  if (waitpid(pid, &status, 0) < 0) {
    unix_error("waitfg: waitpid error");
  }
}

int main(int argc, char **argv) {
  char cmdline[MAXLINE];
  int opt;

  bool cmd_mode = false;
  while ((opt = getopt(argc, argv, "c")) != -1) {
    switch (opt) {
    case 'c':
      cmd_mode = true;
      break;
    }
  }

  while (1) {
    printf("%s ", SHELL);

    if (cmd_mode) {
      // TODO: fix
      strcpy(cmdline, argv[optind]);
      eval(cmdline, true);
      exit(0);
    } else {
      if (signal(SIGINT, intr_handler) == SIG_ERR) {
        perror("Error setting up SIGINT handler");
        return 1;
      }

      Fgets(cmdline, MAXLINE, stdin);

      if (feof(stdin)) {
        graceful_exit();
      }

      eval(cmdline, false);
    }
  }

  return 0;
}
