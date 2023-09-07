section .text
global pangram
pangram:
    xor rcx, rcx

.loop:
    movzx edx, byte [rdi]   ; get char
    cmp edx, 0              ; check for null
    je .return              ; short circuit
    or edx, 32              ; normalize to lowercase
    sub edx, 'a'            ; normalize to an index between 0 and 25
    bts ecx, edx            ; set corresponding letter bit
    inc rdi
    jmp .loop

.return:
    xor rax, rax
    and ecx, 0x03ffffff     ; check for missing bits in the first 26
    cmp ecx, 0x03ffffff     ; compare with all 1s
    sete al                 ; set result
    ret

