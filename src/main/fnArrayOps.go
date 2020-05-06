package main

var array NativeFn = func(vm *VM, args int, argpos int) Obj {

	dtype := ValueType(vm.Pop().(ObjInteger))
	count := int(vm.Pop().(ObjInteger))

	ar := make([]Obj, count)
	for i := 0; i < count; i++ {
		ar[i] = &NULL{}
	}

	return ObjArray{
		ElementCount: count,
		ElementTypes: dtype,
		Elements:     ar,
	}

}
