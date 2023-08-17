section .text
global pangram
pangram:
    xor rcx, rcx

.loop:
    movzx edx, byte [rdi]   ; get char
    cmp edx, 0              ; check for null
    je .return
    cmp edx, 'a'            ; check if uppercase
    jl .check_upper
    cmp edx, 'z'            ; check if letter
    jg .invalid
    sub edx, 'a'            ; lowercase letter; normalize
    bts ecx, edx            ; set corresponding letter bit
    inc rdi
    jmp .loop

.check_upper:
    cmp edx, 'A'            ; check if letter
    jl .invalid
    cmp edx, 'Z'
    jg .invalid
    sub edx, 'A'            ; uppercase letter; normalize
    bts ecx, edx            ; set corresponding letter bit
    inc rdi
    jmp .loop

.invalid:
    inc rdi                 ; skip char and keep looping
    jmp .loop
 
.return:
    xor rax, rax
    and ecx, 0x03ffffff     ; check for missing bits in the first 26
    cmp ecx, 0x03ffffff     ; compare with all 1s
    sete al                 ; set result
    ret

