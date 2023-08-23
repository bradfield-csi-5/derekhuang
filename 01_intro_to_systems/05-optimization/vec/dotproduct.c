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

  printf("%s -> %.2fs to run %d tests (%.2fns per test)\n", func_name,
         time_elapsed, TEST_LOOPS, time_elapsed * 1e9 / TEST_LOOPS);
}

data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;
  long i;

  // we can assume both vectors are the same length
  for (i = 0; i < vec_length(u); i++) {
    get_vec_element(u, i, &u_val);
    get_vec_element(v, i, &v_val);
    sum += u_val * v_val;
  }
  return sum;
}

data_t dotproduct_mov_len(vec_ptr u, vec_ptr v) {
  data_t sum = 0, u_val, v_val;
  long len = vec_length(u);
  long i;

  for (i = 0; i < len; i++) {
    get_vec_element(u, i, &u_val);
    get_vec_element(v, i, &v_val);
    sum += u_val * v_val;
  }
  return sum;
}

data_t dotproduct_direct_access(vec_ptr u, vec_ptr v) {
  data_t sum = 0;
  long len = vec_length(u);
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);
  long i;

  for (i = 0; i < len; i++) {
    sum += vec_u[i] * vec_v[i];
  }
  return sum;
}

data_t dotproduct_loop_unroll_2x1(vec_ptr u, vec_ptr v) {
  data_t sum = 0;
  long len = vec_length(u);
  long limit = len - 1;
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);
  long i;

  for (i = 0; i < limit; i += 2) {
    sum += (vec_u[i] * vec_v[i]) + (vec_u[i + 1] * vec_v[i + 1]);
  }
  for (; i < len; i++) {
    sum += vec_u[i] * vec_v[i];
  }
  return sum;
}

data_t dotproduct_loop_unroll_2x2(vec_ptr u, vec_ptr v) {
  data_t sum0 = 0;
  data_t sum1 = 0;
  long len = vec_length(u);
  long limit = len - 1;
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);
  long i;

  for (i = 0; i < limit; i += 2) {
    sum0 += vec_u[i] * vec_v[i];
    sum1 += vec_u[i + 1] * vec_v[i + 1];
  }
  for (; i < len; i++) {
    sum0 += vec_u[i] * vec_v[i];
  }
  return sum0 + sum1;
}

data_t dotproduct_loop_unroll_3x3(vec_ptr u, vec_ptr v) {
  data_t sum0 = 0;
  data_t sum1 = 0;
  data_t sum2 = 0;
  long len = vec_length(u);
  long limit = len - 2;
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);
  long i;

  for (i = 0; i < limit; i += 3) {
    sum0 += vec_u[i] * vec_v[i];
    sum1 += vec_u[i + 1] * vec_v[i + 1];
    sum2 += vec_u[i + 2] * vec_v[i + 2];
  }
  for (; i < len; i++) {
    sum0 += vec_u[i] * vec_v[i];
  }
  return sum0 + sum1 + sum2;
}

data_t dotproduct_loop_unroll_10x10(vec_ptr u, vec_ptr v) {
  data_t sum0 = 0;
  data_t sum1 = 0;
  data_t sum2 = 0;
  data_t sum3 = 0;
  data_t sum4 = 0;
  data_t sum5 = 0;
  data_t sum6 = 0;
  data_t sum7 = 0;
  data_t sum8 = 0;
  data_t sum9 = 0;
  long len = vec_length(u);
  long limit = len - 9;
  data_t *vec_u = get_vec_start(u);
  data_t *vec_v = get_vec_start(v);
  long i;

  for (i = 0; i < limit; i += 10) {
    sum0 += vec_u[i] * vec_v[i];
    sum1 += vec_u[i + 1] * vec_v[i + 1];
    sum2 += vec_u[i + 2] * vec_v[i + 2];
    sum3 += vec_u[i + 3] * vec_v[i + 3];
    sum4 += vec_u[i + 4] * vec_v[i + 4];
    sum5 += vec_u[i + 5] * vec_v[i + 5];
    sum6 += vec_u[i + 6] * vec_v[i + 6];
    sum7 += vec_u[i + 7] * vec_v[i + 7];
    sum8 += vec_u[i + 8] * vec_v[i + 8];
    sum9 += vec_u[i + 9] * vec_v[i + 9];
  }
  for (; i < len; i++) {
    sum0 += vec_u[i] * vec_v[i];
  }
  return sum0 + sum1 + sum2 + sum3 + sum4 + sum5 + sum6 + sum7 + sum8 + sum9;
}

int main(int argc, char **argv) {
  benchmark(TEST_LOOPS, dotproduct, "original");
  benchmark(TEST_LOOPS, dotproduct_mov_len, "mov_len");
  benchmark(TEST_LOOPS, dotproduct_direct_access, "direct_access");
  benchmark(TEST_LOOPS, dotproduct_loop_unroll_2x1, "loop_unroll_2x1");
  benchmark(TEST_LOOPS, dotproduct_loop_unroll_2x2, "loop_unroll_2x2");
  benchmark(TEST_LOOPS, dotproduct_loop_unroll_3x3, "loop_unroll_3x3");
  benchmark(TEST_LOOPS, dotproduct_loop_unroll_10x10, "loop_unroll_10x10");
}
