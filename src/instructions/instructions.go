package instructions

import (
	."../common"
	"fmt"
)

type Instr interface{
	ToBytes() []byte
	Display()
	GetByteCount() int
}

type Instruction struct {
	OpCode byte
	Operand int16
	OperandCount byte
	ByteCount int
	BytePosition int
}

func (i Instruction) ToBytes() []byte {
	bCode := []byte{i.OpCode}
	if i.OperandCount > 0 {
		bCode = append(bCode,Int16ToBytes(i.Operand)...)
	}
	return bCode
}

func (i Instruction) Display() {
	fmt.Printf("%-15s\t",OpLabel[i.OpCode])
	if i.OperandCount > 0 {
		fmt.Printf("%d",i.Operand)
	}
}

func (i Instruction) GetByteCount() int {
	return i.ByteCount
}

func (i *Instruction) SetOperand(value int16) {
	i.Operand = value
}

type Instructions struct {
	Count int
	OpCode []Instruction
	BytePosition int
	VarResult []ValueType
	Comments []string

	Constants []Obj
	ConstantsCount int16
}

func NewInstructions() Instructions {
	return Instructions {
		Count:        0,
		OpCode:       make([]Instruction,1024),
		BytePosition: 0,
		VarResult:    make([]ValueType,1024),
		Comments:     make([]string,1024),

		ConstantsCount: 0,
		Constants:    make([]Obj,16000) ,
	}
}

func (i *Instructions) AddType(varType ValueType) {
	i.VarResult[i.Count-1] = varType
}

func (i *Instructions) WriteComment(comment string) {
	i.Comments[i.Count-1] = comment
}

func (i *Instructions) WriteInstruction(opcode byte, operand int16) {
	instr := Instruction{
		OpCode:       opcode,
		Operand:      operand,
		OperandCount: 1,
		ByteCount: 3,
	}
	i.OpCode[i.Count] = instr
	i.OpCode[i.Count].BytePosition = i.BytePosition
	i.Count++
	i.BytePosition+=3

}

func (i *Instructions) WriteSimpleInstruction(opcode byte) {
	instr := Instruction{
		OpCode:       opcode,
		OperandCount: 0,
		ByteCount: 1,
	}
	i.OpCode[i.Count] = instr
	i.OpCode[i.Count].BytePosition = i.BytePosition
	i.Count++
	i.BytePosition++
}

func (i *Instructions) SetOperand(instructionNumber int, value int16) {
	i.OpCode[instructionNumber].Operand = value
}

func (i *Instructions) GetInstruction(instructionNumber int) []byte {
	return i.OpCode[instructionNumber].ToBytes()
}

func (i *Instructions) GetType(offset int) ValueType {
	return i.VarResult[i.Count-offset-1]
}

func (i *Instructions) ToByteCode() []byte {
	bCode := make([]byte,0)
	for j:=0;j<i.Count;j++ {
		bCode = append(bCode,i.OpCode[j].ToBytes()...)
	}
	return bCode
}

func (i *Instructions) Display() {
	fmt.Println("=== Instructions ===")
	bcount := 0
	for j:=0;j<i.Count;j++ {
		fmt.Printf("%04d: ",bcount)
		i.OpCode[j].Display()
		fmt.Printf("\t\t\t; %s",i.Comments[j])
		fmt.Println()
		bcount+=i.OpCode[j].GetByteCount()
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
	return i.OpCode[i.Count-1].BytePosition +  i.OpCode[i.Count-1].ByteCount
}

func (i *Instructions) JumpFrom(instrNumber int) int {
	return i.OpCode[instrNumber].BytePosition+i.OpCode[instrNumber].ByteCount
}

func (i *Instructions) JumpFromHere() int {
	return i.OpCode[i.Count-1].BytePosition+i.OpCode[i.Count-1].ByteCount
}
