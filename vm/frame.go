package vm

import (
	"monc/code"
	"monc/object"
)

type Frame struct {
	cl *object.Closure
	ip int
	bp int //base pointer or frame pointer
}

func NewFrame(cl *object.Closure, bp int) *Frame {
	return &Frame{cl: cl, ip: -1, bp: bp}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
