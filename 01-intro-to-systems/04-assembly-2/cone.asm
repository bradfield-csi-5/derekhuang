default rel

section .text
global volume
volume:
    ; v = pi * r^2 * h/3
    ; radius xmm0
    ; height xmm1
    mulss       xmm0, xmm0      ; r^2
    mov         rsi, 3          ; prepare 3 for division
    vcvtsi2ss   xmm2, rsi       ; convert int 3 to single precision floating point
    divss       xmm1, xmm2      ; h/3
    mulss       xmm0, xmm1      ; r^2 * h/3
    mulss       xmm0, [pi]
 	ret

section .data
pi: dd 3.14159265359
