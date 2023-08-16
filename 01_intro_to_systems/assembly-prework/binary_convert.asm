section .text
global binary_convert
binary_convert:
    xor     rax, rax
    mov     rcx, 2

next:
    cmp     byte [rdi], 0
    je      return
    mov     bl, [rdi]
    sub     bl, 48
    mul     rcx
    add     al, bl
    inc     rdi
    jmp     next

return:
    ret

