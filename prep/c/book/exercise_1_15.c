#include <stdio.h>

float fahrenheit_to_celsius(int fahr) { return (5.0 / 9.0) * (fahr - 32.0); }

int main() {
  float fahr, celsius;
  int lower, upper, step;
  lower = 0;
  upper = 300;
  step = 20;

  fahr = lower;
  printf("Fahrenheight  Celsius\n");
  while (fahr <= upper) {
    printf("%12.0f %8.1f\n", fahr, fahrenheit_to_celsius(fahr));
    fahr = fahr + step;
  }
  return 0;
}
