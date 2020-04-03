package main

import (
	"fmt"
)

type Instr interface {
	ToBytes() []byte
	Display()
	GetByteCount() int
}

type Instruction struct {
	OpCode       byte
	Operand      []byte //int16
	OperandCount byte
	ByteCount    int
	BytePosition int
	Line         int
}

func (i Instruction) ToBytes() []byte {
	bCode := []byte{i.OpCode}
	if i.OperandCount > 0 {
		//bCode = append(bCode, Int16ToBytes(i.Operand)...)
		bCode = append(bCode, i.Operand...)
	}
	return bCode
}

func (i Instruction) Display() {
	fmt.Printf(" %d ¦ ", i.Line)
	fmt.Printf("%-15s\t", OpLabel[i.OpCode])

	if i.OperandCount > 0 {

		if i.OpCode == OP_CLOSURE {
			fmt.Printf("%v", i.Operand)
		} else {
			fmt.Printf("%d", BytesToInt16(i.Operand))
		}
	}
}

func (i Instruction) GetByteCount() int {
	return i.ByteCount
}

func (i *Instruction) SetOperand(value int16) {
	i.Operand = Int16ToBytes(value)
}

type Instructions struct {
	Count          int
	OpCode         []Instruction
	BytePosition   int
	Comments       []string
	Constants      []Obj
	ConstantsCount int16
}

func NewInstructions() Instructions {

	size := 255 * 255

	return Instructions{
		Count:        0,
		OpCode:       make([]Instruction, size),
		BytePosition: 0,
		Comments:     make([]string, size),

		ConstantsCount: 0,
		Constants:      make([]Obj, 16000),
	}
}

func (i *Instructions) WriteComment(comment string) {
	i.Comments[i.Count-1] = comment
}

func (i *Instructions) WriteInstruction(opcode byte, operand int16, line int) {
	instr := Instruction{
		OpCode:       opcode,
		Operand:      Int16ToBytes(operand),
		OperandCount: 1,
		ByteCount:    3,
		Line:         line,
	}
	i.OpCode[i.Count] = instr
	i.OpCode[i.Count].BytePosition = i.BytePosition
	i.Count++
	i.BytePosition += 3

}

func (i *Instructions) WriteSimpleInstruction(opcode byte, line int) {
	instr := Instruction{
		OpCode:       opcode,
		OperandCount: 0,
		ByteCount:    1,
		Line:         line,
	}
	i.OpCode[i.Count] = instr
	i.OpCode[i.Count].BytePosition = i.BytePosition
	i.Count++
	i.BytePosition++
}

// These are to replace the existing values
func (i *Instructions) SetOperand32(instructionNumber int, value int32) {
	i.OpCode[instructionNumber].Operand = Int32ToBytes(value)
}

func (i *Instructions) SetByteOperand(instructionNumber int, value byte) {
	i.OpCode[instructionNumber].Operand = []byte{value}
}

func (i *Instructions) SetOperand(instructionNumber int, value int16) {
	i.OpCode[instructionNumber].Operand = Int16ToBytes(value)
}

// These are to "tack on" extra operands
func (i *Instructions) AddOperand32(value int32) {
	i.OpCode[i.Count-1].Operand = append(i.OpCode[i.Count-1].Operand, Int32ToBytes(value)...)
	i.BytePosition += 4
}

func (i *Instructions) AddByteOperand(value byte) {
	i.OpCode[i.Count-1].Operand = append(i.OpCode[i.Count-1].Operand, []byte{value}...)
	i.BytePosition++
}

func (i *Instructions) AddOperand(value int16) {
	i.OpCode[i.Count-1].Operand = append(i.OpCode[i.Count-1].Operand, Int16ToBytes(value)...)
	i.BytePosition += 2
}

func (i *Instructions) GetInstruction(instructionNumber int) []byte {
	return i.OpCode[instructionNumber-1].ToBytes()
}

func (i *Instructions) GetOpcode(instructionNumber int) byte {
	return i.OpCode[instructionNumber-1].OpCode
}

func (i *Instructions) ToByteCode() []byte {
	bCode := make([]byte, 0)
	for j := 0; j < i.Count; j++ {
		bCode = append(bCode, i.OpCode[j].ToBytes()...)
	}
	return bCode
}

func (i *Instructions) ToChunk() *Chunk {
	return &Chunk{
		Code:           i.ToByteCode(),
		Count:          i.BytePosition,
		Constants:      i.Constants,
		ConstantsCount: int(i.ConstantsCount),
	}
}

func (i *Instructions) Display() {

	fmt.Printf("Constants: %d\n", i.ConstantsCount)
	for c := int16(0); c < i.ConstantsCount; c++ {
		fmt.Printf("\tIndex: %d Value: %s\n", c, i.Constants[c].ShowValue())
	}

	bcount := 0
	for j := 0; j < i.Count; j++ {
		fmt.Printf("%04d: ", bcount)
		i.OpCode[j].Display()
		fmt.Printf("\t\t\t; %s", i.Comments[j])
		fmt.Println()
		bcount += i.OpCode[j].GetByteCount()
	}
}

// Calculate the number of bytes between two instructions
// Useful for any jumps and branching
func (i *Instructions) CalcByteDiff(fromInstr int, toInstr int) int {
	return i.OpCode[toInstr].BytePosition - i.OpCode[fromInstr].BytePosition
}

func (i *Instructions) CurrentBytePosition() int {
	return i.OpCode[i.Count-1].BytePosition
}

func (i *Instructions) NextBytePosition() int {
	return i.OpCode[i.Count-1].BytePosition + i.OpCode[i.Count-1].ByteCount
}

func (i *Instructions) JumpFrom(instrNumber int) int {
	return i.OpCode[instrNumber].BytePosition + i.OpCode[instrNumber].ByteCount
}

func (i *Instructions) JumpFromHere() int {
	return i.OpCode[i.Count-1].BytePosition + i.OpCode[i.Count-1].ByteCount
}
