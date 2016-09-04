//	Copyright © 1994-1999 Lucent Technologies Inc.  All rights reserved.
//	Portions Copyright © 1995-1997 C H Forsyth (forsyth@terzarima.net)
//	Portions Copyright © 1997-1999 Vita Nuova Limited
//	Portions Copyright © 2000-2008 Vita Nuova Holdings Limited (www.vitanuova.com)
//	Portions Copyright © 2004,2006 Bruce Ellis
//	Portions Copyright © 2005-2007 C H Forsyth (forsyth@terzarima.net)
//	Revisions Copyright © 2000-2008 Lucent Technologies Inc. and others
//	Portions Copyright © 2009 The Go Authors.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package riscv

import "cmd/internal/obj"

//go:generate go run ../stringer.go -i $GOFILE -o anames.go -p riscv

const (
	// Base register numberings.
	REG_X0 = obj.RBaseRISCV + iota
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15
	REG_X16
	REG_X17
	REG_X18
	REG_X19
	REG_X20
	REG_X21
	REG_X22
	REG_X23
	REG_X24
	REG_X25
	REG_X26
	REG_X27
	REG_X28
	REG_X29
	REG_X30
	REG_X31

	// FP register numberings.
	REG_F0
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_F8
	REG_F9
	REG_F10
	REG_F11
	REG_F12
	REG_F13
	REG_F14
	REG_F15
	REG_F16
	REG_F17
	REG_F18
	REG_F19
	REG_F20
	REG_F21
	REG_F22
	REG_F23
	REG_F24
	REG_F25
	REG_F26
	REG_F27
	REG_F28
	REG_F29
	REG_F30
	REG_F31

	// Special/control registers.
	// TODO(myenik) Read more and add the ones we need...

	// This marks the end of the register numbering.
	REG_END

	// General registers reassigned to ABI names.
	REG_ZERO = REG_X0
	REG_RA   = REG_X1
	REG_SP   = REG_X2
	REG_GP   = REG_X3 // aka REG_SB
	REG_TP   = REG_X4 // aka REG_G
	REG_T0   = REG_X5
	REG_T1   = REG_X6
	REG_T2   = REG_X7
	REG_S0   = REG_X8
	REG_S1   = REG_X9
	REG_A0   = REG_X10
	REG_A1   = REG_X11
	REG_A2   = REG_X12
	REG_A3   = REG_X13
	REG_A4   = REG_X14
	REG_A5   = REG_X15
	REG_A6   = REG_X16
	REG_A7   = REG_X17
	REG_S2   = REG_X18
	REG_S3   = REG_X19
	REG_S4   = REG_X20
	REG_S5   = REG_X21
	REG_S6   = REG_X22
	REG_S7   = REG_X23
	REG_S8   = REG_X24
	REG_S9   = REG_X25
	REG_S10  = REG_X26
	REG_S11  = REG_X27
	REG_T3   = REG_X28
	REG_T4   = REG_X29
	REG_T5   = REG_X30
	REG_T6   = REG_X31

	// Go runtime register names.
	REG_G    = REG_X4 // G pointer.
	REG_RT1  = REG_S2 // Reserved for runtime (duffzero and duffcopy).
	REG_RT2  = REG_S3 // Reserved for runtime (duffcopy).
	REG_CTXT = REG_S4 // Context for closures.

	// ABI names for floating point registers.
	REG_FT0  = REG_F0
	REG_FT1  = REG_F1
	REG_FT2  = REG_F2
	REG_FT3  = REG_F3
	REG_FT4  = REG_F4
	REG_FT5  = REG_F5
	REG_FT6  = REG_F6
	REG_FT7  = REG_F7
	REG_FS0  = REG_F8
	REG_FS1  = REG_F9
	REG_FA0  = REG_F10
	REG_FA1  = REG_F11
	REG_FA2  = REG_F12
	REG_FA3  = REG_F13
	REG_FA4  = REG_F14
	REG_FA5  = REG_F15
	REG_FA6  = REG_F16
	REG_FA7  = REG_F17
	REG_FS2  = REG_F18
	REG_FS3  = REG_F19
	REG_FS4  = REG_F20
	REG_FS5  = REG_F21
	REG_FS6  = REG_F22
	REG_FS7  = REG_F23
	REG_FS8  = REG_F24
	REG_FS9  = REG_F25
	REG_FS10 = REG_F26
	REG_FS11 = REG_F27
	REG_FT8  = REG_F28
	REG_FT9  = REG_F29
	REG_FT10 = REG_F30
	REG_FT11 = REG_F31
)

