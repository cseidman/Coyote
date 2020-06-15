package main

// We use this to assign a unique if to every scope defined in the application
// so that we can know later on that variables declared in a given scope are different
// to variables of the same name are indeed different despite being in the same depth
var ScopeId = -1

// Keeps track of the last expression's return value
type ExpressionData struct {
	Value   ValueType
	ObjType VarType
	//DataType string
	Dimensions int // Relevant only for arrays and matrices
}

var ExpressionValue = make([]ExpressionData, 255)
var ExpressionValueId int

func PushExpressionValue(data ExpressionData) {
	ExpressionValue[ExpressionValueId] = data

	ExpressionValueId++
}

func PopExpressionValue() ExpressionData {
	ExpressionValueId--
	return ExpressionValue[ExpressionValueId]
}

type VarTable struct {
	Symbol []Variable
}

type Variable interface {
	GetScopeType() VariableScope
}

// Describes properties
type PropertyVar struct {
	EnclosingClass *ClassVar
	Access         AccessorType
	Name           string
	Index          int16
	ExprData       ExpressionData
	HasValue       bool
}
func (v *PropertyVar) GetScopeType() VariableScope {
	return CLASS_PROPERTY
}

// This is the global variable space
type Global struct {
	Module *ObjModule
	name          string
	IsInitialized bool
	Class         *ClassVar
	ExprData      ExpressionData
}
func (v *Global) GetScopeType() VariableScope {
	return GLOBAL
}

var GlobalVars = make([]Global, 65000)
var GlobalCount = int16(0)

type Local struct {
	Module *ObjModule
	name          string
	depth         int
	isCaptured    bool
	scopeId       int
	IsInitialized bool
	Class         *ClassVar
	Function	  *FunctionVar
	ExprData  ExpressionData

}
func (v *Local) GetScopeType() VariableScope {
	return LOCAL
}

type Upvalue struct {
	Index    int16
	IsLocal  bool
	ExprData  ExpressionData
	Class    *ClassVar
}
func (v *Upvalue) GetScopeType() VariableScope {
	return UPVALUE
}

var namedRegisters = make(map[string]int16)

type register struct {
	isUsed bool
}
func (v *register) GetScopeType() VariableScope {
	return REGISTER
}
var registers = make([]register, 256)
