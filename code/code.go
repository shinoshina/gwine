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
	OpNull
	OpArray
	OpHash
	OpIndex
	OpGetBuiltin

	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop

	OpTrue
	OpFalse

	OpEqual
	OpNEqual
	OpGT
	OpLT

	OpMinus
	OpBang

	OpJumpIfNotTrue
	OpJump
	OpCall
	OpReturn
	OpReturnValue

	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal
)

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpNull:     {"OpNull", []int{}},
	OpArray:    {"OpArray", []int{2}},
	OpHash:     {"OpHash", []int{2}},
	OpIndex:    {"OpIndex", []int{}},
	OpGetBuiltin: {"OpGetBuiltin",[]int{1}},

	OpAdd: {"OpAdd", []int{}},
	OpSub: {"OpSub", []int{}},
	OpMul: {"OpMul", []int{}},
	OpDiv: {"OpDiv", []int{}},
	OpPop: {"OpPop", []int{}},

	OpTrue:  {"OpTrue", []int{}},
	OpFalse: {"OpFalse", []int{}},

	OpEqual:  {"OpEqual", []int{}},
	OpNEqual: {"OpNEqual", []int{}},
	OpGT:     {"OpGT", []int{}},
	OpLT:     {"OpLT", []int{}},

	OpMinus: {"OpMinus", []int{}},
	OpBang:  {"OpBang", []int{}},

	OpJumpIfNotTrue: {"OpJumpIfNotTrue", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpCall:          {"OpCall", []int{1}},
	OpReturn:        {"OpReturn", []int{}},
	OpReturnValue:   {"OpReturnValue", []int{}},

	OpGetGlobal: {"OpGetGlobal", []int{2}},
	OpSetGlobal: {"OpSetGlobal", []int{2}},
	OpGetLocal:  {"OpGetLocal", []int{1}},
	OpSetLocal:  {"OpSetLocal", []int{1}},
}

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	length := 1
	for _, width := range def.OperandWidths {
		length += width
	}
	instruction := make([]byte, length)
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
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
func ReadUint8(ins Instructions) uint8{
	return uint8(ins[0])
}
func (ins Instructions) String() string {
	var out bytes.Buffer

	for i := 0; i < len(ins); {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s \n", err)
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
		return fmt.Sprintf("ERROR: operand len dismatch")
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR unhandled operand count %s", def.Name)
}
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}
