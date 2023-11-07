# Threads
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 256
Naive: 0.082s
Parallel: 0.029s
2.83x speedup between naive and parallel
```
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 512
Naive: 0.675s
Parallel: 0.206s
3.27x speedup between naive and parallel
```
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 1024
Naive: 8.077s
Parallel: 1.671s
4.83x speedup between naive and parallel
```

# Threads + Better Cache Utilization (swap `j` and `k`)
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 256
Naive: 0.082s
Parallel: 0.014s
6.02x speedup between naive and parallel
```
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 512
Naive: 0.672s
Parallel: 0.104s
6.45x speedup between naive and parallel
```
```
❯ cc benchmark.c matrix-multiply.c; ./a.out 1024
Naive: 8.423s
Parallel: 0.847s
9.95x speedup between naive and parallel
```
