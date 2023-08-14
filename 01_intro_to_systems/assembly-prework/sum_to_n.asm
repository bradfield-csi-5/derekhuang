section .text
global sum_to_n
sum_to_n:
        mov     rax, 0          ; init result to 0
        jmp     .test

.loop:
        add     rax, rdi
        dec     rdi

.test:
        cmp     rdi, 0          ; n > 0?
        jg      .loop           ; keep looping
