package main

import (
	"io/ioutil"
	"log"
	"os"
)

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
		return ObjString(string(file))
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
