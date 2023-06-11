#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MAX_WORDS 10

int main(int argc, char **argv) {
  int c;
  int word_len;
  int arg_start = 1;
  int arg_offset = 1;
  int longest = 0;
  int printed = 0;
  int vertical = 0;
  int histogram[MAX_WORDS] = {0};
  char *words[MAX_WORDS];

  while ((c = getopt(argc, argv, "v")) != -1) {
    switch (c) {
      case 'v':
        vertical = 1;
        break;
      default:
        abort();
    }
  }

  if (vertical == 1) {
    arg_start = 2;
    arg_offset = 2;
  }

  for (; arg_start < argc; arg_start++) {
    word_len = (int)strlen(argv[arg_start]);
    longest = fmax(longest, word_len);
    histogram[arg_start - arg_offset] = word_len;
    words[arg_start - arg_offset] = argv[arg_start];
  }

  printf("Histogram\n");
  printf("-------------\n");

  if (vertical == 1) {
    while (longest > 0) {
      for (int i = 0; i < MAX_WORDS; i++) {
        if (histogram[i] >= longest) {
          printf("# ");
        } else {
          printf("  ");
        }
      }
      printf("\n");
      longest--;
    }

    for (int i = 0; i < MAX_WORDS; i++) {
      if (histogram[i] > 0) {
        printf("%d ", i);
      }
    }
    printf("\n");
  } else {
    for (int i = 0; i < MAX_WORDS; i++) {
      if (histogram[i] > 0) {
        printf("%d: ", i);
        for (int j = 0; j < histogram[i]; j++) {
          printf("# ");
        }
        printf("\n");
      }
    }
  }
  printf("\n");

  printf("Words\n");
  printf("-------------\n");
  for (int i = 0; i < MAX_WORDS; i++) {
    if (histogram[i] > 0) {
      printf("%d: %s\n", i, words[i]);
    }
  }
}
