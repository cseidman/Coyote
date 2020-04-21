package main

import (
	"io/ioutil"
	"log"
	"os"
)

var FunctionRegister = make(map[string]*ObjNative)

func RegisterNative(name string, ofn NativeFn) {
	FunctionRegister[name] = NewNative(&ofn)
}

func RegisterFunctions() {
	RegisterNative("Add", AddIt)
	RegisterNative("OpenFile", OpenFile)
}

func ResolveNativeFunction(name string) *ObjNative {
	if val, ok := FunctionRegister[name]; ok {
		return val
	}
	return nil
}

func FuncToNative(fn *NativeFn) ObjNative {
	return ObjNative{Function: NewNative(fn).Function}
}

// Function definitions ************************
var AddIt NativeFn = func(vm *VM, args int, argpos int) Obj {
	y := vm.Pop().(ObjInteger)
	x := vm.Pop().(ObjInteger)

	return x + y
}

var OpenFile NativeFn = func(vm *VM, args int, argpos int) Obj {
	fileName := vm.Pop().(*ObjString).Value
	file, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Error opening file '%s'", fileName)
	}

	class := &ObjClass{
		Fields: make(map[string]Obj),
	}
	// Load some properties
	class.Fields["position"] = ObjInteger(0)

	// Build the native methods here ************************

	// read(<start:int>, <bytes:int>) returns []bytes

	var fnFread NativeFn = func(vm *VM, args int, argpos int) Obj {

		byteCount := int64(vm.Pop().(ObjInteger))
		startFrom := int64(vm.Pop().(ObjInteger))

		b := make([]byte, byteCount)
		_, err := file.ReadAt(b, startFrom)
		if err != nil {
			log.Panicf("Error reading file '%s'", fileName)
		}
		arObj := make([]Obj, byteCount)
		for i := int64(0); i < byteCount; i++ {
			arObj[byteCount-i-1] = oByte(b[i])
		}
		return ObjArray{
			ElementCount: int(byteCount),
			ElementTypes: VAL_BYTE,
			Elements:     arObj,
		}
	}
	class.Fields["read"] = FuncToNative(&fnFread)

	// readall: Returns the full contents of the file
	var fnFreadAll NativeFn = func(vm *VM, args int, argpos int) Obj {
		file, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Panicf("Error reading file '%s'", fileName)
		}
		return ObjString{Value: string(file)}
	}
	class.Fields["readall"] = FuncToNative(&fnFreadAll)

	return class

}
