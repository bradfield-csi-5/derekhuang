#include <assert.h>
#include <dirent.h>
#include <grp.h>
#include <libgen.h>
#include <limits.h>
#include <pwd.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>

bool show_hidden = false;
bool long_format = false;
bool human_fmt = false;
bool multiple_paths = false;

void print_file(char *name, struct stat stbuf);
void print_mode(mode_t mode);
void stat_walk(char *name, bool recurse);
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
    stat_walk(".", true);
  } else {
    while (optind < argc) {
      if (multiple_paths) {
        struct stat stbuf;
        if (stat(argv[optind], &stbuf) == -1) {
          fprintf(stderr, "stat_walk: can't access %s\n", argv[optind]);
          return EXIT_FAILURE;
        }
        printf("%s:\n", argv[optind]);
      }
      stat_walk(argv[optind], true);
      if (long_format) {
        printf("\n");
      } else {
        printf("\n\n");
      }
      optind++;
    }
  }

  return EXIT_SUCCESS;
}

void stat_walk(char *name, bool recurse) {
  struct stat stbuf;

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "stat_walk: can't access %s\n", name);
    return;
  }

  if (S_ISDIR(stbuf.st_mode)) {
    if (recurse) {
      dirwalk(name);
    } else {
      print_file(name, stbuf);
    }
  } else {
    print_file(name, stbuf);
  }
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
    if (strcmp(dp->d_name, ".") == 0 || strcmp(dp->d_name, "..") == 0) {
      continue;
    }
    if (strlen(dir) + strlen(dp->d_name) + 2 > sizeof(name)) {
      fprintf(stderr, "dirwalk: name %s/%s too long\n", dir, dp->d_name);
    } else {
      snprintf(name, PATH_MAX, "%s/%s", dir, dp->d_name);
      stat_walk(name, false);
    }
  }
}

void print_file(char *name, struct stat stbuf) {
  char *bn = basename(strdup(name));
  if (!show_hidden && bn[0] == '.') {
    return;
  }
  if (long_format) {
    struct passwd *oi = getpwuid(stbuf.st_uid);
    struct group *gi = getgrgid(stbuf.st_gid);
    struct tm *ti = localtime(&stbuf.st_mtimespec.tv_sec);
    print_mode(stbuf.st_mode);
    printf(" %4u", stbuf.st_nlink);
    printf(" %s", oi->pw_name);
    printf(" %s", gi->gr_name);
    printf(" %10lld", stbuf.st_size);
    if (ti != NULL) {
      char tibuf[15];
      strftime(tibuf, sizeof(tibuf), "%b %d %H:%M", ti);
      printf(" %s", tibuf);
    }
    printf(" %s\n", bn);
  } else {
    printf("%s\t", bn);
  }
}

void print_mode(mode_t mode) {
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
