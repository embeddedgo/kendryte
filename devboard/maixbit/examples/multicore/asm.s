#include "textflag.h"

#define CSRR(CSR,RD) WORD $(0x2073 + RD<<7 + CSR<<20)
#define mhartid 0xF14
#define s0 8

// func hartid() int
TEXT ·hartid(SB),NOSPLIT|NOFRAME,$0
	CSRR  (mhartid, s0)
	MOV   S0, ret+0(FP)
	RET

// func  loop(n int)
TEXT ·loop(SB),NOSPLIT|NOFRAME,$0
	MOV  n+0(FP), S0
	BEQ  ZERO, S0, end
	ADD  $-1, S0
	JMP  -2(PC)
end:
	RET
