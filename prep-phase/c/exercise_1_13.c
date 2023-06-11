#include <math.h>
#include <stdio.h>
#include <string.h>

#define MAX_WORDS 10

int main(int argc, char **argv) {
  int c, word_len, longest = 0, printed = 0;
  int histogram[MAX_WORDS] = {0};
  char *words[MAX_WORDS];

  for (int i = 1; i < argc; i++) {
    word_len = (int)strlen(argv[i]);
    longest = fmax(longest, word_len);
    histogram[i - 1] = word_len;
    words[i - 1] = argv[i];
  }

  printf("Histogram\n");
  printf("-------------\n");

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
  printf("\n\n");

  printf("Words\n");
  printf("-------------\n");
  for (int i = 0; i < MAX_WORDS; i++) {
    if (histogram[i] > 0) {
      printf("%d: %s\n", i, words[i]);
    }
  }
}
