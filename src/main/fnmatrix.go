package main

import (
	"gonum.org/v1/gonum/mat"
	"fmt"
)

type ObjMatrix struct {
	Data mat.Matrix
}

func MakeDense() {

	x := &ObjMatrix{
		Data: mat.NewDense(3, 5, nil),
	}
	fmt.Println(x)
}

var Maxtrix NativeFn = func(vm *VM, args int, argpos int) Obj {

	dataArray := vm.Pop()
	cols := vm.Pop()
	rows := vm.Pop()

	for i := args - 1; i > 0; i-- {
		argVals[i-1] = vm.Pop().ToValue()
	}

}