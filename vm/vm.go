package vm

import (
	"fmt"
	"gwine/code"
	"gwine/compiler"
	"gwine/object"
)

const StackSize = 2048
const GlobalsSize = 65536
const MaxFrames = 1024

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int // offset of top object + 1 , 0 refers  stack empty

	globals []object.Object

	frames     []*Frame
	frameIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFm := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFm

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:     frames,
		frameIndex: 1,
	}
}
func NewWithGlobalStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFm := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFm
	return &VM{
		constants: bytecode.Constants,

		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    s,
		frames:     frames,
		frameIndex: 1,
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

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}
func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.frameIndex] = f
	vm.frameIndex++
}
func (vm *VM) popFrame() *Frame {
	if vm.frameIndex == 0 {
		return nil
	}
	vm.frameIndex--
	return vm.frames[vm.frameIndex]
}
func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			index := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.constants[index])
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(object.NullObj)
			if err != nil {
				return err
			}
		case code.OpArray:
			numElements := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			array := vm.buildArray(vm.sp-int(numElements), vm.sp)
			vm.sp -= int(numElements)

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElements := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			hash, err := vm.buildHash(vm.sp-int(numElements), vm.sp)
			if err != nil {
				return err
			}
			vm.sp -= int(numElements)
			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNEqual, code.OpGT, code.OpLT:
			err := vm.executeComparison(op)
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
		case code.OpJump:
			jumpto := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = jumpto - 1
		case code.OpJumpIfNotTrue:
			jumpto := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			condition := vm.pop()
			if !isTrue(condition) {
				vm.currentFrame().ip = jumpto - 1
			}
		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			fn, ok := vm.stack[vm.sp-1-int(numArgs)].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non function")
			}
			if int(numArgs) != fn.NumParameters {
				return fmt.Errorf("wrong number of argument")
			}
			frame := NewFrame(fn, vm.sp-int(numArgs))
			vm.pushFrame(frame)
			vm.sp = frame.basePointer + fn.NumLocals
		case code.OpReturnValue:
			rv := vm.pop()
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(rv)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1
			err := vm.push(object.NullObj)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			vm.stack[vm.currentFrame().basePointer+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1
			err := vm.push(vm.stack[vm.currentFrame().basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (vm *VM) buildArray(start, end int) object.Object {
	eles := make([]object.Object, end-start)
	for i := start; i < end; i++ {
		eles[i-start] = vm.stack[i]
	}
	return &object.Array{Elements: eles}
}
func (vm *VM) buildHash(start, end int) (object.Object, error) {
	pairs := make(map[object.HashKey]object.HashPair)

	for i := start; i < end; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key %s", key.Type())
		}
		pairs[hashKey.HashKey()] = pair
	}
	return &object.Hash{Pairs: pairs}, nil
}
func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		array := left.(*object.Array)
		index := index.(*object.Integer).Value
		if index < 0 || int(index) > len(array.Elements)-1 {
			return vm.push(object.NullObj)
		}
		return vm.push(array.Elements[index])
	case left.Type() == object.HASH_OBJ:
		hashmap := left.(*object.Hash)
		key, ok := index.(object.Hashable)
		if !ok {
			return fmt.Errorf("unusable as hash key %s", index.Type())
		}

		pair, ok := hashmap.Pairs[key.HashKey()]
		if !ok {
			return vm.push(object.NullObj)
		}
		return vm.push(pair.Value)
	default:
		return fmt.Errorf("index operator not supported %s", left.Type())
	}
}
func (vm *VM) executeBinaryOperation(op code.Opcode) error {

	r := vm.pop()
	l := vm.pop()

	switch {
	case l.Type() == object.INTEGER_OBJ && r.Type() == object.INTEGER_OBJ:
		return vm.executeBinaryIntegerOperation(op, l, r)
	case l.Type() == object.STRING_OBJ && r.Type() == object.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, l, r)
	default:
		return fmt.Errorf("operator %d operand %s , %s type dismatch", op, r.Type(), l.Type())
	}
}
func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value

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
func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {

	if op != code.OpAdd {
		return fmt.Errorf("unknown operator %d", op)
	}
	lv := left.(*object.String).Value
	rv := right.(*object.String).Value

	return vm.push(&object.String{Value: lv + rv})

}
func (vm *VM) executeComparison(op code.Opcode) error {

	r := vm.pop()
	l := vm.pop()

	if l.Type() == object.INTEGER_OBJ && r.Type() == object.INTEGER_OBJ {
		lv := l.(*object.Integer).Value
		rv := r.(*object.Integer).Value

		switch op {
		case code.OpEqual:
			return vm.push(nativeBoolToBooleanObject(lv == rv))
		case code.OpNEqual:
			return vm.push(nativeBoolToBooleanObject(lv != rv))
		case code.OpGT:
			return vm.push(nativeBoolToBooleanObject(lv > rv))
		case code.OpLT:
			return vm.push(nativeBoolToBooleanObject(lv < rv))
		default:
			return fmt.Errorf("unknown operator %d", op)
		}
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(l == r))
	case code.OpNEqual:
		return vm.push(nativeBoolToBooleanObject(l != r))
	default:
		return fmt.Errorf("unknown operator %d", op)
	}
}
func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case object.True:
		return vm.push(object.False)
	case object.False:
		return vm.push(object.True)
	case object.NullObj:
		return vm.push(object.True)
	default:
		return vm.push(object.False)
	}
}
func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type %s for -", operand.Type())
	}
	v := operand.(*object.Integer).Value
	return vm.push(&object.Integer{Value: -v})
}
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return object.True
	} else {
		return object.False
	}
}
func isTrue(obj object.Object) bool {

	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}

}
