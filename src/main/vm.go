package main

import (
	"bytes"
	"fmt"
	"math"
)

type CallFrame struct {
	function *ObjFunction

	ip      int
	slots   []Obj
	slotptr int
}

type VM struct {
	Frames []CallFrame
	Code   []byte
	fp     int

	Frame *CallFrame

	Stack   []Obj
	Globals []Obj

	sp     int
	prevSp int

	Registers   []int64
	ObjRegister []Obj
}

func (v *VM) GetByteCode() *[]byte {
	return &v.Frame.function.Code.Code
}

func (v *VM) PushFrame() {
	v.Frame = &v.Frames[v.fp]
	v.Frame.ip = -1
	v.fp++

}

func (v *VM) PopFrame() {
	v.fp--
	v.Frame = &v.Frames[v.fp]
	v.Code = v.Frame.function.Code.Code[:]
}

func (v *VM) ReadConstant(index int16) Obj {
	return v.Frame.function.Code.Constants[index]
}

func (v *VM) Push(value Obj) {

	v.Stack[v.sp] = value
	v.sp++
}

func (v *VM) Pop() Obj {
	v.sp--
	return v.Stack[v.sp] //v.Frame.slots[v.Frame.slotptr]
}

func (v *VM) Peek(distance int) Obj {
	return v.Stack[v.sp-distance-1] //v.Frame.slots[v.Frame.slotptr-distance-1]
}

func (v *VM) GetOperandValue() int16 {
	val := BytesToInt16(v.Code[(v.Frame.ip + 1):(v.Frame.ip + 3)])
	v.Frame.ip += 2
	return val
}

func (v *VM) GetOperand() Obj {
	val := BytesToInt16(v.Code[(v.Frame.ip + 1):(v.Frame.ip + 3)])
	v.Frame.ip += 2
	return v.ReadConstant(val)
}

/* --------------------------------------
First call into the VM
 ---------------------------------------*/
func Exec(source *string) {

	vm := VM{
		fp:        0,
		Stack:     make([]Obj, 1024),
		Globals:   make([]Obj, 1024),
		Registers: make([]int64, 256),
		Frames:    make([]CallFrame, 1024),
	}

	fn := Compile(source)

	vm.Frame = &vm.Frames[vm.fp]
	vm.Frame.function = fn
	vm.Frame.slots = vm.Stack[:]
	vm.Code = fn.Code.Code[:]
	vm.Frame.ip = -1

	vm.fp++

	vm.Interpret()

}

func (v *VM) Interpret() {

	fmt.Println("=== VM Run ===")
	var opCode byte
	for {
		v.Frame.ip++
		opCode = v.Code[v.Frame.ip]
		if opCode == OP_HALT {
			break
		}
		v.Dispatch(opCode)
	}
}

func (v *VM) Scan() {
	// Get the object we're scanning from the stack
	obj := v.Peek(0)

	switch obj.Type() {
	case VAL_ARRAY:
		v.ScanArray()
	}

}

// Handles the scan
func (v *VM) ScanArray() {

	bytes := int(v.GetOperandValue())
	counterReg := int(v.GetOperandValue())
	localIndex := int(v.GetOperandValue())

	// Initialize the register
	v.Registers[counterReg] = 0

	// Get the object with an iterator
	obj := v.Pop().(*ObjArray)

	startIp := v.Frame.ip
	stackPtr := v.sp

mainLoop:
	for i := 0; i < obj.ElementCount; i++ {

		// if there was a target variable assigned
		if localIndex != -1 {
			v.Frame.slots[localIndex] = obj.GetElement(int64(i))
		}

		v.Registers[counterReg]++

		for {
			v.Frame.ip++

			if v.Code[v.Frame.ip] == OP_BREAK {
				break mainLoop
			} else if v.Code[v.Frame.ip] == OP_CONTINUE {
				v.Frame.ip = startIp
				continue mainLoop
			}
			v.Dispatch(v.Code[v.Frame.ip])
		}
		// Get back to where we were
		v.Frame.ip = startIp
	}
	v.Frame.ip = startIp + bytes
	v.sp = stackPtr
}

