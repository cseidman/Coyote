package main
import ("fmt")

// Display Data Frame data ----------------------------------
var DfBrowse NativeFn = func(vm *VM, args int, argpos int) Obj {
	rowCount := int64(0)
	if args > 1 {
		rowCount = int64(vm.Pop().(ObjInteger))
	}
	df := vm.Pop().(*ObjDataFrame)
	df.PrintData(rowCount)

	return nil
}

// Database operations ---------------------------------------------
// Open a database
var OpenDatabase NativeFn = func(vm *VM, args int, argpos int) Obj {

	dbPath := string(vm.Pop().(ObjString))
	name := string(vm.Pop().(ObjString))

	vm.db = OpenDb(dbPath)
	vm.DbList[name] = vm.db

	return nil
}
// Select the context to the given database
var UseDatabase NativeFn = func(vm *VM, args int, argpos int) Obj {

	name := string(vm.Pop().(ObjString))
	if ctxdb,ok := vm.DbList[name]; ok {
		vm.db = ctxdb
	} else {
		fmt.Printf("Database '%s' does not exist. No context switch was made", name)
	}

	return nil
}




