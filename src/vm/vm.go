package vm

import (
	. "../common"
	. "../compiler"
	. "../value"
	"bytes"
	"fmt"
)

type VM struct {
	Position int64
	ByteCode []byte

	ip int
	sp int

	Stack     []Obj
	Globals   []Obj
	Constants []Obj

	Registers []int64

	compiler Compiler
}

func (v *VM) ReadConstant(index int16) Obj {
	return v.Constants[index]
}

func (v *VM) Push(value Obj) {
	v.Stack[v.sp] = value
	v.sp++
}

func (v *VM) Pop() Obj {
	v.sp--
	return v.Stack[v.sp]
}

func (v *VM) Peek(distance int) Obj {
	return v.Stack[v.sp-1]
}

func (v *VM) GetOperandValue() int16 {
	val := BytesToInt16(v.ByteCode[(v.ip + 1):(v.ip + 3)])
	v.ip += 2
	return val
}

func (v *VM) GetOperand() Obj {
	val := BytesToInt16(v.ByteCode[(v.ip + 1):(v.ip + 3)])
	v.ip += 2
	return v.ReadConstant(val)
}

func Exec(source *string) {

	vm := VM{
		Position:  0,
		Stack:     make([]Obj, 1024),
		Globals:   make([]Obj, 1024),
		Registers: make([]int64, 256),
		ip:        -1,
	}

	instr := Compile(source)
	vm.Constants = instr.Constants
	vm.ByteCode = instr.ToByteCode()

	vm.Interpret()

}

func (v *VM) ForLoop(fromReg int, bytes int) {
	step := v.Pop().(*ObjInteger).Value
	to := v.Pop().(*ObjInteger).Value
	v.Registers[fromReg] = v.Pop().(*ObjInteger).Value

	startIp := v.ip
	for i := v.Registers[fromReg]; i <= to; i += step {
		v.Registers[fromReg] = i
		for {
			v.ip++
			if (v.ip - startIp) > bytes {
				break
			}
			v.Dispatch()
			i = v.Registers[fromReg]
		}
		v.ip = startIp
	}
	v.ip += bytes
}

func (v *VM) Dispatch() {
	for i := 0; i < 5; i++ {
		if i >= v.sp {
			fmt.Printf("[%s] ", "  ")
		} else {
			fmt.Printf("[%2s] ", v.Stack[i].ShowValue())
		}
	}
	fmt.Print(" | ")
	//fmt.Printf("%s",strings.Repeat("\t",v.))
	fmt.Printf("%04d %s\n", v.ip, OpLabel[v.ByteCode[v.ip]])

	switch v.ByteCode[v.ip] {
	case OP_RETURN:
		return
	case OP_ICONST, OP_FCONST, OP_SCONST:
		v.Push(v.GetOperand())
	case OP_PUSH:
		v.Push(&ObjInteger{int64(v.GetOperandValue())})
	case OP_IADD:
		rval := v.Pop().(*ObjInteger).Value
		lval := v.Pop().(*ObjInteger).Value

		v.Push(&ObjInteger{Value: rval + lval})
	case OP_FADD:
		rval := v.Pop().(*ObjFloat).Value
		lval := v.Pop().(*ObjFloat).Value

		v.Push(&ObjFloat{Value: rval + lval})
	case OP_SADD:
		rval := v.Pop().(*ObjString).Value
		lval := v.Pop().(*ObjString).Value

		v.Push(&ObjString{Value: rval + lval})
	case OP_ISUBTRACT:
		rval := v.Pop().(*ObjInteger).Value
		lval := v.Pop().(*ObjInteger).Value

		v.Push(&ObjInteger{Value: lval - rval})
	case OP_FSUBTRACT:
		rval := v.Pop().(*ObjFloat).Value
		lval := v.Pop().(*ObjFloat).Value

		v.Push(&ObjFloat{Value: lval - rval})
	case OP_IMULTIPLY:
		rval := v.Pop().(*ObjInteger).Value
		lval := v.Pop().(*ObjInteger).Value

		v.Push(&ObjInteger{Value: rval * lval})
	case OP_FMULTIPLY:
		rval := v.Pop().(*ObjFloat).Value
		lval := v.Pop().(*ObjFloat).Value

		v.Push(&ObjFloat{Value: rval * lval})
	case OP_IDIVIDE:
		rval := v.Pop().(*ObjInteger).Value
		lval := v.Pop().(*ObjInteger).Value

		v.Push(&ObjInteger{Value: lval / rval})
	case OP_FDIVIDE:
		rval := v.Pop().(*ObjFloat).Value
		lval := v.Pop().(*ObjFloat).Value

		v.Push(&ObjFloat{Value: lval / rval})
	case OP_NIL:
		v.Push(&NULL{})
	case OP_INCREMENT:
		//val := v.Pop()

		//v.Push(val)
	case OP_PREINCREMENT:
		//val := v.Pop()+1
		//v.Push(val)
	case OP_GET_GLOBAL:
		idx := v.GetOperandValue()
		v.Push(v.Globals[idx])
	case OP_SET_GLOBAL:
		idx := v.GetOperandValue()
		v.Globals[idx] = v.Peek(0)
	case OP_SET_LOCAL:
		slot := v.GetOperandValue()
		v.Stack[slot] = v.Peek(0)
	case OP_GET_LOCAL:
		slot := v.GetOperandValue()
		v.Push(v.Stack[slot])
	case OP_SET_REGISTER:
		idx := v.GetOperandValue()
		v.Registers[idx] = v.Peek(0).(*ObjInteger).Value
	case OP_GET_REGISTER:
		idx := v.GetOperandValue()
		v.Push(&ObjInteger{Value: v.Registers[idx]})
	case OP_POP:
		v.sp--
	case OP_PRINT:
		fmt.Println(v.Pop().ShowValue())
	case OP_JUMP_IF_FALSE:
		val := v.Pop()
		jmpIndex := v.GetOperandValue()
		if val.(*ObjBool).Value == false {
			v.ip += int(jmpIndex)
		}

	case OP_JUMP:
		v.ip += int(v.GetOperandValue())
	case OP_FOR_LOOP:
		rgInit := int(v.Pop().(*ObjInteger).Value)
		bytevals := int(v.GetOperandValue())
		v.ForLoop(rgInit, bytevals)
	case OP_LESS:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Compare(lval, rval) == -1})
	case OP_GREATER:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Compare(lval, rval) == 1})
	case OP_NOT_EQUAL:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: !bytes.Equal(rval, lval)})
	case OP_EQUAL:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Equal(rval, lval)})

	default:
		fmt.Printf("Unhandled command: %s\n", OpLabel[v.ByteCode[v.ip]])
		return
	}
}

func (v *VM) Interpret() {

	fmt.Println("=== VM Run ===")
	codeLength := len(v.ByteCode)
	for {

		v.ip++
		if v.ip >= codeLength {
			break
		}

		v.Dispatch()

	}
}
