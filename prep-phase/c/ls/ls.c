#include <assert.h>
#include <limits.h>
#include <stdbool.h>
#include <stdio.h>
#include <sys/stat.h>

#include "dirent.h"

void print_file(char *name, bool recurse);
void print_permissions(mode_t mode);
void traverse(char *dir);

int main(int argc, char **argv) {
  if (argc == 1) {
    print_file(".", true);
  } else {
    while (--argc > 0) {
      print_file(*++argv, true);
    }
  }
  return 0;
}

void print_file(char *name, bool recurse) {
  struct stat stbuf;

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "print_file: can't access %s\n", name);
    return;
  }

  if (S_ISDIR(stbuf.st_mode)) {
    if (recurse) {
      traverse(name);
    } else {
      print_permissions(stbuf.st_mode);
      printf("%10lld %s\n", stbuf.st_size, name);
    }
  } else {
    print_permissions(stbuf.st_mode);
    printf("%10lld %s\n", stbuf.st_size, name);
  }
}

void print_permissions(mode_t mode) {
  // owner
  printf((mode & S_IRUSR) ? "r" : "-");
  printf((mode & S_IWUSR) ? "w" : "-");
  printf((mode & S_IXUSR) ? "x" : "-");

  // group
  printf((mode & S_IRGRP) ? "r" : "-");
  printf((mode & S_IWGRP) ? "w" : "-");
  printf((mode & S_IXGRP) ? "x" : "-");

  // others
  printf((mode & S_IROTH) ? "r" : "-");
  printf((mode & S_IWOTH) ? "w" : "-");
  printf((mode & S_IXOTH) ? "x" : "-");
}

void traverse(char *dir) {
  struct dirent *dp;
  DIR *dfd;
  char name[PATH_MAX];

  if ((dfd = opendir(dir)) == NULL) {
    fprintf(stderr, "traverse: can't open %s\n", dir);
    return;
  }

  while ((dp = readdir(dfd)) != NULL) {
    snprintf(name, PATH_MAX, "%s/%s", dir, dp->d_name);
    print_file(name, false);
  }
}
