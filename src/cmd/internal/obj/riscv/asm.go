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

// preprocess is responsible for:
// * Updating the SP on function entry and exit
// * Rewriting RET to a real return instruction
func preprocess(ctxt *obj.Link, cursym *obj.LSym) {
	log.Printf("preprocess: ctxt: %+v", ctxt)

	ctxt.Cursym = cursym

	if cursym.Text == nil || cursym.Text.Link == nil {
		return
	}

	stackSize := cursym.Text.To.Offset

	// TODO(prattmic): explain what these are really for,
	// once I figure it out.
	cursym.Args = cursym.Text.To.Val.(int32)
	cursym.Locals = int32(stackSize)

	var q *obj.Prog
	for p := cursym.Text; p != nil; p = p.Link {
		log.Printf("p: %+v", p)

		switch p.As {
		case obj.ATEXT:
			// Function entry. Setup stack.
			// TODO(prattmic): handle calls to morestack.
			q = p
			q = obj.Appendp(ctxt, q)
			q.As = AADDI
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_SP
			q.From3 = &obj.Addr{}
			q.From3.Type = obj.TYPE_CONST
			q.From3.Offset = -stackSize
			q.To.Type = obj.TYPE_REG
			q.To.Reg = REG_SP
			q.Spadj = -stackSize
			// TODO(prattmic): Other fields, like Reg?
		case obj.ARET:
			// Function exit. Stack teardown and exit.
			q = p
			q = obj.Appendp(ctxt, q)
			q.As = AADDI
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_SP
			q.From3 = &obj.Addr{}
			q.From3.Type = obj.TYPE_CONST
			q.From3.Offset = stackSize
			q.To.Type = obj.TYPE_REG
			q.To.Reg = REG_SP
			q.Spadj = stackSize
			// TODO(prattmic): Other fields, like Reg?

			q = obj.Appendp(ctxt, q)
			q.As = AJAL
			q.From.Type = obj.TYPE_REG
			q.From.Reg = REG_RA
			q.To.Type = obj.TYPE_REG
			q.To.Reg = REG_ZERO
			// TODO(prattmic): Other fields, like Reg?
		}
	}
}

// TODO(myenik)
func assemble(ctxt *obj.Link, cursym *obj.LSym) {
	log.Printf("assemble: ctxt: %+v", ctxt)

	for ; cursym != nil; cursym = cursym.Next {
		log.Printf("cursym: %+v", cursym)
	}
}
