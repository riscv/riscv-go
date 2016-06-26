package riscv

import (
	"log"

	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/riscv"
)

// This table gives the basic information about instruction generated by the
// compiler and processed in the optimizer.  Instructions not generated may be
// omitted.
//
// NOTE(prattmic): I believe that the gc.Size flags are used only for non-SSA
// peephole optimizations, and can thus be omitted for RISCV.
var progmap = map[obj.As]obj.ProgInfo{
	obj.ATYPE:     {Flags: gc.Pseudo | gc.Skip},
	obj.ATEXT:     {Flags: gc.Pseudo},
	obj.AFUNCDATA: {Flags: gc.Pseudo},
	obj.APCDATA:   {Flags: gc.Pseudo},
	obj.AUNDEF:    {Flags: gc.Break},
	obj.AUSEFIELD: {Flags: gc.OK},
	obj.ACHECKNIL: {Flags: gc.LeftRead},
	obj.AVARDEF:   {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARKILL:  {Flags: gc.Pseudo | gc.RightWrite},
	obj.AVARLIVE:  {Flags: gc.Pseudo | gc.LeftRead},
	obj.ARET:      {Flags: gc.Break},
	obj.AJMP:      {Flags: gc.Jump | gc.Break | gc.KillCarry},

	// NOP is an internal no-op that also stands for USED and SET
	// annotations.
	obj.ANOP: {Flags: gc.LeftRead | gc.RightWrite},

	// RISCV simple three operand instructions
	riscv.AADD: {Flags: gc.LeftRead | gc.RegRead | gc.RightWrite},
	riscv.AAND: {Flags: gc.LeftRead | gc.RegRead | gc.RightWrite},
	riscv.AOR:  {Flags: gc.LeftRead | gc.RegRead | gc.RightWrite},
	riscv.ASUB: {Flags: gc.LeftRead | gc.RegRead | gc.RightWrite},
	riscv.AXOR: {Flags: gc.LeftRead | gc.RegRead | gc.RightWrite},

	// RISCV instructions
	riscv.AADDI:  {Flags: gc.LeftRead | gc.RightWrite},
	riscv.ALD:    {Flags: gc.LeftRead | gc.RightWrite | gc.Move},
	riscv.ASD:    {Flags: gc.LeftRead | gc.RightWrite | gc.Move},
	riscv.AMOV:   {Flags: gc.LeftRead | gc.RightWrite | gc.Move},
	riscv.ASEQZ:  {Flags: gc.LeftRead | gc.RightWrite},
	riscv.ASCALL: {Flags: gc.OK},

	// RISCV conditional branches
	riscv.ABEQ:  {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
	riscv.ABNE:  {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
	riscv.ABGE:  {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
	riscv.ABGEU: {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
	riscv.ABLT:  {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
	riscv.ABLTU: {Flags: gc.Cjmp | gc.LeftRead | gc.RegRead},
}

func proginfo(p *obj.Prog) {
	info, ok := progmap[p.As]
	if !ok {
		log.Printf("proginfo missing prog %s", obj.Aconv(p.As))
		return
	}

	p.Info = info
}
