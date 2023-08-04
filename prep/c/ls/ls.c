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

char *get_human_readable(off_t bytes);
void print_file(char *name, struct stat stbuf);
void stat_walk(char *name, bool recurse);
void dirwalk(char *dir);

bool show_hidden = false;
bool long_format = false;
bool human_readable = false;

int main(int argc, char **argv) {
  int opt;
  bool multiple_files;

  while ((opt = getopt(argc, argv, "Ahl")) != -1) {
    switch (opt) {
    case 'l':
      long_format = true;
      break;
    case 'A':
      show_hidden = true;
      break;
    case 'h':
      human_readable = true;
      break;
    default:
      fprintf(stderr, "Usage: %s [-Ahl] [file...]\n", *argv);
      exit(EXIT_FAILURE);
    }
  }

  multiple_files = argc - optind > 1;

  if (argc == 1) {
    stat_walk(".", true);
  } else {
    while (optind < argc) {
      // print paths as headers when more than one are passed
      if (multiple_files) {
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
  struct stat lstbuf;
  struct stat sbuf;

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "stat_walk: can't stat %s\n", name);
    return;
  }

  if (lstat(name, &lstbuf) == -1) {
    fprintf(stderr, "stat_walk: can't lstat %s\n", name);
    return;
  }

  sbuf = S_ISLNK(lstbuf.st_mode) ? lstbuf : stbuf;

  if (S_ISDIR(sbuf.st_mode)) {
    if (recurse) {
      dirwalk(name);
    } else {
      print_file(name, sbuf);
    }
  } else {
    print_file(name, sbuf);
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

char *get_human_readable(off_t bytes) {
  int si = 0;
  char suffixes[] = {'B', 'K', 'M', 'G', 'T', 'P'};
  char *res;
  while (bytes > 1024) {
    bytes /= 1024;
    si++;
  }
  asprintf(&res, "%5.1lld%c", bytes, suffixes[si]);
  return res;
}

void print_file(char *name, struct stat stbuf) {
  bool is_link = S_ISLNK(stbuf.st_mode);
  char *bn = basename(strdup(name));
  char target[PATH_MAX];
  ssize_t bytes_read;

  if (!show_hidden && bn[0] == '.') {
    return;
  }

  if (is_link) {
    bytes_read = readlink(name, target, sizeof(target) - 1);

    if (bytes_read == -1) {
      fprintf(stderr, "print_file: failed to read symlink for %s\n", name);
      return;
    }

    target[bytes_read] = '\0';
  }

  if (long_format) {
    struct passwd *oi = getpwuid(stbuf.st_uid);
    struct group *gi = getgrgid(stbuf.st_gid);
    struct tm *ti = localtime(&stbuf.st_mtimespec.tv_sec);

    // file mode
    printf(S_ISDIR(stbuf.st_mode) ? "d" : "-");

    // owner
    printf((stbuf.st_mode & S_IRUSR) ? "r" : "-");
    printf((stbuf.st_mode & S_IWUSR) ? "w" : "-");
    printf((stbuf.st_mode & S_IXUSR) ? "x" : "-");

    // group
    printf((stbuf.st_mode & S_IRGRP) ? "r" : "-");
    printf((stbuf.st_mode & S_IWGRP) ? "w" : "-");
    printf((stbuf.st_mode & S_IXGRP) ? "x" : "-");

    // others
    printf((stbuf.st_mode & S_IROTH) ? "r" : "-");
    printf((stbuf.st_mode & S_IWOTH) ? "w" : "-");
    printf((stbuf.st_mode & S_IXOTH) ? "x" : "-");

    // hard links
    printf(" %4u", stbuf.st_nlink);

    // owner name
    printf(" %s", oi->pw_name);

    // group name
    printf(" %s", gi->gr_name);

    // size in bytes
    if (human_readable) {
      printf(" %s", get_human_readable(stbuf.st_size));
    } else {
      printf(" %8lld", stbuf.st_size);
    }

    // last updated
    if (ti != NULL) {
      char tibuf[15];
      strftime(tibuf, sizeof(tibuf), "%b %d %H:%M", ti);
      printf(" %s", tibuf);
    }

    if (is_link) {
      printf(" %s ->", bn);
      printf(" %s\n", target);
    } else {
      printf(" %s\n", bn);
    }
  } else {
    printf("%s\t", bn);
  }
}
