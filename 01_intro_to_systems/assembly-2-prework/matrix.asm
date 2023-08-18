section .text
global index
index:
    ; rdi: matrix
    ; rsi: rows
    ; rdx: cols
    ; rcx: rindex
    ; r8: cindex
    imul rcx, rdx
    lea rax, [rcx * 4]
    lea rax, [rax + r8 * 4]
    mov eax, [rdi + rax]
    ret
