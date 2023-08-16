section .text
global binary_convert
binary_convert:
    xor     rax, rax
    mov     cl, 2

next:
    cmp     byte [rdi], 0
    je      return
    mov     bl, [rdi]
    sub     bl, 48
    mul     cl
    add     al, bl
    inc     dil
    jmp     next

return:
    ret

