// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riscv

import (
	"cmd/compile/internal/gc"
	"cmd/internal/obj"
	"cmd/internal/obj/riscv"
)

func ginsnop() {
	// Hardware nop is ADD $0, ZERO
	p := gc.Prog(riscv.AADD)
	p.From.Type = obj.TYPE_CONST
	p.From.Offset = 0
	p.From3 = &obj.Addr{Type: obj.TYPE_REG, Reg: riscv.REG_ZERO}
	p.To = *p.From3
}