// This handles the FOR loop
func (v *VM) ForLoop() {

	fromReg := int(v.Pop().(*ObjInteger).Value)
	bytes := int(v.GetOperandValue())

	step := v.Pop().(*ObjInteger).Value
	to := v.Pop().(*ObjInteger).Value
	fromVal := v.Pop().(*ObjInteger).Value
	// We're positioned at the current instruction
	startIp := v.Frame.ip
	stackPtr := v.sp

mainLoop:
	for i := fromVal; i <= to; i += step {
		// Update the register we use as a variable
		v.Registers[fromReg] = i
		for {
			// Get next instruction
			v.Frame.ip++

			if v.Code[v.Frame.ip] == OP_BREAK {
				break mainLoop
			} else if v.Code[v.Frame.ip] == OP_CONTINUE {
				v.Frame.ip = startIp
				continue mainLoop
			}
			// Execute the instruction
			v.Dispatch(v.Code[v.Frame.ip])
		}
		// Get back to where we were
		v.Frame.ip = startIp
	}
	v.Frame.ip = startIp + bytes
	v.sp = stackPtr
}

func (v *VM) DebugInfo(opCode byte) {

	for i := 0; i < 8; i++ {
		// Loop over the slots of the current frame
		if i >= v.sp {
			// If there are less than 5 elements in the stack then just print blanks
			fmt.Printf("[%s] ", "    ")
		} else {
			// Print the value of the current frame
			if v.Stack[i] == nil {
				fmt.Printf("[%s] ", "null")
			} else {
				fmt.Printf("[%4s] ", v.Stack[i].ShowValue())
			}
		}
	}

	fmt.Print(" | ")
	fmt.Printf("%d:%d:%d", v.sp, v.Frame.slotptr, v.fp-1)
	fmt.Print(" | ")
	fmt.Printf("%04d %s", v.Frame.ip, OpLabel[opCode])
	switch opCode {
	case OP_GET_LOCAL:
		slot := BytesToInt16(v.Code[(v.Frame.ip + 1):(v.Frame.ip + 3)])
		fmt.Printf("\t[%d]", slot)
	case OP_CALL:
		fmt.Println()
	}
	fmt.Println()

}

/* -------------------------------------------------------
Main dispatch loop
 --------------------------------------------------------*/
