// Copyright 2020 The Embedded Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

#define AMOXOR(rs, rd, rt) WORD $(0x2600202F + rt<<7 + rt<<15 + rs<<20)
#define AMOOR(rs, rd, rt) WORD $(0x4600202F + rt<<7 + rt<<15 + rs<<20)
#define AMOAND(rs, rd, rt) WORD $(0x6600202F + rt<<7 + rt<<15 + rs<<20)

#define zero 0
#define s0 8
#define a0 10


// func AtomicSetBits(r *mmio.U32, mask uint32)
TEXT ·AtomicSetBits(SB), NOSPLIT, $0-12
	MOV    r+0(FP), A0
	MOVW   mask+8(FP), S0
	AMOOR  (s0, zero, a0)
	RET


// func AtomicClearBits(r *mmio.U32, mask uint32)
TEXT ·AtomicClearBits(SB), NOSPLIT, $0-12
	MOV     r+0(FP), A0
	MOVW    mask+8(FP), S0
	XOR     $-1, S0, S0
	AMOAND  (s0, zero, a0)
	RET

// func AtomicToggleBits(r *mmio.U32, mask uint32)
TEXT ·AtomicToggleBits(SB), NOSPLIT, $0-12
	MOV     r+0(FP), A0
	MOVW    mask+8(FP), S0
	AMOXOR  (s0, zero, a0)
	RET
