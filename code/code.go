package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte
type Definition struct {
	Name          string
	OperandWidths []int
}

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpIf
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpClosure
)

var definitions = map[Opcode]*Definition{
	OpConstant:    {"OpConstant", []int{2}},
	OpAdd:         {"OpAdd", []int{}},   // addition
	OpPop:         {"OpPop", []int{}},   // pop the topmost element
	OpSub:         {"OpSub", []int{}},   // subtraction
	OpMul:         {"OpMul", []int{}},   // multiplication
	OpDiv:         {"OpDiv", []int{}},   // division
	OpTrue:        {"OpTrue", []int{}},  // push true to stack
	OpFalse:       {"OpFalse", []int{}}, // push false to stack
	OpEqual:       {"OpEqual", []int{}},
	OpNotEqual:    {"OpNotEqual", []int{}},
	OpGreaterThan: {"OpGreaterThan", []int{}},
	OpMinus:       {"OpMinus", []int{}},
	OpBang:        {"OpBang", []int{}},
	OpJump:        {"OpJump", []int{2}},
	OpJumpIf:      {"OpJumpIf", []int{2}}, // jump over if stack top is not truthy
	OpNull:        {"OpNull", []int{}},    // put vm.Null on the stack
	OpSetGlobal:   {"OpSetGlobal", []int{2}},
	OpGetGlobal:   {"OpGetGlobal", []int{2}},
	OpArray:       {"OpArray", []int{2}}, // operand is the array length
	OpHash:        {"OpHash", []int{2}},  // operand specifies the number of keys and values
	OpIndex:       {"OpIndex", []int{}},
	OpCall:        {"OpCall", []int{1}},
	OpReturnValue: {"OpReturnValue", []int{}},
	OpReturn:      {"OpReturn", []int{}},
	OpSetLocal:    {"OpSetLocal", []int{1}},
	OpGetLocal:    {"OpGetLocal", []int{1}},
	OpGetBuiltin:  {"OpGetBuiltin", []int{1}},
	/*
	   OpClosure has two operands
	   - 2 bytes wide constant index: specifies where in the constant pool
	     we can find the *object.CompiledFn
	   - 1 byte wide count: specifies how many free variables sit on the
	     stack
	*/
	OpClosure: {"OpClosure", []int{2, 1}},
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n ", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make encodes into a bytecode instruction
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}

// ReadOperands decodes the operands of the encoded bytecode instruction by [code.Make]
// and returns a slice of operands and the number of bytes taken by the operands
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(instructions Instructions) uint16 {
	return binary.BigEndian.Uint16(instructions)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
