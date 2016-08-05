// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import "cmd/internal/obj/riscv"

func init() {
	var regNamesRISCV []string
	var gpMask, fpMask, gpspMask, gpspsbMask regMask
	regNamed := make(map[string]regMask)

	// Build the list of register names, creating an appropriately indexed
	// regMask for the gp and fp registers as we go.
	addreg := func(r int) regMask {
		mask := regMask(1) << uint(len(regNamesRISCV))
		name := riscv.RegNames[int16(r)]
		regNamesRISCV = append(regNamesRISCV, name)
		regNamed[name] = mask
		return mask
	}
	for r := riscv.REG_X0; r <= riscv.REG_X31; r++ {
		mask := addreg(r)
		// Add general purpose registers to gpMask.
		switch r {
		// Special registers that we must leave alone.
		// TODO: Is this list right?
		case riscv.REG_ZERO, riscv.REG_RA, riscv.REG_G:
		case riscv.REG_SB:
			gpspsbMask |= mask
		case riscv.REG_SP:
			gpspMask |= mask
			gpspsbMask |= mask
		default:
			gpMask |= mask
			gpspMask |= mask
			gpspsbMask |= mask
		}
	}
	for r := riscv.REG_F0; r <= riscv.REG_F31; r++ {
		mask := addreg(r)
		fpMask |= mask
	}

	if len(regNamesRISCV) > 64 {
		// regMask is only 64 bits.
		panic("Too many RISCV registers")
	}

	regCtxt := regNamed["CTXT"]
	callerSave := gpMask | fpMask | regNamed["g"]

	var (
		gpstore = regInfo{inputs: []regMask{gpspsbMask, gpspMask, 0}} // SB in first input so we can load from a global, but not in second to avoid using SB as a temporary register
		gp01    = regInfo{outputs: []regMask{gpMask}}
		// FIXME(prattmic): This is a hack to get things to build, but it probably
		// not correct.
		gp11   = regInfo{inputs: []regMask{gpMask}, outputs: []regMask{gpMask}}
		gp21   = regInfo{inputs: []regMask{gpMask, gpMask}, outputs: []regMask{gpMask}}
		gp20   = regInfo{inputs: []regMask{gpMask, gpMask}, outputs: []regMask{}}
		gpload = regInfo{inputs: []regMask{gpspsbMask, 0}, outputs: []regMask{gpMask}}
		gp11sb = regInfo{inputs: []regMask{gpspsbMask}, outputs: []regMask{gpMask}}

		fp11 = regInfo{inputs: []regMask{fpMask}, outputs: []regMask{fpMask}}
		fp21 = regInfo{inputs: []regMask{fpMask, fpMask}, outputs: []regMask{fpMask}}
		gpfp = regInfo{inputs: []regMask{gpMask}, outputs: []regMask{fpMask}}
		fpgp = regInfo{inputs: []regMask{fpMask}, outputs: []regMask{gpMask}}

		call        = regInfo{clobbers: callerSave}
		callClosure = regInfo{inputs: []regMask{gpspMask, regCtxt, 0}, clobbers: callerSave}
		callInter   = regInfo{inputs: []regMask{gpMask}, clobbers: callerSave}
	)

	RISCVops := []opData{
		{name: "ADD", argLength: 2, reg: gp21, asm: "ADD", commutative: true}, // arg0 + arg1
		{name: "ADDI", argLength: 1, reg: gp11sb, asm: "ADDI", aux: "Int64"},  // arg0 + auxint
		{name: "SUB", argLength: 2, reg: gp21, asm: "SUB"},                    // arg0 - arg1

		// M extension. H means high (i.e., it returns the top bits of
		// the result). U means unsigned. W means word (i.e., 32-bit).
		{name: "MUL", argLength: 2, reg: gp21, asm: "MUL", commutative: true, typ: "Int64"}, // arg0 * arg1
		{name: "MULW", argLength: 2, reg: gp21, asm: "MULW", commutative: true, typ: "Int32"},
		{name: "MULH", argLength: 2, reg: gp21, asm: "MULH", commutative: true, typ: "Int64"},
		{name: "MULHU", argLength: 2, reg: gp21, asm: "MULHU", commutative: true, typ: "UInt64"},
		{name: "DIV", argLength: 2, reg: gp21, asm: "DIV", typ: "Int64"}, // arg0 / arg1
		{name: "DIVU", argLength: 2, reg: gp21, asm: "DIVU", typ: "UInt64"},
		{name: "DIVW", argLength: 2, reg: gp21, asm: "DIVW", typ: "Int32"},
		{name: "DIVUW", argLength: 2, reg: gp21, asm: "DIVUW", typ: "UInt32"},
		{name: "REM", argLength: 2, reg: gp21, asm: "REM", typ: "Int64"}, // arg0 % arg1
		{name: "REMU", argLength: 2, reg: gp21, asm: "REMU", typ: "UInt64"},
		{name: "REMW", argLength: 2, reg: gp21, asm: "REMW", typ: "Int32"},
		{name: "REMUW", argLength: 2, reg: gp21, asm: "REMUW", typ: "UInt32"},

		{name: "MOVmem", argLength: 1, reg: gp11sb, asm: "MOV", aux: "SymOff"}, // arg0 + auxint + offset encoded in aux
		// auxint+aux == add auxint and the offset of the symbol in aux (if any) to the effective address

		{name: "MOVBconst", reg: gp01, asm: "MOV", typ: "UInt8", aux: "Int8", rematerializeable: true},   // 8 low bits of auxint
		{name: "MOVWconst", reg: gp01, asm: "MOV", typ: "UInt16", aux: "Int16", rematerializeable: true}, // 16 low bits of auxint
		{name: "MOVLconst", reg: gp01, asm: "MOV", typ: "UInt32", aux: "Int32", rematerializeable: true}, // 32 low bits of auxint
		{name: "MOVQconst", reg: gp01, asm: "MOV", typ: "UInt64", aux: "Int64", rematerializeable: true}, // auxint

		// Loads: load <size> bits from arg0+auxint+aux and extend to 64 bits; arg1=mem
		{name: "LB", argLength: 2, reg: gpload, asm: "MOVB", aux: "SymOff", typ: "Int8"},     //  8 bits, sign extend
		{name: "LH", argLength: 2, reg: gpload, asm: "MOVH", aux: "SymOff", typ: "Int16"},    // 16 bits, sign extend
		{name: "LW", argLength: 2, reg: gpload, asm: "MOVW", aux: "SymOff", typ: "Int32"},    // 32 bits, sign extend
		{name: "LD", argLength: 2, reg: gpload, asm: "MOV", aux: "SymOff", typ: "Int64"},     // 64 bits
		{name: "LBU", argLength: 2, reg: gpload, asm: "MOVBU", aux: "SymOff", typ: "UInt8"},  //  8 bits, zero extend
		{name: "LHU", argLength: 2, reg: gpload, asm: "MOVHU", aux: "SymOff", typ: "UInt16"}, // 16 bits, zero extend
		{name: "LWU", argLength: 2, reg: gpload, asm: "MOVWU", aux: "SymOff", typ: "UInt32"}, // 32 bits, zero extend

		// Stores: store <size> lowest bits in arg1 to arg0+auxint+aux; arg2=mem
		// TODO: rename SB_ to SB when https://go-review.googlesource.com/24649 goes in.
		{name: "SB_", argLength: 3, reg: gpstore, asm: "MOVB", aux: "SymOff", typ: "Mem"}, //  8 bits
		{name: "SH", argLength: 3, reg: gpstore, asm: "MOVH", aux: "SymOff", typ: "Mem"},  // 16 bits
		{name: "SW", argLength: 3, reg: gpstore, asm: "MOVW", aux: "SymOff", typ: "Mem"},  // 32 bits
		{name: "SD", argLength: 3, reg: gpstore, asm: "MOV", aux: "SymOff", typ: "Mem"},   // 64 bits

		// Shift ops
		{name: "SLL", argLength: 2, reg: gp21, asm: "SLL"},                 // arg0 << aux1
		{name: "SRA", argLength: 2, reg: gp21, asm: "SRA"},                 // arg0 >> aux1, signed
		{name: "SRL", argLength: 2, reg: gp21, asm: "SRL"},                 // arg0 >> aux1, unsigned
		{name: "SLLI", argLength: 1, reg: gp11, asm: "SLLI", aux: "Int64"}, // arg0 << auxint
		{name: "SRAI", argLength: 1, reg: gp11, asm: "SRAI", aux: "Int64"}, // arg0 >> auxint, signed
		{name: "SRLI", argLength: 1, reg: gp11, asm: "SRLI", aux: "Int64"}, // arg0 >> auxint, unsigned

		// Bitwise ops
		{name: "XOR", argLength: 2, reg: gp21, asm: "XOR", commutative: true}, // arg0 ^ arg1
		{name: "XORI", argLength: 1, reg: gp11, asm: "XORI", aux: "Int64"},    // arg0 ^ auxint
		{name: "OR", argLength: 2, reg: gp21, asm: "OR", commutative: true},   // arg0 | arg1
		{name: "ORI", argLength: 1, reg: gp11, asm: "ORI", aux: "Int64"},      // arg0 | auxint
		{name: "AND", argLength: 2, reg: gp21, asm: "AND", commutative: true}, // arg0 & arg1
		{name: "ANDI", argLength: 1, reg: gp11, asm: "ANDI", aux: "Int64"},    // arg0 & auxint

		// Generate boolean values
		{name: "SEQZ", argLength: 1, reg: gp11, asm: "SEQZ"},                 // arg0 == 0, result is 0 or 1
		{name: "SNEZ", argLength: 1, reg: gp11, asm: "SNEZ"},                 // arg0 != 0, result is 0 or 1
		{name: "SLT", argLength: 2, reg: gp21, asm: "SLT"},                   // arg0 < arg1, result is 0 or 1
		{name: "SLTI", argLength: 1, reg: gp11, asm: "SLTI", aux: "Int64"},   // arg0 < auxint, result is 0 or 1
		{name: "SLTU", argLength: 2, reg: gp21, asm: "SLTU"},                 // arg0 < arg1, unsigned, result is 0 or 1
		{name: "SLTIU", argLength: 1, reg: gp11, asm: "SLTIU", aux: "Int64"}, // arg0 < auxint, unsigned, result is 0 or 1

		// Flag pseudo-ops.
		// RISC-V doesn't have flags, but SSA wants branches to be flag-based.
		// These are pseudo-ops that contain the arguments to the comparison.
		// There is a single branching block type, BRANCH,
		// which accepts one of these values as its control.
		// During code generation we consult the control value
		// to emit the correct conditional branch instruction.
		{name: "BEQ", argLength: 2, reg: gp20, asm: "BEQ", typ: "Flags"},   // branch if arg0 == arg1
		{name: "BNE", argLength: 2, reg: gp20, asm: "BNE", typ: "Flags"},   // branch if arg0 != arg1
		{name: "BLT", argLength: 2, reg: gp20, asm: "BLT", typ: "Flags"},   // branch if arg0 < arg1
		{name: "BLTU", argLength: 2, reg: gp20, asm: "BLTU", typ: "Flags"}, // branch if arg0 < arg1, unsigned
		{name: "BGE", argLength: 2, reg: gp20, asm: "BGE", typ: "Flags"},   // branch if arg0 >= arg1
		{name: "BGEU", argLength: 2, reg: gp20, asm: "BGEU", typ: "Flags"}, // branch if arg0 >= arg1, unsigned

		// MOVconvert converts between pointers and integers.
		// We have a special op for this so as to not confuse GC
		// (particularly stack maps). It takes a memory arg so it
		// gets correctly ordered with respect to GC safepoints.
		{name: "MOVconvert", argLength: 2, reg: gp11, asm: "MOV"}, // arg0, but converted to int/ptr as appropriate; arg1=mem

		// Calls
		{name: "CALLstatic", argLength: 1, reg: call, aux: "SymOff"},        // call static function aux.(*gc.Sym). arg0=mem, auxint=argsize, returns mem
		{name: "CALLclosure", argLength: 3, reg: callClosure, aux: "Int64"}, // call function via closure. arg0=codeptr, arg1=closure, arg2=mem, auxint=argsize, returns mem
		{name: "CALLdefer", argLength: 1, reg: call, aux: "Int64"},          // call deferproc. arg0=mem, auxint=argsize, returns mem
		{name: "CALLgo", argLength: 1, reg: call, aux: "Int64"},             // call newproc. arg0=mem, auxint=argsize, returns mem
		{name: "CALLinter", argLength: 2, reg: callInter, aux: "Int64"},     // call fn by pointer. arg0=codeptr, arg1=mem, auxint=argsize, returns mem

		// Lowering pass-throughs
		{name: "LoweredNilCheck", argLength: 2, reg: regInfo{inputs: []regMask{gpspMask}}},                                                         // arg0=ptr,arg1=mem, returns void.  Faults if ptr is nil.
		{name: "LoweredGetClosurePtr", reg: regInfo{outputs: []regMask{regCtxt}}},                                                                  // scheduler ensures only at beginning of entry block
		{name: "LoweredExitProc", argLength: 2, typ: "Mem", reg: regInfo{inputs: []regMask{gpMask, 0}, clobbers: regNamed["A0"] | regNamed["A7"]}}, // arg0=mem, auxint=return code

		// F extension.
		{name: "FADDS", argLength: 2, reg: fp21, asm: "FADDS", commutative: true, typ: "Float32"},  // arg0 + arg1
		{name: "FSUBS", argLength: 2, reg: fp21, asm: "FSUBS", commutative: false, typ: "Float32"}, // arg0 - arg1
		{name: "FMULS", argLength: 2, reg: fp21, asm: "FMULS", commutative: true, typ: "Float32"},  // arg0 * arg1
		{name: "FDIVS", argLength: 2, reg: fp21, asm: "FDIVS", commutative: false, typ: "Float32"}, // arg0 / arg1
		{name: "FSQRTS", argLength: 1, reg: fp11, asm: "FDIVS", typ: "Float32"},                    // sqrt(arg0)
		{name: "FNEGS", argLength: 1, reg: fp11, asm: "FNEGS", typ: "Float32"},                     // -arg0
		{name: "FCVTSW", argLength: 1, reg: gpfp, asm: "FCVTSW", typ: "Float32"},                   // float32(arg0)
		{name: "FCVTSL", argLength: 1, reg: gpfp, asm: "FCVTSL", typ: "Float32"},                   // float32(arg0)
		{name: "FCVTWS", argLength: 1, reg: fpgp, asm: "FCVTWS", typ: "Int32"},                     // int32(arg0)
		{name: "FCVTLS", argLength: 1, reg: fpgp, asm: "FCVTLS", typ: "Int64"},                     // int64(arg0)
	}

	RISCVblocks := []blockData{
		{name: "BRANCH"}, // see flag pseudo-ops above.
	}

	archs = append(archs, arch{
		name:            "RISCV",
		pkg:             "cmd/internal/obj/riscv",
		genfile:         "../../riscv/ssa.go",
		ops:             RISCVops,
		blocks:          RISCVblocks,
		regnames:        regNamesRISCV,
		gpregmask:       gpMask,
		fpregmask:       fpMask,
		flagmask:        0,  // no flags
		framepointerreg: -1, // not used
	})
}
