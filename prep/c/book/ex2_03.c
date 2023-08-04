#include <ctype.h>
#include <math.h>
#include <stdio.h>
#define MAX_DIGITS 10

double htoi(char s[], int i);

int main() {
  char s[MAX_DIGITS];
  int i = 0;
  int c;
  int res;
  while (((c = getchar()) != EOF) && i < MAX_DIGITS) {
    if (isxdigit(c) || c == 'x' || c == 'X') {
      s[i] = c;
    }
    i++;
  }
  printf("\n");
  res = htoi(s, i);
  printf("htoi: %d\n", res);
  return res;
}

double htoi(char s[], int i) {
  int num = 0;
  int len = i;
  double sum = 0.0;
  for (int j = 0; j < len; j++) {
    switch (s[j]) {
    case '1':
    case '2':
    case '3':
    case '4':
    case '5':
    case '6':
    case '7':
    case '8':
    case '9':
      num = s[j] - '0';
      break;
    case 'A':
    case 'a':
      num = 10;
      break;
    case 'B':
    case 'b':
      num = 11;
      break;
    case 'C':
    case 'c':
      num = 12;
      break;
    case 'D':
    case 'd':
      num = 13;
      break;
    case 'E':
    case 'e':
      num = 14;
      break;
    case 'F':
    case 'f':
      num = 15;
      break;
    default:
      // 0 or x
      num = 0;
      break;
    }
    if (num > 0) {
      // printf("sum before: %f\n", sum);
      // printf("multiplying num %d by pow(16, %d)\n", num, i - 1);
      sum += ceil(num * pow(16, i - 1));
      // printf("sum after: %f\n", sum);
    }
    i--;
  }
  return sum;
}
