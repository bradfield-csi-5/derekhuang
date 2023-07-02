#include <assert.h>
#include <dirent.h>
#include <limits.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

bool show_hidden = false;
bool long_format = false;
bool human_fmt = false;
bool multiple_paths = false;

void print_file(char *name, bool recurse);
void print_permissions(mode_t mode);
void dirwalk(char *dir);

int main(int argc, char **argv) {
  int opt;

  while ((opt = getopt(argc, argv, "Ahl")) != -1) {
    switch (opt) {
      case 'l':
        long_format = true;
        break;
      case 'A':
        show_hidden = true;
        break;
      case 'h':
        human_fmt = true;
        break;
      default:
        fprintf(stderr, "Usage: %s [-Ahl] [file...]\n", *argv);
        exit(EXIT_FAILURE);
    }
  }

  multiple_paths = argc - optind > 1;

  if (argc == 1) {
    print_file(".", true);
  } else {
    while (optind < argc) {
      if (multiple_paths) {
        struct stat stbuf;
        if (stat(argv[optind], &stbuf) == -1) {
          fprintf(stderr, "print_file: can't access %s\n", argv[optind]);
          return EXIT_FAILURE;
        }
        printf("%s:\n", argv[optind]);
      }
      print_file(argv[optind], true);
      printf("\n");
      optind++;
    }
  }

  return EXIT_SUCCESS;
}

void print_file(char *name, bool recurse) {
  struct stat stbuf;

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "print_file: can't access %s\n", name);
    return;
  }

  if (S_ISDIR(stbuf.st_mode)) {
    if (recurse) {
      dirwalk(name);
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
  printf(S_ISDIR(mode) ? "d" : "-");

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

void dirwalk(char *dir) {
  struct dirent *dp;
  DIR *dfd;
  char name[PATH_MAX];

  if ((dfd = opendir(dir)) == NULL) {
    fprintf(stderr, "dirwalk: can't open %s\n", dir);
    return;
  }

  while ((dp = readdir(dfd)) != NULL) {
    snprintf(name, PATH_MAX, "%s/%s", dir, dp->d_name);
    print_file(name, false);
  }
}
