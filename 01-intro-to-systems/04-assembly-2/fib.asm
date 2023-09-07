section .text
global fib
fib:
	cmp rdi, 1
	jg .rec
	mov eax, edi	; if n <= 1 return
	ret

.rec:
	push rdi		; save n
	dec rdi			; n - 1
	call fib
	pop rdi			; get n for next recursive call
	push rax		; save retval
	sub rdi, 2		; n - 2
	call fib
	pop rdi			; get retval and put in rdi
	add eax, edi	; add the results of recursive calls
	ret
