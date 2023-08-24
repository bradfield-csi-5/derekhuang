#include "vec.h"
#include <stdint.h>
#include <stdio.h>
#include <time.h>

const long TEST_LOOPS = 10000;

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
