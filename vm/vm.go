package vm

import (
	"fmt"
	"monc/code"
	"monc/compiler"
	"monc/object"
)

const StackSize = 2048
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int // pointer to next value. stack top is stack[sp-1]
	globals      []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		sp:           0,
		globals:      make([]object.Object, GlobalSize),
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpHash:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements
			err = vm.push(hash)
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements
			err := vm.push(array)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpConstant:
			// ReadUint16 is faster than ReadOperands
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.OpAdd, code.OpMul, code.OpDiv, code.OpSub:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}

		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1

		case code.OpJumpIf:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpNull:
			if err := vm.push(Null); err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) buildHash(head, tail int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := head; i < tail; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]
		pair := object.HashPair{Key: key, Value: value}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}
		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) buildArray(head, tail int) object.Object {
	elements := make([]object.Object, tail-head)

	for i := head; i < tail; i++ {
		elements[i-head] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func (vm *VM) executeMinusOperator() error {
	oprnd := vm.pop()

	if oprnd.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", oprnd.Type())
	}

	value := oprnd.(*object.Integer).Value

	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	/* We consider every object other than False as a truthy value
	 * and therefore push a False on the stack
	 */

	switch operand {
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(rightVal == leftVal))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(rightVal != leftVal))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	} else {
		return False
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	switch {

	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return vm.execBiIntOp(op, left, right)

	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return vm.execBiStrOP(op, left, right)

	default:
		return fmt.Errorf("unsupported tepes for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) execBiStrOP(op code.Opcode, lObj, rObj object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	lVal := lObj.(*object.String).Value
	rVal := rObj.(*object.String).Value

	return vm.push(&object.String{Value: lVal + rVal})
}

func (vm *VM) execBiIntOp(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal

	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

// for testing OpPop
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}
