package main

var FunctionRegister = make(map[string]*ObjNative)

func RegisterNative(name string, ofn NativeFn, returnData ExpressionData, hasReturnValue bool) {
	FunctionRegister[name] = NewNative(&ofn)
	FunctionRegister[name].ReturnType = returnData
	FunctionRegister[name].hasReturn = hasReturnValue
}

func RegisterFunctions() {
	RegisterNative("OpenFile", OpenFile, ExpressionData{VAL_INTEGER, VAR_SCALAR,0}, true)
	RegisterNative("print", Out, ExpressionData{VAL_INTEGER, VAR_UNKNOWN,0},false)
	RegisterNative("println", Outln, ExpressionData{VAL_NIL, VAR_UNKNOWN,0},false)
	RegisterNative("printf", Outf, ExpressionData{VAL_NIL, VAR_UNKNOWN,0}, false)
	RegisterNative("Matrix", Matrix, ExpressionData{VAL_CLASS, VAR_CLASS,0}, true)
	RegisterNative("newarray", array, ExpressionData{VAL_ARRAY, VAR_ARRAY,0}, true)
	RegisterNative("mean", mean, ExpressionData{VAL_FLOAT, VAR_SCALAR,0}, true)
	RegisterNative("wmean", wmean, ExpressionData{VAL_FLOAT, VAR_SCALAR,0}, true)
	RegisterNative("transpose", Transpose, ExpressionData{VAL_MATRIX, VAR_MATRIX,0}, true)
	// Stats - Distribution
	RegisterNative("dnorm", dnorm, ExpressionData{VAL_FLOAT, VAR_ARRAY,0}, true)
	// Dataframe and database
	RegisterNative("showdata", DfBrowse, ExpressionData{VAL_NIL, VAR_SCALAR,0}, false)
	RegisterNative("opendb", OpenDatabase, ExpressionData{VAL_NIL, VAR_SCALAR,0}, false)
	RegisterNative("use", UseDatabase, ExpressionData{VAL_NIL, VAR_SCALAR,0}, false)
}

func ResolveNativeFunction(name string) *ObjNative {
	if val, ok := FunctionRegister[name]; ok {
		return val
	}
	return nil
}

func FuncToNative(fn *NativeFn) ObjNative {
	return ObjNative{Function: NewNative(fn).Function}
}

// Supporting functions
func convert2FloatArray(array *ObjArray) *[]float64 {
	ar := make([]float64, array.ElementCount)
	for i := 0; i < array.ElementCount; i++ {
		ar[i] = float64(array.Elements[i].(ObjFloat))
	}
	return &ar
}
