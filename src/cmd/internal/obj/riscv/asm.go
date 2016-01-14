// Copyright © 2015 The Go Authors.  All rights reserved.
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

import (
	"log"

	"cmd/internal/obj"
)

// TODO(bbaren)
type Optab struct {
	size int8
}

// progedit is called individually for each Prog.
// TODO(myenik)
func progedit(ctxt *obj.Link, p *obj.Prog) {
	log.Printf("progedit: ctxt: %+v p: %#v p: %s", ctxt, p, p)

	// Rewrite branches as TYPE_BRANCH
	switch p.As {
	case AJAL,
		AJALR,
		ABEQ,
		ABNE,
		ABLT,
		ABLTU,
		ABGE,
		ABGEU,
		obj.ARET,
		obj.ADUFFZERO,
		obj.ADUFFCOPY:
		if p.To.Sym != nil {
			p.To.Type = obj.TYPE_BRANCH
		}
	}
}

// TODO(myenik)
func follow(ctxt *obj.Link, s *obj.LSym) {
	log.Printf("follow: ctxt: %+v", ctxt)

	for ; s != nil; s = s.Next {
		log.Printf("s: %+v", s)
	}
}

// TODO(myenik)
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {
	log.Printf("preprocess: ctxt: %+v", ctxt)

	for ; cursym != nil; cursym = cursym.Next {
		log.Printf("cursym: %+v", cursym)
	}
}

// TODO(bbaren): Looks up an operation in the (currently nonexistent) operation
// table.
func oplook(ctxt *obj.Link, p *obj.Prog) *Optab {
	log.Printf("oplook: ctxt: %+v p: %+v", ctxt, p)
	return nil
}

// TODO(bbaren): Encodes a machine instruction.
func asmout(ctxt *obj.Link, p *obj.Prog, o *Optab) uint32 {
	log.Printf("asmout: ctxt: %+v p: %+v o: %+v", ctxt, p, o)
	return 0
}

func assemble(ctxt *obj.Link, cursym *obj.LSym) {
	log.Printf("assemble: ctxt: %+v", ctxt)

	if cursym.Text == nil || cursym.Text.Link == nil {
		// We're being asked to assemble an external function or an ELF
		// section symbol.  Do nothing.
		return
	}

	ctxt.Cursym = cursym

	// Lay out code, keeping track of how many bytes this symbol will wind
	// up using.
	cursym.Size = int64(0)
	ctxt.Autosize = int32(cursym.Text.To.Offset + 4)
	bp := cursym.P // output pointer
	for p := cursym.Text; p != nil; p = p.Link {
		ctxt.Curp = p
		ctxt.Pc = cursym.Size
		p.Pc = cursym.Size

		o := oplook(ctxt, p)
		m := o.size

		// All operations should be 32 bits wide.
		if m%4 != 0 || p.Pc%4 != 0 {
			ctxt.Diag("!pc invalid: %v size=%d", p, m)
		}

		if m == 0 {
			ctxt.Diag("zero-width instruction\n%v", p)
			continue
		}

		out := asmout(ctxt, p, o)
		bp[0] = byte(out)
		bp = bp[1:]
		bp[0] = byte(out >> 8)
		bp = bp[1:]
		bp[0] = byte(out >> 16)
		bp = bp[1:]
		bp[0] = byte(out >> 24)
		bp = bp[1:]

		cursym.Size += int64(m)
	}
}
