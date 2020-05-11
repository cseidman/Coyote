package main

import (
	"gonum.org/v1/gonum/mat"
)

type ObjMatrix struct {
	Data mat.Matrix
	Rows int
	Cols int
}

func (o ObjMatrix) ShowValue() string    { return "<matrix>" }
func (o ObjMatrix) Type() ValueType      { return VAL_CLASS }
func (o ObjMatrix) ToBytes() []byte      { panic("implement me") }
func (o ObjMatrix) ToValue() interface{} { return "<matrix>" }
func (o ObjMatrix) Print() string {
	return "<nmatrix>"
}

var Matrix NativeFn = func(vm *VM, args int, argpos int) Obj {

	dataArray := vm.Pop().(*ObjArray)
	cols := int(vm.Pop().(ObjInteger))
	rows := int(vm.Pop().(ObjInteger))

	ar := convert2FloatArray(dataArray)

	return ObjMatrix{
		Rows: rows,
		Cols: cols,
		Data: mat.NewDense(rows, cols, *ar),
	}
}

var Transpose NativeFn = func(vm *VM, args int, argpos int) Obj {
	m :=  vm.Pop().(ObjMatrix)
	return ObjMatrix {
		Cols: m.Rows,
		Rows: m.Cols,
		Data : m.Data.T(),
	}

}
