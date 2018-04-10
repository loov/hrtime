#include "textflag.h"

//+build amd64,!gccgo

// func rdtscpAsm() uint64
TEXT ·rdtscpAsm(SB),NOSPLIT,$0-8
	BYTE $0x0F; BYTE $0x01; BYTE $0xF9 // RDTSCP
	SHLQ $32, DX
	ADDQ DX, AX
	MOVQ AX, ret+0(FP)
	RET

// func ·cpuidAsm(op, op2 uint32) (eax, ebx, ecx, edx uint32)
TEXT ·cpuidAsm(SB),NOSPLIT,$8-16
	MOVL  op1+0(FP), AX
	MOVL  op2+4(FP), CX
	CPUID
	MOVL  AX, eax+8(FP)
	MOVL  BX, ebx+12(FP)
	MOVL  CX, ecx+16(FP)
	MOVL  DX, edx+20(FP)
	RET
