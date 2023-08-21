section .text
global binary_convert
binary_convert:
    xor     rax, rax
    mov     rcx, 2      ; store multiplicand

loop:
    mov     bl, [rdi]       ; move first char
    sub     bl, 48          ; convert char into int
    mul     rcx             ; multiply and store in rax
    add     al, bl          ; add result with current digit
    inc     rdi             ; move to next char
    cmp     byte [rdi], 0   ; check if char is null
    jne     loop            ; keep looping if not
    ret
