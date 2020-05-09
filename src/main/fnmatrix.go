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

	ar := make([]float64, dataArray.ElementCount)
	for i := 0; i < dataArray.ElementCount; i++ {
		ar[i] = float64(dataArray.Elements[i].(ObjFloat))
	}
	mat := ObjMatrix{
		Rows: rows,
		Cols: cols,
		Data: mat.NewDense(rows, cols, ar),
	}

	oClass := &ObjClass{
		Fields: make(map[string]Obj),
	}

	oClass.Fields["matrix"] = mat
	oClass.FieldCount++

	// Transpose
	var fnTranspose NativeFn = func(vm *VM, args int, argpos int) Obj {
		m := new(ObjMatrix)
		return *m
	}
	oClass.Fields["Transpose"] = FuncToNative(&fnTranspose)
	oClass.FieldCount++

	return oClass
}
