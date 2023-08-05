package vm

import (
	"monc/code"
	"monc/object"
)

type Frame struct {
	fn *object.CompiledFn
	ip int
}

func NewFrame(fn *object.CompiledFn) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
