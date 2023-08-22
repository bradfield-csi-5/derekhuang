#include "vec.h"
#include <stdint.h>
#include <stdio.h>
#include <time.h>

#define TEST_LOOPS 10000

void benchmark(long n, data_t (*f)(vec_ptr, vec_ptr), char *func_name) {
  clock_t baseline_start, baseline_end, test_start, test_end;
  double clocks_elapsed, time_elapsed;

  baseline_start = clock();
  int i;
  vec_ptr u = new_vec(n);
  vec_ptr v = new_vec(n);

  for (long i = 0; i < n; i++) {
    set_vec_element(u, i, i + 1);
    set_vec_element(v, i, i + 1);
  }

  for (i = 0; i < TEST_LOOPS; i++) {
  }
  baseline_end = clock();

  test_start = clock();
  for (i = 0; i < TEST_LOOPS; i++) {
    (*f)(u, v);
  }
  test_end = clock();

  clocks_elapsed = test_end - test_start - (baseline_end - baseline_start);
  time_elapsed = clocks_elapsed / CLOCKS_PER_SEC;

  printf("func: %s -> %.2fs to run %d tests (%.2fns per test)\n", func_name,
         time_elapsed, TEST_LOOPS, time_elapsed * 1e9 / TEST_LOOPS);
}

data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;

  for (long i = 0; i < vec_length(u);
       i++) { // we can assume both vectors are same length
    get_vec_element(u, i, &u_val);
    get_vec_element(v, i, &v_val);
    sum += u_val * v_val;
  }
  return sum;
}

data_t dotproduct_mov_len(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;
  long len = vec_length(u);

  for (long i = 0; i < len; i++) {
    get_vec_element(u, i, &u_val);
    get_vec_element(v, i, &v_val);
    sum += u_val * v_val;
  }
  return sum;
}

data_t dotproduct_direct_access(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;
  long len = vec_length(u);
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);

  for (long i = 0; i < len; i++) {
    sum += vec_u[i] * vec_v[i];
  }
  return sum;
}

int main(int argc, char **argv) {
  benchmark(TEST_LOOPS, dotproduct, "original");
  benchmark(TEST_LOOPS, dotproduct_mov_len, "mov_len");
  benchmark(TEST_LOOPS, dotproduct_direct_access, "direct_access");
}