// Prog.Mark flags.
const (
	// NEED_PCREL_RELOC is set on AUIPC instructions to indicate that it
	// is a AUIPC+arithmetic pair that needs a R_PCRELRISCV relocation.
	NEED_PCREL_RELOC = 1 << 0
)

// RISC-V mnemonics, as defined in the "opcodes" and "opcodes-pseudo" files of
// riscv-opcodes, as well as some fake mnemonics (e.g., MOV) used only in the
// assembler.
//
// If you modify this table, you MUST run 'go generate' to regenerate anames.go!
const (
	// User ISA

	// 2.4: Integer Computational Instructions
	ASLLIRV32 = obj.ABaseRISCV + obj.A_ARCHSPECIFIC + iota
	ASRLIRV32
	ASRAIRV32

	// 2.5: Control Transfer Instructions
	AJAL
	AJALR
	ABEQ
	ABNE
	ABLT
	ABLTU
	ABGE
	ABGEU

	// 2.7: Memory Model
	AFENCE
	AFENCEI

	// 4.2: Integer Computational Instructions
	AADDI
	ASLTI
	ASLTIU
	AANDI
	AORI
	AXORI
	ASLLI
	ASRLI
	ASRAI
	ALUI
	AAUIPC
	AADD
	ASLT
	ASLTU
	AAND
	AOR
	AXOR
	ASLL
	ASRL
	ASUB
	ASRA
	AADDIW
	ASLLIW
	ASRLIW
	ASRAIW
	AADDW
	ASLLW
	ASRLW
	ASUBW
	ASRAW

	// 4.3: Load and Store Instructions
	ALD
	ALW
	ALWU
	ALH
	ALHU
	ALB
	ALBU
	ASD
	ASW
	ASH
	ASB

	// 4.4: System Instructions
	ARDCYCLE
	ARDCYCLEH
	ARDTIME
	ARDTIMEH
	ARDINSTRET
	ARDINSTRETH

	// 5.1: Multiplication Operations
	AMUL
	AMULH
	AMULHU
	AMULHSU
	AMULW
	ADIV
	ADIVU
	AREM
	AREMU
	ADIVW
	ADIVUW
	AREMW
	AREMUW

	// 6.2: Load-Reserved/Store-Conditional Instructions
	ALRD
	ASCD
	ALRW
	ASCW

	// 6.3: Atomic Memory Operations
	AAMOSWAPD
	AAMOADDD
	AAMOANDD
	AAMOORD
	AAMOXORD
	AAMOMAXD
	AAMOMAXUD
	AAMOMIND
	AAMOMINUD
	AAMOSWAPW
	AAMOADDW
	AAMOANDW
	AAMOORW
	AAMOXORW
	AAMOMAXW
	AAMOMAXUW
	AAMOMINW
	AAMOMINUW

	// 7.2: Floating-Point Control and Status Register
	AFRCSR
	AFSCSR
	AFRRM
	AFSRM
	AFRFLAGS
	AFSFLAGS
	AFSRMI
	AFSFLAGSI

	// 7.5: Single-Precision Load and Store Instructions
	AFLW
	AFSW

	// 7.6: Single-Precision Floating-Point Computational Instructions
	AFADDS
	AFSUBS
	AFMULS
	AFDIVS
	AFMINS
	AFMAXS
	AFSQRTS
	AFMADDS
	AFMSUBS
	AFNMADDS
	AFNMSUBS

	// 7.7: Single-Precision Floating-Point Conversion and Move Instructions
	AFCVTWS
	AFCVTLS
	AFCVTSW
	AFCVTSL
	AFCVTWUS
	AFCVTLUS
	AFCVTSWU
	AFCVTSLU
	AFSGNJS
	AFSGNJNS
	AFSGNJXS
	AFMVSX
	AFMVXS

	// 7.8: Single-Precision Floating-Point Compare Instructions
	AFEQS
	AFLTS
	AFLES

	// 7.9: Single-Precision Floating-Point Classify Instruction
	AFCLASSS

	// 8.2: Double-Precision Load and Store Instructions
	AFLD
	AFSD

	// 8.3: Double-Precision Floating-Point Computational Instructions
	AFADDD
	AFSUBD
	AFMULD
	AFDIVD
	AFMIND
	AFMAXD
	AFSQRTD
	AFMADDD
	AFMSUBD
	AFNMADDD
	AFNMSUBD

	// 8.4: Double-Precision Floating-Point Conversion and Move Instructions
	AFCVTWD
	AFCVTLD
	AFCVTDW
	AFCVTDL
	AFCVTWUD
	AFCVTLUD
	AFCVTDWU
	AFCVTDLU
	AFCVTSD
	AFCVTDS
	AFSGNJD
	AFSGNJND
	AFSGNJXD
	AFMVXD
	AFMVDX

	// 8.5: Double-Precision Floating-Point Compare Instructions
	AFEQD
	AFLTD
	AFLED

	// 8.6: Double-Precision Floating-Point Classify Instruction
	AFCLASSD

	// Privileged ISA

	// 2.1: Instructions to Access CSRs
	ACSRRW
	ACSRRS
	ACSRRC
	ACSRRWI
	ACSRRSI
	ACSRRCI

	// 3.2.1: Instructions to Change Privilege Level
	AECALL
	ASCALL
	AEBREAK
	ASBREAK
	AERET
	ASRET

	// 3.2.2: Trap Redirection Instructions
	AMRTS
	AMRTH
	AHRTS

	// 3.2.3: Wait for Interrupt
	AWFI

	// 4.2.1: Supervisor Memory-Management Fence Instruction
	ASFENCEVM

	// Fake instructions.  These get translated by the assembler into other
	// instructions, based on their operands.
	AFNEGS
	AMOV
	AMOVB
	AMOVBU
	AMOVF
	AMOVH
	AMOVHU
	AMOVW
	AMOVWU
	ASEQZ
	ASNEZ
)

