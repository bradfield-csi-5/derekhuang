#include <stdio.h>

#define MAXLINE 1024

int main(int argc, char **argv) {
  char buf[MAXLINE];
  FILE *fp;

  if ((fp = fopen("net.cap", "r")) == NULL) {
    perror("Error opening file");
    return 1;
  }

  while (fread(buf, sizeof(char), MAXLINE, fp)) {
    printf("%s\n", buf);
  }

  fclose(fp);
}
