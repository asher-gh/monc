package vm

import (
	"monc/code"
	"monc/object"
)

type Frame struct {
	fn *object.CompiledFn
	ip int
	bp int //base pointer or frame pointer
}

func NewFrame(fn *object.CompiledFn, bp int) *Frame {
	return &Frame{fn: fn, ip: -1, bp: bp}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
