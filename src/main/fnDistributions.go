package main

import (
	. "gonum.org/v1/gonum/stat/distuv"
	//rnd "golang.org/x/exp/rand"
)

// Standard correlation
var dnorm NativeFn = func(vm *VM, args int, argpos int) Obj {

	sigma := vm.Pop().(ObjFloat)
	mu := vm.Pop().(ObjFloat)
	size := int(vm.Pop().(ObjInteger))

	dist := Normal{
		Mu:    float64(mu),
		Sigma: float64(sigma),
		//Src: rnd.NewSource(0), //rnd.NewSource(0),
	}

	data := make([]Obj, size)
	for i:=0;i<size;i++{
		data[i] = ObjFloat(dist.Rand())
	}

	return ObjArray{
		ElementCount: size,
		ElementTypes: VAL_FLOAT,
		Elements:     data,
	}
}
