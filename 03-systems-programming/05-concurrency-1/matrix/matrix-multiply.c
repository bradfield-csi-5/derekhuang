#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

/*
  A naive implementation of matrix multiplication.

  DO NOT MODIFY THIS FUNCTION, the tests assume it works correctly, which it
  currently does
*/
void matrix_multiply(double **C, double **A, double **B, int a_rows, int a_cols,
                     int b_cols) {
  for (int i = 0; i < a_rows; i++) {
    for (int j = 0; j < b_cols; j++) {
      C[i][j] = 0;
      for (int k = 0; k < a_cols; k++)
        C[i][j] += A[i][k] * B[k][j];
    }
  }
}

typedef struct mult_params {
  double **A;
  double **B;
  double **C;
  int a_rows;
  int a_cols;
  int b_cols;
  int i0;
  int i1;
} mult_params;

const int NUM_THREADS = 4;

void *mult_row(void *vargp) {
  mult_params *args = (mult_params *)vargp;
  for (int i = args->i0; i < args->i1; i++) {
    for (int j = 0; j < args->b_cols; j++) {
      args->C[i][j] = 0;
      for (int k = 0; k < args->a_cols; k++) {
        args->C[i][j] += args->A[i][k] * args->B[k][j];
      }
    }
  }
  return NULL;
}

void parallel_matrix_multiply(double **c, double **a, double **b, int a_rows,
                              int a_cols, int b_cols) {
  mult_params params[NUM_THREADS];
  pthread_t threads[NUM_THREADS];
  for (int t = 0; t < NUM_THREADS; t++) {
    params[t].C = c;
    params[t].A = a;
    params[t].B = b;
    params[t].a_rows = a_rows;
    params[t].a_cols = a_cols;
    params[t].b_cols = b_cols;
    params[t].i0 = a_rows * t / NUM_THREADS;
    params[t].i1 = a_rows * (t + 1) / NUM_THREADS;
    pthread_create(&threads[t], NULL, mult_row, &params[t]);
  }
  for (int t = 0; t < NUM_THREADS; t++) {
    pthread_join(threads[t], NULL);
  }
}
