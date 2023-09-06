/*
 * A few headers from csapp.h
 * https://csapp.cs.cmu.edu/3e/ics3/code/include/csapp.h
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MAXARGS 16
#define MAXLINE 1024

void unix_error(char *msg);
void app_error(char *msg);

char *Fgets(char *ptr, int n, FILE *stream);
pid_t Fork(void);
