package main

import (
	"io/ioutil"
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
	y := vm.Pop().(*ObjInteger).Value
	x := vm.Pop().(*ObjInteger).Value
	return ObjInteger{x + y}
}

var OpenFile NativeFn = func(vm *VM, args int, argpos int) Obj {
	fileName := vm.Pop().(*ObjString).Value
	//file,err := os.Open(fileName)
	//if err !=nil {
	//	log.Panicf("Error opening file '%s'",fileName)
	//}

	class := &ObjClass{
		Fields: make(map[string]Obj),
	}

	var fnFread NativeFn = func(vm *VM, args int, argpos int) Obj {

		file, _ := ioutil.ReadFile(fileName)

		return ObjString{Value: string(file)}
	}

	class.Fields["read"] = FuncToNative(&fnFread)

	return class

}