func (v *VM) Dispatch(opCode byte) {

	v.DebugInfo(opCode)

	// Main switch
	switch opCode {
	case OP_RETURN:
		// Before this frames disappears, get the
		// return value on the stack
		result := v.Pop()

		// If this is the first element on the stack, then
		// just clear and return
		v.fp--
		if v.fp == 0 {
			v.Pop()
			return
		}
		// The stack pointer

		// Push the results we got above on to this frames stack
		prev := v.Frame.slotptr // Get the previous frame's stack position
		v.Frame = &v.Frames[v.fp-1]
		v.sp = prev
		v.Code = v.Frame.function.Code.Code[:]

		v.Push(result)
		return

	case OP_HALT:
		return
	case OP_CALL:
		// Get the parameters
		argCount := int(v.GetOperandValue())
		function := v.Peek(argCount)

		// Push the code into this new frame
		v.Frame = &v.Frames[v.fp]
		v.Frame.ip = -1
		v.fp++

		v.Frame.function = function.(*ObjFunction)
		// Reference to the code block
		v.Code = v.Frame.function.Code.Code[:]
		// This frame's slots line up with the stacks at the point where
		// the function and the parameters begin on the stack
		start := v.sp - argCount - 1
		//fmt.Printf("Start: %d STackValue: %v\n",start,v.Stack[start].ShowValue())
		v.Frame.slots = v.Stack[start:]
		v.Frame.slotptr = start

	case OP_ICONST, OP_FCONST, OP_SCONST, OP_FN_CONST:
		v.Push(v.GetOperand())

	case OP_PUSH:
		v.Push(&ObjInteger{int64(v.GetOperandValue())})

	case OP_PUSH_0:
		v.Push(&ObjInteger{0})

	case OP_PUSH_1:
		v.Push(&ObjInteger{1})

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

	case OP_INEGATE:
		val := -v.Pop().(*ObjInteger).Value
		v.Push(&ObjInteger{Value: -val})

	case OP_FNEGATE:
		val := -v.Pop().(*ObjFloat).Value
		v.Push(&ObjFloat{Value: -val})

	case OP_GET_GLOBAL:
		idx := v.GetOperandValue()
		v.Push(v.Globals[idx])

	case OP_SET_GLOBAL:
		idx := v.GetOperandValue()
		v.Globals[idx] = v.Pop() //v.Peek(0)

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

	case OP_GET_ALOCAL:
		elem := v.Pop().(*ObjInteger).Value
		slot := v.GetOperandValue()
		v.Push(v.Frame.slots[slot].(*ObjArray).Elements[elem])

	case OP_SET_ALOCAL:
		slot := v.GetOperandValue()
		val := v.Peek(0)
		elem := v.Pop().(*ObjInteger).Value
		v.Stack[slot].(*ObjArray).Elements[elem] = val

	case OP_GET_AGLOBAL:
		elem := v.Pop().(*ObjInteger).Value
		idx := v.GetOperandValue()
		v.Push(v.Globals[idx].(*ObjArray).Elements[elem])

	case OP_SET_AGLOBAL:
		idx := v.GetOperandValue()
		val := v.Pop()
		elem := int(v.Peek(0).(*ObjInteger).Value)
		v.Globals[idx].(*ObjArray).Elements[elem] = val

	case OP_POP:
		v.sp--

	case OP_PRINT:
		fmt.Println(v.Pop().ShowValue())

	case OP_JUMP_IF_FALSE:
		val := v.Peek(0)
		jmpIndex := v.GetOperandValue()
		if val.(*ObjBool).Value == false {
			v.Frame.ip += int(jmpIndex)
		}

	case OP_JUMP:
		v.Frame.ip += int(v.GetOperandValue())

	case OP_FOR_LOOP:
		v.ForLoop()

	case OP_LESS:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Compare(lval, rval) < 0})

	case OP_LESS_EQUAL:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Compare(lval, rval) <= 0})

	case OP_GREATER:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Compare(lval, rval) > 0})

	case OP_NOT_EQUAL:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: !bytes.Equal(rval, lval)})

	case OP_EQUAL:
		rval := v.Pop().ToBytes()
		lval := v.Pop().ToBytes()

		v.Push(&ObjBool{Value: bytes.Equal(rval, lval)})

	case OP_IEXP:
		pwr := v.Pop().(*ObjInteger).Value
		lval := v.Pop().(*ObjInteger).Value

		v.Push(&ObjInteger{Value: int64(math.Pow(float64(lval), float64(pwr)))})

	case OP_FEXP:
		pwr := v.Pop().(*ObjFloat).Value
		lval := v.Pop().(*ObjFloat).Value

		v.Push(&ObjFloat{Value: math.Pow(lval, pwr)})

	case OP_TRUE:
		v.Push(&ObjBool{Value: true})

	case OP_FALSE:
		v.Push(&ObjBool{Value: false})

	case OP_CONTINUE:
	case OP_ARRAY:
		elements := v.Pop().(*ObjInteger).Value
		dType := byte(v.GetOperandValue())

		o := make([]Obj, elements)
		for i := elements - 1; i >= 0; i-- {
			o[i] = v.Pop()
		}
		v.Push(&ObjArray{
			ElementCount: int(elements),
			ElementTypes: ValueType(dType),
			Elements:     o,
		})
	case OP_SCAN:
		v.Scan()
	default:
		fmt.Printf("Unhandled command: %s\n", OpLabel[(*v.GetByteCode())[v.Frame.ip]])
		return
	}
}
