package main

import "fmt"

// Print operations -------------------------------------------
var Outf NativeFn = func(vm *VM, args int, argpos int) Obj {
	// If there's only one parameter the print as is
	//fmt.Printf(formattedValue(vm, args))
	argVals := make([]interface{}, args-1)
	for i := args - 1; i > 0; i-- {
		argVals[i-1] = vm.Pop().ToValue()
	}
	// The first argument is the template
	format := string(vm.Pop().(ObjString))
	fmt.Printf(format, argVals...)

	return nil
}

var Out NativeFn = func(vm *VM, args int, argpos int) Obj {
	x := vm.Pop()
	fmt.Print(x.ShowValue())
	vm.sp--
	return nil
}

var Outln NativeFn = func(vm *VM, args int, argpos int) Obj {
	x := vm.Pop()
	fmt.Println(x.ShowValue())
	vm.sp--
	return nil
}

// Supporting function
func formattedValue(vm *VM, args int) string {

	if args == 1 {
		x := vm.Pop()
		return fmt.Sprint(x.ShowValue())
	} else {
		argVals := make([]interface{}, args-1)
		for i := args - 1; i > 0; i-- {
			argVals[i-1] = vm.Pop().ToValue()
		}
		// The first argument is the template
		format := string(vm.Pop().ToBytes())
		return fmt.Sprintf(format, argVals...)
	}
}
