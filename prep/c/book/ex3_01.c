/*
Our binary search makes two tests inside the loop, when one would suffice (at
the price of more tests outside). Write a version with only one test inside the
loop and measure the difference in run-time.
*/

#include <stdio.h>
#include <time.h>

int binsearch_tt(int x, int v[], int n) {
  int low, high, mid;

  low = 0;
  high = n - 1;
  while (low <= high) {
    mid = (low + high) / 2;
    if (x < v[mid]) {
      high = mid - 1;
    } else if (x > v[mid]) {
      low = mid + 1;
    } else { /* found match */
      return mid;
    }
  }
  return -1; /* no match */
}

int binsearch_ot(int x, int v[], int n) {
  int low, high, mid;
  low = 0;
  high = n - 1;
  while (low <= high) {
    mid = (low + high) / 2;
    if (x < v[mid]) {
      high = mid - 1;
    } else {
      low = mid + 1;
    }
  }
  return v[mid] == x ? mid : -1;
}

int main(int argc, char **argv) {
  int nums[] = {-3, -1, 0, 1, 4, 5};
  clock_t tt_start, tt_end;
  clock_t ot_start, ot_end;
  double tt_cpu_time_used;
  double ot_cpu_time_used;
  for (int i = 0; i < 10; i++) {
    tt_start = clock();
    binsearch_tt(1, nums, 6);
    tt_end = clock();
    tt_cpu_time_used = ((double)(tt_end - tt_start)) / CLOCKS_PER_SEC;
    printf("Run %d - tt_cpu_time_used: %f\n", i + 1, tt_cpu_time_used);
    ot_start = clock();
    binsearch_ot(1, nums, 6);
    ot_end = clock();
    ot_cpu_time_used = ((double)(ot_end - ot_start)) / CLOCKS_PER_SEC;
    printf("Run %d - ot_cpu_time_used: %f\n", i + 1, ot_cpu_time_used);
    printf("\n");
  }
}
