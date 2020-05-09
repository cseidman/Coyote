package main

import (
	"gonum.org/v1/gonum/stat"
)

// Standard mean
var mean NativeFn = func(vm *VM, args int, argpos int) Obj {

	array := vm.Pop().(*ObjArray)
	ar := convert2FloatArray(array)
	m := stat.Mean(*ar, nil)
	return ObjFloat(m)
}

// Weighted mean
var wmean NativeFn = func(vm *VM, args int, argpos int) Obj {

	weights := vm.Pop().(*ObjArray)
	array := vm.Pop().(*ObjArray)

	ar := convert2FloatArray(array)
	wt := convert2FloatArray(weights)

	m := stat.Mean(*ar, *wt)
	return ObjFloat(m)
}

// Circular mean
// Standard mean
var cmean NativeFn = func(vm *VM, args int, argpos int) Obj {
	array := vm.Pop().(*ObjArray)
	ar := convert2FloatArray(array)
	m := stat.CircularMean(*ar, nil)
	return ObjFloat(m)
}

// Weighted mean
var wcmean NativeFn = func(vm *VM, args int, argpos int) Obj {
	weights := vm.Pop().(*ObjArray)
	array := vm.Pop().(*ObjArray)

	ar := convert2FloatArray(array)
	wt := convert2FloatArray(weights)

	m := stat.CircularMean(*ar, *wt)
	return ObjFloat(m)
}
