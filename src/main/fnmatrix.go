package main

import (
	"gonum.org/v1/gonum/mat"
)

type ObjMatrix struct {
	Data mat.Matrix
}

func (o ObjMatrix) ShowValue() string {
	panic("implement me")
}

func (o ObjMatrix) Type() ValueType {
	panic("implement me")
}

func (o ObjMatrix) ToBytes() []byte {
	panic("implement me")
}

func (o ObjMatrix) ToValue() interface{} {
	panic("implement me")
}

var Matrix NativeFn = func(vm *VM, args int, argpos int) Obj {

	dataArray := vm.Pop().(ObjArray)
	cols := int(vm.Pop().(ObjInteger))
	rows := int(vm.Pop().(ObjInteger))

	ar := make([]float64, dataArray.ElementCount)
	for i := 0; i < dataArray.ElementCount; i++ {
		ar[i] = float64(dataArray.Elements[i].(ObjFloat))
	}
	return ObjMatrix{
		Data: mat.NewDense(rows, cols, ar),
	}
}
