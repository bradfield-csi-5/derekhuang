#include <math.h>
#include <stdio.h>
#include <string.h>

#define MAX_WORDS 10

int main(int argc, char **argv) {
  int c, longest = 0, printed = 0;
  int word_len;
  int histogram[MAX_WORDS] = {0};

  for (int i = 1; i < argc; i++) {
    word_len = (int)strlen(argv[i]);
    longest = fmax(longest, word_len);
    histogram[i - 1] = word_len;
  }

  while (longest > 0) {
    for (int i = 0; i < MAX_WORDS; i++) {
      if (histogram[i] >= longest) {
        printf("#");
      } else {
        printf(" ");
      }
    }
    printf("\n");
    longest--;
  }
}
