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
	JE ACEND
	// check tail
	SUBQ $4, SI
	JL ACTAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
ACLOOP:	// Unrolled x2 d[i]|d[i+1] += c
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	ADDPD X0, X1
	ADDPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE ACLOOP
ACTAIL:	// Catch len % 4 == 0
	ADDQ $4, SI
	JE ACEND
ACTL:	// Calc the last values individually d[i] += c
	MOVSD 0(R10), X1
	ADDSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG ACTL
ACEND:
	RET

// func subtrC(c float64, d []float64)
TEXT ·subtrC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE SCEND
	// check tail
	SUBQ $4, SI
	JL SCTAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
SCLOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	SUBPD X0, X1
	SUBPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE SCLOOP
SCTAIL:
	ADDQ $4, SI
	JE SCEND
SCTL:	
	MOVSD 0(R10), X1
	SUBSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG SCTL
SCEND:
	RET

// func multC(c float64, d []float64)
TEXT ·multC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE MCEND
	// check tail
	SUBQ $4, SI
	JL MCTAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
MCLOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	MULPD X0, X1
	MULPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE MCLOOP
MCTAIL:
	ADDQ $4, SI
	JE MCEND
MCTL:	
	MOVSD 0(R10), X1
	MULSD X0,X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG MCTL
MCEND:
	RET

// func divC(c float64, d []float64)
TEXT ·divC(SB), NOSPLIT, $0
	//data ptr
	MOVQ d+8(FP), R10
	// n = data len
	MOVQ d_len+16(FP), SI
	// zero len return
	CMPQ SI, $0
	JE DCEND
	// check tail
	SUBQ $4, SI
	JL DCTAIL
	// load multiplier
	MOVSD c+0(FP), X0
	SHUFPD $0, X0, X0
DCLOOP:	// load d[i] | d[i+1]
	MOVUPD 0(R10), X1
	MOVUPD 16(R10), X2
	DIVPD X0, X1
	DIVPD X0, X2
	MOVUPD X1, 0(R10)
	MOVUPD X2, 16(R10)
	ADDQ $32, R10
	SUBQ $4, SI
	JGE DCLOOP	
DCTAIL:
	ADDQ $4, SI
	JE DCEND
DCTL:	
	MOVSD 0(R10), X1
	DIVSD X0, X1
	MOVSD X1, 0(R10)
	ADDQ $8, R10
	SUBQ $1, SI
	JG DCTL
DCEND:
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
	JE AEND
	// check tail
	SUBQ $2, SI
	JL ATAIL
ALD:
	CMPQ DI, $1
	JE ALT
	SUBQ $2, DI
	JGE ALO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
ALO:
	MOVUPD (R9), X1
	ADDQ $16, R9
	JMP ALOOP
ALT:
	MOVLPD (R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD (R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
ALOOP:	
	MOVUPD (R8), X0
	ADDPD X1, X0
	MOVUPD X0, (R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE ALD
ATAIL:
	ADDQ $2, SI
	JE AEND
ATL:	
	MOVSD (R8), X0
	MOVSD (R9), X1
	ADDSD X1,X0
	MOVSD X0, (R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG ATL
AEND:
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
	JE SEND
	// check tail
	SUBQ $2, SI
	JL STAIL
SLD:
	SUBQ $1, DI
	JE SLT
	SUBQ $1, DI
	JGE SLO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
SLO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP SLOOP
SLT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
SLOOP:	
	MOVUPD 0(R8), X0
	SUBPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE SLD
STAIL:
	ADDQ $2, SI
	JE SEND
STL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	SUBSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG STL
SEND:
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
	JE MEND
	// check tail
	SUBQ $2, SI
	JL MTAIL
MLD:
	SUBQ $1, DI
	JE MLT
	SUBQ $1, DI
	JGE MLO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
MLO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP MLOOP
MLT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
MLOOP:	
	MOVUPD 0(R8), X0
	MULPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE MLD
MTAIL:
	ADDQ $2, SI
	JE MEND
MTL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	MULSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG MTL
MEND:
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
	JE DEND
	// check tail
	SUBQ $2, SI
	JL DTAIL
DLD:
	SUBQ $1, DI
	JE DLT
	SUBQ $1, DI
	JGE DLO
	MOVQ R10, R9
	MOVQ R11, DI
	SUBQ $2, DI
DLO:
	MOVUPD 0(R9), X1
	ADDQ $16, R9
	JMP DLOOP
DLT:
	MOVLPD 0(R9), X1
	MOVQ R10, R9
	MOVQ R11, DI
	MOVHPD 0(R9), X1
	SUBQ $1, DI
	ADDQ $8, R9
DLOOP:	
	MOVUPD 0(R8), X0
	DIVPD X1, X0
	MOVUPD X0, 0(R8)
	ADDQ $16, R8
	SUBQ $2, SI
	JGE DLD
DTAIL:
	ADDQ $2, SI
	JE DEND
DTL:	
	MOVSD 0(R8), X0
	MOVSD 0(R9), X1
	DIVSD X1,X0
	MOVSD X0, 0(R8)
	ADDQ $8, R8
	ADDQ $8, R9
	SUBQ $1, SI
	JG DTL
DEND:
	RET