// All unary instructions which write to their arguments (as opposed to reading
// from them) go here.  The assembly parser uses this information to populate
// its AST in a semantically reasonable way.
//
// Any instructions not listed here is assumed to either be non-unary or to read
// from its argument.
var unaryDst = map[obj.As]bool{
	ARDCYCLE:    true,
	ARDCYCLEH:   true,
	ARDTIME:     true,
	ARDTIMEH:    true,
	ARDINSTRET:  true,
	ARDINSTRETH: true,
}

// Operands
const (
	// No operand.  This constant goes in any operand slot which is unused
	// (e.g., the two source register slots in RDCYCLE).
	C_NONE = iota

	// An integer register, either numbered (e.g., R12) or with an ABI name
	// (e.g., T2).
	C_REGI

	// An integer immediate.
	C_IMMI

	// A relative address.
	C_RELADDR

	// A memory address contained in an integer register.
	C_MEM

	// The size of a TEXT section.
	C_TEXTSIZE
)

// Instruction encoding masks
const (
	// ITypeImmMask is a mask including only the immediate portion of
	// I-type instructions.
	ITypeImmMask = 0xfff00000

	// UTypeImmMask is a mask including only the immediate portion of
	// U-type instructions.
	UTypeImmMask = 0xfffff000

	// UJTypeImmMask is a mask including only the immediate portion of
	// UJ-type instructions.
	UJTypeImmMask = UTypeImmMask
)
