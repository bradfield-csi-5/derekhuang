#include "vendor/unity.h"

#include "vec.h"

extern data_t dotproduct(vec_ptr, vec_ptr);
extern data_t dotproduct_mov_len(vec_ptr, vec_ptr);
extern data_t dotproduct_direct_access(vec_ptr, vec_ptr);
extern data_t dotproduct_loop_unroll_2x1(vec_ptr, vec_ptr);
extern data_t dotproduct_loop_unroll_2x2(vec_ptr, vec_ptr);
extern data_t dotproduct_loop_unroll_3x3(vec_ptr, vec_ptr);
extern data_t dotproduct_loop_unroll_10x10(vec_ptr, vec_ptr);

void setUp(void) {}

void tearDown(void) {}

void test_empty(void) {
  vec_ptr u = new_vec(0);
  vec_ptr v = new_vec(0);

  TEST_ASSERT_EQUAL(0, dotproduct(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_mov_len(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_direct_access(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_loop_unroll_2x1(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_loop_unroll_2x2(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_loop_unroll_3x3(u, v));
  TEST_ASSERT_EQUAL(0, dotproduct_loop_unroll_10x10(u, v));

  free_vec(u);
  free_vec(v);
}

void test_basic(void) {
  vec_ptr u = new_vec(3);
  vec_ptr v = new_vec(3);

  set_vec_element(u, 0, 1);
  set_vec_element(u, 1, 2);
  set_vec_element(u, 2, 3);
  set_vec_element(v, 0, 4);
  set_vec_element(v, 1, 5);
  set_vec_element(v, 2, 6);

  TEST_ASSERT_EQUAL(32, dotproduct(u, v));
  TEST_ASSERT_EQUAL(32, dotproduct_mov_len(u, v));
  TEST_ASSERT_EQUAL(32, dotproduct_direct_access(u, v));
  TEST_ASSERT_EQUAL(32, dotproduct_loop_unroll_2x1(u, v));
  TEST_ASSERT_EQUAL(32, dotproduct_loop_unroll_2x2(u, v));
  TEST_ASSERT_EQUAL(32, dotproduct_loop_unroll_10x10(u, v));

  free_vec(u);
  free_vec(v);
}

void test_longer(void) {
  long n = 1000000;
  vec_ptr u = new_vec(n);
  vec_ptr v = new_vec(n);

  for (long i = 0; i < n; i++) {
    set_vec_element(u, i, i + 1);
    set_vec_element(v, i, i + 1);
  }

  long expected = (2 * n * n * n + 3 * n * n + n) / 6;
  TEST_ASSERT_EQUAL(expected, dotproduct(u, v));
  TEST_ASSERT_EQUAL(expected, dotproduct_mov_len(u, v));
  TEST_ASSERT_EQUAL(expected, dotproduct_direct_access(u, v));
  TEST_ASSERT_EQUAL(expected, dotproduct_loop_unroll_2x1(u, v));
  TEST_ASSERT_EQUAL(expected, dotproduct_loop_unroll_2x2(u, v));
  TEST_ASSERT_EQUAL(expected, dotproduct_loop_unroll_10x10(u, v));

  free_vec(u);
  free_vec(v);
}

int main(void) {
  UNITY_BEGIN();

  RUN_TEST(test_empty);
  RUN_TEST(test_basic);
  RUN_TEST(test_longer);

  return UNITY_END();
}
