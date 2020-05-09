package main

import (
	"gonum.org/v1/gonum/stat"
)

// Standard correlation
var correlate NativeFn = func(vm *VM, args int, argpos int) Obj {

	yarray := vm.Pop().(*ObjArray)
	yar := convert2FloatArray(yarray)

	xarray := vm.Pop().(*ObjArray)
	xar := convert2FloatArray(xarray)

	c := stat.Correlation(*xar, *yar, nil)
	return ObjFloat(c)
}

// Weighted correlation
var wcorrelate NativeFn = func(vm *VM, args int, argpos int) Obj {

	weights := vm.Pop().(*ObjArray)
	w := convert2FloatArray(weights)

	yarray := vm.Pop().(*ObjArray)
	yar := convert2FloatArray(yarray)

	xarray := vm.Pop().(*ObjArray)
	xar := convert2FloatArray(xarray)

	c := stat.Correlation(*xar, *yar, *w)
	return ObjFloat(c)
}
