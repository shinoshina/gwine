package vm

import (
	"fmt"
	"gwine/code"
	"gwine/compiler"
	"gwine/object"
)

const StackSize = 2048

type VM struct {
	instructions code.Instructions
	constants    []object.Object

	stack []object.Object
	sp    int // offset of top object + 1 , 0 refers  stack empty
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}
func (vm *VM) Top() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}
func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}
func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	vm.sp--
	return vm.stack[vm.sp]

}
func (vm *VM) LastPoped() object.Object {
	return vm.stack[vm.sp]
}
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			index := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[index])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual ,code.OpNEqual,code.OpGT,code.OpLT:
			err := vm.executeComparison(op)
			if err != nil{
				return err
			}
		case code.OpPop:
			vm.pop()
		case code.OpTrue:
			err := vm.push(object.True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(object.False)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (vm *VM) executeBinaryOperation(op code.Opcode) error {

	r := vm.pop()
	l := vm.pop()

	if r.Type() != object.INTEGER_OBJ || l.Type() != object.INTEGER_OBJ {

		return fmt.Errorf("operator %d operand %s , %s type dismatch", op, r.Type(), l.Type())
	}

	rv := r.(*object.Integer).Value
	lv := l.(*object.Integer).Value

	var result int64
	switch op {
	case code.OpAdd:
		result = lv + rv
	case code.OpSub:
		result = lv - rv
	case code.OpMul:
		result = lv * rv
	case code.OpDiv:
		result = lv / rv
	default:
		return fmt.Errorf("unknown operator %d", op)
	}
	return vm.push(&object.Integer{Value: result})
}
func (vm *VM) executeComparison(op code.Opcode) error{

	r := vm.pop()
	l := vm.pop()


	if l.Type() == object.INTEGER_OBJ && r.Type() == object.INTEGER_OBJ {
		lv := l.(*object.Integer).Value
		rv := r.(*object.Integer).Value
		
		switch op{
		case code.OpEqual:
			return vm.push(nativeBoolToBooleanObject(lv == rv))
		case code.OpNEqual:
			return vm.push(nativeBoolToBooleanObject(lv != rv))
		case code.OpGT:
			return vm.push(nativeBoolToBooleanObject(lv > rv))
		case code.OpLT:
			return vm.push(nativeBoolToBooleanObject(lv < rv))
		default:
			return fmt.Errorf("unknown operator %d",op)
		}
	}

	switch op{
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(l == r))
	case code.OpNEqual:
		return vm.push(nativeBoolToBooleanObject(l != r))
	default:
		return fmt.Errorf("unknown operator %d",op)
	}


}
func nativeBoolToBooleanObject(input bool) *object.Boolean{
	if input {
		return object.True
	}else {
		return object.False
	}
}
