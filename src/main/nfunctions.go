package main

import (
	"fmt"
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
	RegisterNative("print", Out)
	RegisterNative("println", Outln)
	RegisterNative("printf", Outf)
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

// Print operations -------------------------------------------
var Outf NativeFn = func(vm *VM, args int, argpos int) Obj {
	// If there's only one parameter the print as is
	fmt.Print(formattedValue(vm, args))
	return nil
}

var Out NativeFn = func(vm *VM, args int, argpos int) Obj {
	x := vm.Pop()
	fmt.Print(x.ShowValue())
	return nil
}

var Outln NativeFn = func(vm *VM, args int, argpos int) Obj {
	x := vm.Pop()
	fmt.Println(x.ShowValue())
	return nil
}

// Supporting function
func formattedValue(vm *VM, args int) string {
	if args == 0 {
		x := vm.Pop()
		return fmt.Sprint(x.ShowValue())
	} else {
		argVals := make([]interface{}, args-1)
		for i := args - 1; i > 0; i-- {
			argVals[i-1] = vm.Pop().ToValue()
		}
		// The first argument is the template
		format := string(vm.Pop().(ObjString))
		return fmt.Sprintf(format, argVals...)
	}
}

// File operations ----------------------------------------------
var OpenFile NativeFn = func(vm *VM, args int, argpos int) Obj {
	fileName := string(vm.Pop().(ObjString))
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
		file, _ := ioutil.ReadFile(fileName)
		return ObjString{Value: string(file)}
	}
	class.Fields["read"] = FuncToNative(&fnFread)

	// readall: Returns the full contents of the file
	var fnFreadAll NativeFn = func(vm *VM, args int, argpos int) Obj {
		file, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Panicf("Error reading file '%s'", fileName)
		}
		return ObjString(file)
	}
	class.Fields["readall"] = FuncToNative(&fnFreadAll)

	return class

}
