// +build !noasm !appengine

#define NOSPLIT 4

// func addC(c float64, d []float64)
TEXT ·addC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE END
	// check tail
	SUBQ $4, SI
	JL TAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
LOOP:	// Unrolled x2 d[i]|d[i+1] += c
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	ADDPD X0, X1
	ADDPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE LOOP
TAIL:	// Catch len % 4 == 0
	ADDQ $4, SI
	JE END
TL:	// Calc the last values individually d[i] += c
	MOVSD 0(R10), X1
	ADDSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG TL
END:
	RET

// func subtrC(c float64, d []float64)
TEXT ·subtrC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE END
	// check tail
	SUBQ $4, SI
	JL TAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
LOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	SUBPD X0, X1
	SUBPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE LOOP	
TAIL:
	ADDQ $4, SI
	JE END
TL:	
	MOVSD 0(R10), X1
	SUBSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG TL
END:
	RET

// func multC(c float64, d []float64)
TEXT ·multC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE END
	// check tail
	SUBQ $4, SI
	JL TAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
LOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	MULPD X0, X1
	MULPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE LOOP	
TAIL:
	ADDQ $4, SI
	JE END
TL:	
	MOVSD 0(R10), X1
	MULSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG TL
END:
	RET

// func divC(c float64, d []float64)
TEXT ·divC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE END
	// check tail
	SUBQ $4, SI
	JL TAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
LOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	DIVPD X0, X1
	DIVPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE LOOP	
TAIL:
	ADDQ $4, SI
	JE END
TL:	
	MOVSD 0(R10), X1
	DIVSD X0, X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG TL
END:
	RET
	
// func add(a,b []float64)
TEXT ·add(SB), NOSPLIT, $0
	//a data ptr
	MOVQ a_base+0(FP), R8
	//a len
	MOVQ a_len+8(FP), SI
	//b data ptr
	MOVQ b_base+24(FP), R9
	MOVQ R9, R10
	//b len
	MOVQ b_len+32(FP), DI
	MOVQ DI, R11
	// zero len return
	CMPQ SI, $0
	JE END
	// check tail
	SUBQ $2, SI
	JL TAIL
LD:
	CMPQ DI, $1
	JE LT
	SUBQ $2, DI
	JGE LO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
LO:
	MOVUPD (R9), X1
	ADDQ $16, R9
	JMP LOOP
LT:
	MOVLPD (R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD (R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
LOOP:	
	MOVUPD (R8), X0
	ADDPD X1, X0
	MOVUPD X0, (R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE LD
TAIL:
	ADDQ $2, SI
	JE END
TL:	
	MOVSD (R8), X0
	MOVSD (R9), X1
	ADDSD X1,X0
	MOVSD X0, (R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG TL
END:
	RET

// func subtr(a,b []float64)
TEXT ·subtr(SB), NOSPLIT, $0
	//a data ptr
	MOVQ a_base+0(FP), R8
	//a len
	MOVQ a_len+8(FP), SI
	//b data ptr
	MOVQ b_base+24(FP), R9
	MOVQ R9, R10
	//b len
	MOVQ b_len+32(FP), DI
	MOVQ DI, R11
	// zero len return
	MOVQ $0, AX
	CMPQ AX, SI
	JE END
	// check tail
	SUBQ $2, SI
	JL TAIL
LD:
	SUBQ $1, DI
	JE LT
	SUBQ $1, DI
	JGE LO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
LO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP LOOP
LT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
LOOP:	
	MOVUPD 0(R8), X0
	SUBPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE LD
TAIL:
	ADDQ $2, SI
	JE END
TL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	SUBSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG TL
END:
	RET

// func mult(a,b []float64)
TEXT ·mult(SB), NOSPLIT, $0
	//a data ptr
	MOVQ a_base+0(FP), R8
	//a len
	MOVQ a_len+8(FP), SI
	//b data ptr
	MOVQ b_base+24(FP), R9
	MOVQ R9, R10
	//b len
	MOVQ b_len+32(FP), DI
	MOVQ DI, R11
	// zero len return
	MOVQ $0, AX
	CMPQ AX, SI
	JE END
	// check tail
	SUBQ $2, SI
	JL TAIL
LD:
	SUBQ $1, DI
	JE LT
	SUBQ $1, DI
	JGE LO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
LO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP LOOP
LT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
LOOP:	
	MOVUPD 0(R8), X0
	MULPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE LD
TAIL:
	ADDQ $2, SI
	JE END
TL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	MULSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG TL
END:
	RET

// func div(a,b []float64)
TEXT ·div(SB), NOSPLIT, $0
	//a data ptr
	MOVQ a_base+0(FP), R8
	//a len
	MOVQ a_len+8(FP), SI
	//b data ptr
	MOVQ b_base+24(FP), R9
	MOVQ R9, R10
	//b len
	MOVQ b_len+32(FP), DI
	MOVQ DI, R11
	// zero len return
	MOVQ $0, AX
	CMPQ AX, SI
	JE END
	// check tail
	SUBQ $2, SI
	JL TAIL
LD:
	SUBQ $1, DI
	JE LT
	SUBQ $1, DI
	JGE LO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
LO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP LOOP
LT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
LOOP:	
	MOVUPD 0(R8), X0
	DIVPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE LD	
TAIL:
	ADDQ $2, SI
	JE END
TL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	DIVSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG TL
END:
	RET
