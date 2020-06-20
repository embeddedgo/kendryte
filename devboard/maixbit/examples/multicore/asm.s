#include "textflag.h"

#define CSRR(CSR,RD) WORD $(0x2073 + RD<<7 + CSR<<20)
#define mhartid 0xF14
#define s0 8

// func hartid() int
TEXT ·hartid(SB),NOSPLIT|NOFRAME,$0
	CSRR  (mhartid, s0)
	MOV   48(g), A0    // g.m
	MOV   160(A0), A0  // m.p
	MOVW  (A0), S1     // p.id
	SLL   $1, S1
	OR    S1, S0
	MOV   S0, ret+0(FP)
	RET

// func loop(n int)
TEXT ·loop(SB),NOSPLIT|NOFRAME,$0
	MOV  n+0(FP), S0
	BEQ  ZERO, S0, end
	ADD  $-1, S0
	BNE  ZERO, S0, -1(PC)
end:
	RET
