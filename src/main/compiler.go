package main

import (
	"fmt"
	"strconv"
	"strings"
)

/* ------------------------------------------------------
These are the starting location on the stack
-------------------------------------------------------- */
const (
	GLOBAL_OFFSET    = 0
	LOCAL_OFFSET     = 16000
	STACK_OFFSET     = 1024000
	MAX_MEMORY_SLOTS = 2048000
)

// Current Class. This is a cheap way to know what class we're
// referring to when we come to a . modifier
var CurrentClass *ClassVar

// Keeps track of break and continue instruction locations
type Break struct {
	StartLoc int  // Continue will bump up to here
	CanPatch bool // Flag to indicate if this is waiting to be patched
}

var Breaks = make([]Break, 255)
var BreakPtr int

var StartLoop = make([]int, 255)
var StartPtr int


// Utility
func (c *Compiler) ClearCR() {
	for {
		if !c.Match(TOKEN_CR) {
			break
		}
	}
}

// Keep track of loop types currently in play
const (
	LOOP_WHILE byte = iota
	LOOP_FOR
	LOOP_SCAN
)

var LoopType = make([]byte, 255)
var LoopPtr int

func PushLoop(loopType byte) {
	LoopType[LoopPtr] = loopType
	LoopPtr++
}

func PopLoop() byte {
	LoopPtr--
	return LoopType[LoopPtr]
}

func PeekLoop() byte {
	return LoopType[LoopPtr-1]
}

func PushClass() *ClassVar {
	ClassVarId++
	return &Classes[ClassVarId-1]
}

func PopClass() {
	ClassVarId--
}

var ClassVarId int
var Classes = make([]ClassVar, 128)

type ClassVar struct {
	Id            int
	Class         *ObjClass
	Enclosing     *ClassVar
	Properties    []PropertyVar
	PropertyCount int16
}

func NewClassVar() ClassVar {
	ClassVarId++
	return ClassVar{
		Id:            ClassVarId - 1,
		Enclosing:     nil,
		Properties:    make([]PropertyVar, 65556),
		PropertyCount: 0,
	}
}

// Function structure
type FunctionVar struct {
	Enclosing  *FunctionVar
	paramCount int16
	returnType ValueType
	instr      Instructions

	Locals   []Local
	Upvalues []Upvalue

	Classes []Local

	LocalCount   int16
	UpvalueCount int16
}

func (f *FunctionVar) ConvertToObj() *ObjFunction {
	FunctionId++
	return &ObjFunction{
		Arity:        f.paramCount,
		Code:         f.instr.ToChunk(),
		UpvalueCount: int(f.UpvalueCount),
		FuncType:     TYPE_FUNCTION,
		Id:           FunctionId,
	}
}


type Compiler struct {
	Current *FunctionVar

	Parser     *Parser
	Rules      []ParseRule
	registers  []register
	ScopeDepth int
	DebugMode  bool
}

func Compile(source *string, dbgMode bool) *ObjFunction {

	// First parse the source
	parser := NewParser(source)

	compiler := NewCompiler(&parser)
	compiler.DebugMode = dbgMode

	// Global function store
	RegisterFunctions()

	compiler.Advance()
	for !compiler.Match(TOKEN_EOF) {
		compiler.Evaluate()
	}

	compiler.EmitOp(OP_HALT)
	if dbgMode {
		fmt.Println("=== Instructions ===")
		compiler.CurrentInstructions().Display()
	}

	fn := &ObjFunction{
		Arity:        0,
		Code:         compiler.CurrentInstructions().ToChunk(),
		UpvalueCount: 0,
		FuncType:     TYPE_SCRIPT,
		Id:           0,
	}
	if compiler.Parser.HadError {
		return nil
	}
	return fn
}

/* -------------------------------------------------------
Creates and initializes a app ready for use
 ------------------------------------------------------- */
func NewCompiler(parser *Parser) *Compiler {

	var compiler Compiler

	compiler.Parser = parser
	compiler.Init()
	compiler.LoadRules()
	return &compiler
}

func (c *Compiler) Init() {

	c.Current = &FunctionVar{
		paramCount: 0,
		Enclosing:  nil,
		returnType: VAL_NIL,
		instr:      NewInstructions(),
		Locals:     make([]Local, 65000),
		Upvalues:   make([]Upvalue, 65000),
		LocalCount: 0,
	}

	c.ScopeDepth = 0

}

func (c *Compiler) CurrentInstructions() *Instructions {
	return &c.Current.instr
}

func (c *Compiler) FreeRegister(location int16) {
	registers[location].isUsed = false
}

func (c *Compiler) GetFreeRegister() int16 {
	for i := int16(0); i < 256; i++ {
		if !registers[i].isUsed {
			registers[i] = register{
				isUsed: true,
			}
			return i
		}
	}
	return -1
}

func (c *Compiler) MakeConstant(value Obj) int16 {

	// This way we can re-use the constants table for
	// strings and save some space
	if value.Type() == VAL_STRING {
		var i int16
		for i = 0; i < c.CurrentInstructions().ConstantsCount; i++ {
			// it's a string
			if c.CurrentInstructions().Constants[i].Type() == VAL_STRING {
				// If the strings match
				if c.CurrentInstructions().Constants[i].(ObjString) == value.(ObjString) {
					return i
				}
			}
		}

	}

	constant := c.AddConstant(value)
	if constant > 16000 {
		c.Error("Too many constants")
		return 0
	}
	return constant
}

func (c *Compiler) AddConstant(val Obj) int16 {

	c.CurrentInstructions().Constants[c.CurrentInstructions().ConstantsCount] = val
	c.CurrentInstructions().ConstantsCount++
	return c.CurrentInstructions().ConstantsCount - 1
}

func (c *Compiler) Advance() {
	c.Parser.Prev2 = c.Parser.Previous
	c.Parser.Previous = c.Parser.Current
	c.Parser.Current = c.Parser.NextToken()
}

func (c *Compiler) ResolveGlobal(tok *Token) (int16, *ExpressionData) {
	var i int16
	for i = 0; i < GlobalCount; i++ {
		if tok.ToString() == GlobalVars[i].name {
			return i, &GlobalVars[i].ExprData
		}
	}
	c.ErrorAtCurrent(fmt.Sprintf("Global variable '%s' not found", tok.ToString()))
	return -1, nil
}

func (c *Compiler) ResolveLocal(fn *FunctionVar, name string) (int16, *ExpressionData) {
	// Work our way backwards from the bottom of the locals store so we can
	// identify the variable at lowest scope relative to this one
	for i := fn.LocalCount - 1; i >= 0; i-- {
		if c.IdentifiersEqual(name, fn.Locals[i].name) {
			if fn.Locals[i].depth == -1 {
				c.Error("Cannot read local variable in its own initializer.")
			}
			return i, &fn.Locals[i].ExprData
		}
	}
	return -1, nil
}

func (c *Compiler) AddUpvalue(fn *FunctionVar, index int16, isLocal bool) int16 {

	upvCount := fn.UpvalueCount

	// First check to see if this value has already been tagged as
	// an upvalue. No point in storing it again.
	for i := int16(0); i < upvCount; i++ {
		upVal := fn.Upvalues[i]
		if upVal.Index == index && upVal.IsLocal == isLocal {
			// Found one already used. Just send the index back
			return i
		}
	}

	if upvCount >= 16556 {
		c.Error("Too many closure variables in this function")
	}

	// Create a new closure
	fn.Upvalues[upvCount].IsLocal = isLocal
	fn.Upvalues[upvCount].Index = index
	fn.Upvalues[upvCount].ExprData = fn.Enclosing.Locals[index].ExprData
	fn.UpvalueCount++

	return fn.UpvalueCount - 1

}

func (c *Compiler) AddLocal(name string) int16 {

	if c.Current.LocalCount == 16000 {
		c.Error("Too many local variables in function.")
		return -1
	}

	c.Current.Locals[c.Current.LocalCount].name = name
	c.Current.Locals[c.Current.LocalCount].depth = c.ScopeDepth
	c.Current.Locals[c.Current.LocalCount].isCaptured = false
	c.Current.Locals[c.Current.LocalCount].scopeId = ScopeId

	c.Current.LocalCount++
	return c.Current.LocalCount - 1

}

type Address struct {
	scope       int
	objtype     VarType
	baseAddress int16
	index       int
}

func (c *Compiler) AddressOfVariable() Address {

	addr := Address{
		scope:       c.ScopeDepth,
		objtype:     0,
		baseAddress: 0,
	}

	c.Consume(TOKEN_IDENTIFIER, "Expecting variable name after ")

	addr.objtype = VAR_SCALAR
	tok := c.Parser.Previous

	if c.Check(TOKEN_LEFT_BRACKET) {
		addr.objtype = VAR_ARRAY
	}

	// If it's a local variable, we look for that before globals
	idx, _ := c.ResolveLocal(c.Current, tok.ToString())
	// -1 means it wasn't found
	if idx != -1 {
		idx, _ = c.ResolveGlobal(&tok)
	}

	addr.baseAddress = idx
	return addr
}

func (c *Compiler) ResolveUpvalue(fn *FunctionVar, name string) (int16, *ExpressionData) {

	if fn.Enclosing == nil {
		// There is no function above this one
		return -1, nil
	}

	// Resolve the local variable for the enclosing function
	local, exprValue := c.ResolveLocal(fn.Enclosing, name)

	// We found the value in the enclosing function, so we
	// go ahead and add it to the current function's upvalue store
	if local != -1 {
		fn.Enclosing.Locals[local].isCaptured = true
		uix := c.AddUpvalue(fn, local, true)
		return uix, exprValue
	}

	// On the other hand .. if we still didn't find it in the enclosing function
	// then maybe it's in a function higher up the scope
	upVal, exprValue := c.ResolveUpvalue(fn.Enclosing, name)
	if upVal != -1 {
		// And if we found it, we add the upvalue to this upvalue
		uix := c.AddUpvalue(fn, upVal, false)
		return uix, exprValue
	}
	return -1, nil
}

func (c *Compiler) NamedList(tok Token) {
	// Get the list
	idx, _, varscope := c.ResolveVariable(tok)

	// Get the key
	c.Consume(TOKEN_IDENTIFIER, "Expect key after '$'")
	key := c.Parser.Previous.ToString()

	kIdx := c.MakeConstant(ObjString(key))
	c.EmitInstr(OP_PUSH, kIdx)

	if c.Match(TOKEN_EQUAL) {
		c.Expression()
		if varscope == GLOBAL {
			c.EmitInstr(OP_SET_HGLOBAL, idx)
		} else {
			c.EmitInstr(OP_SET_HLOCAL, idx)
		}
		c.WriteComment(fmt.Sprintf("Array name %s Index %d", tok.ToString(), idx))
	} else {
		if varscope == GLOBAL {
			c.EmitInstr(OP_GET_HGLOBAL, idx)
		} else {
			c.EmitInstr(OP_GET_HLOCAL, idx)
		}
		c.WriteComment(fmt.Sprintf("List name '%s' Index '%s'", tok.ToString(), key))
	}

}

func (c *Compiler) NamedArray(tok Token) {

	// This gets the address of the entire array
	idx, expData, varscope := c.ResolveVariable(tok)
	// Push the array on to the stack
	if varscope == GLOBAL {
		c.EmitInstr(OP_GET_GLOBAL, idx)
	} else {
		c.EmitInstr(OP_GET_LOCAL, idx)
	}
	c.WriteComment(fmt.Sprintf("Array name %s Location %d of type %s", tok.ToString(), idx,ValueTypeLabel[expData.Value]))
	dims := 0
	c.Consume(TOKEN_LEFT_BRACKET,"Expect '[' after array name")
	for {
		c.Expression()
		dims++
		if !c.Match(TOKEN_COMMA) {
			break
		}
	}
	c.Consume(TOKEN_RIGHT_BRACKET,"']' must follow array index reference")

	if c.Match(TOKEN_EQUAL) {
		c.Expression()
		PopExpressionValue()
		if varscope == GLOBAL {
			c.EmitInstr(OP_SET_AGLOBAL, idx)
		} else {
			c.EmitInstr(OP_SET_ALOCAL, idx)
		}
		c.WriteComment(fmt.Sprintf("Array name %s Index %d", tok.ToString(), idx))
	} else {
		c.EmitInstr(OP_AINDEX,int16(dims))
		c.WriteComment(fmt.Sprintf("Getting array index with %d dimensions",dims))
		PushExpressionValue(expData)
	}

}

func (c *Compiler) ResolveVariable(tok Token) (int16, ExpressionData, VariableScope) {

	var idx int16
	var expData *ExpressionData
	var vScope VariableScope
	var ok bool

	if idx, expData = c.ResolveLocal(c.Current, tok.ToString()); idx != -1 {
		vScope = LOCAL
	} else if idx, expData = c.ResolveGlobal(&tok); idx != -1 {
		vScope = GLOBAL
	} else if idx, expData = c.ResolveUpvalue(c.Current, tok.ToString()); idx != -1 {
		vScope = UPVALUE
	} else if idx, ok = namedRegisters[tok.ToString()]; ok {
		vScope = REGISTER
	} else {
		c.Error("Variable not found")
	}
	return idx, *expData, vScope
}

func (c *Compiler) NamedVariable(canAssign bool) {

	tok := c.Parser.Previous
	// Above all, check to see if this name is a built-in function
	nativeFunction := ResolveNativeFunction(tok.ToString())
	if nativeFunction != nil {
		c.CallNative(nativeFunction)
		return
	}

	objType := VAR_SCALAR

	var getOp byte
	var setOp byte
	var valType ValueType

	isLocal := false
	isGlobal := false
	isUpvalue := false

	isHasOperand := false

	//If this is a list
	if c.Match(TOKEN_DOLLAR) {
		c.NamedList(tok)
		return
	}

	// If this is an array .
	if c.Check(TOKEN_LEFT_BRACKET) {
		c.NamedArray(tok)
		return
	}

	// If it's a local variable, we look for that before globals
	idx, _ := c.ResolveLocal(c.Current, tok.ToString())
	// -1 means it wasn't found
	if idx != -1 {
		isLocal = true
		// If this is an expression of an array element
		switch idx {
		case 0:
			getOp = OP_GET_LOCAL_0
		case 1:
			getOp = OP_GET_LOCAL_1
		case 2:
			getOp = OP_GET_LOCAL_2
		case 3:
			getOp = OP_GET_LOCAL_3
		case 4:
			getOp = OP_GET_LOCAL_4
		case 5:
			getOp = OP_GET_LOCAL_5
		default:
			getOp = OP_GET_LOCAL
			isHasOperand = true
		}
		setOp = OP_SET_LOCAL

	} else if idx, _ = c.ResolveUpvalue(c.Current, tok.ToString()); idx != -1 {
		isUpvalue = true
		getOp = OP_GET_UPVALUE
		setOp = OP_SET_UPVALUE
		isHasOperand = true

	} else {
		// Is it in a register?
		ridx, ok := namedRegisters[tok.ToString()]
		if ok {
			idx = ridx
			getOp = OP_GET_REGISTER
			setOp = OP_SET_REGISTER
			valType = VAL_INTEGER
			isHasOperand = true
		} else {
			isGlobal = true
			idx, _ = c.ResolveGlobal(&tok)
			setOp = OP_SET_GLOBAL
			if idx != -1 {
				switch idx {
				case 0:
					getOp = OP_GET_GLOBAL_0
				case 1:
					getOp = OP_GET_GLOBAL_1
				case 2:
					getOp = OP_GET_GLOBAL_2
				case 3:
					getOp = OP_GET_GLOBAL_3
				case 4:
					getOp = OP_GET_GLOBAL_4
				case 5:
					getOp = OP_GET_GLOBAL_5
				default:
					getOp = OP_GET_GLOBAL
					isHasOperand = true
				}
			}
		}
	}

	if canAssign && c.Match(TOKEN_EQUAL) {

		c.Expression()
		data := PopExpressionValue()

		valType = data.Value
		objType = data.ObjType

		if isGlobal {

			GlobalVars[idx].IsInitialized = true
			//GlobalVars[idx].datatype = valType
			//GlobalVars[idx].objtype = objType

			if objType == VAR_CLASS {

				GlobalVars[idx].Class = CurrentClass
			}

			if GlobalVars[idx].ExprData.Value != valType || GlobalVars[idx].ExprData.ObjType != objType {
				gVar := ValueTypeLabel[GlobalVars[idx].ExprData.Value]
				gObj := VarTypeLabel[ GlobalVars[idx].ExprData.ObjType]
				errStr := fmt.Sprintf("Variable %s is a %s of type %s: cannot assign a %s of type %s",
					tok.ToString(),gObj,gVar,VarTypeLabel[objType],ValueTypeLabel[valType])
				c.Error(errStr)
			}

		} else if isLocal {
			if c.Current.Locals[idx].IsInitialized && valType != c.Current.Locals[idx].ExprData.Value {
				c.Error("Cannot assign incompatible variable")
			}

			c.Current.Locals[idx].IsInitialized = true
			c.Current.Locals[idx].ExprData.Value = valType
			c.Current.Locals[idx].ExprData.ObjType = objType

			if objType == VAR_CLASS {
				c.Current.Locals[idx].Class = CurrentClass
			}

		} else if isUpvalue {
			c.Current.Upvalues[idx].ExprData.Value = valType
			if objType == VAR_CLASS {
				c.Current.Upvalues[idx].Class = CurrentClass
			}
		}
		c.EmitInstr(setOp,idx)
		c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[setOp], tok.ToString(), idx, valType))
	} else {
		if isHasOperand {
			c.EmitInstr(getOp, idx)
		} else {
			c.EmitOp(getOp)
		}
		if isGlobal {
			valType = GlobalVars[idx].ExprData.Value
			objType = GlobalVars[idx].ExprData.ObjType
			if objType == VAR_CLASS {
				CurrentClass = GlobalVars[idx].Class
			}

		} else if isLocal {
			valType = c.Current.Locals[idx].ExprData.Value
			objType = c.Current.Locals[idx].ExprData.ObjType
			if objType == VAR_CLASS {
				CurrentClass = c.Current.Locals[idx].Class
			}
		} else if isUpvalue {
			valType = c.Current.Upvalues[idx].ExprData.Value
			if objType == VAR_CLASS {
				CurrentClass = c.Current.Upvalues[idx].Class
			}
		}
		c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[getOp], tok.ToString(), idx, valType))
	}
	PushExpressionValue(ExpressionData{Value: valType, ObjType: objType})
}

func (c *Compiler) IdentifierConstant() int16 {
	return c.MakeConstant(ObjString(string(c.Parser.Previous.Value)))
}

func (c *Compiler) DefineProperty() {
	c.Consume(TOKEN_IDENTIFIER, "Expect method name")
	// Store the token here
	//tok := c.Parser.Previous
	var valType ValueType

	c.Consume(TOKEN_COLON, "Expect ':' with type aftermethod")
	valType = c.GetDataType().Value
	if valType == VAL_NIL {
		c.ErrorAtCurrent("Invalid data type")
	}

	//var index int16
	//for

}

func (c *Compiler) DefineParameter() {
	c.Consume(TOKEN_IDENTIFIER, "Expect parameter name")
	// Store the token here
	tok := c.Parser.Previous
	var valType ValueType

	c.Consume(TOKEN_COLON, "Expect ':' with type after parameter")
	valType = c.GetDataType().Value
	if valType == VAL_NIL {
		c.ErrorAtCurrent("Invalid data type")
	}
	var index int16

	for i := c.Current.LocalCount - 1; i >= 0; i-- {
		if c.Current.Locals[i].depth != -1 &&
			c.Current.Locals[i].depth < c.ScopeDepth {
			break
		}

		if c.IdentifiersEqual(tok.ToString(), c.Current.Locals[i].name) &&
			c.Current.Locals[i].scopeId == ScopeId {
			c.Error(fmt.Sprintf("Variable with the name %s already declared in this scope.", tok.ToString()))
		}
	}

	index = c.AddLocal(tok.ToString())
	c.Current.Locals[index].ExprData.Value = valType

}

func (c *Compiler) CheckArrayType(valueType ValueType) *ExpressionData {
	var expd ExpressionData
	if c.Match(TOKEN_LEFT_BRACKET) {
		dims:=0
		for {
			dims++
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.Consume(TOKEN_RIGHT_BRACKET,"Expect ']' after array dimension declaration")
		expd.ObjType = VAR_ARRAY
		expd.Value = valueType
		expd.Dimensions = dims
	} else 	{
		expd.ObjType = VAR_SCALAR
		expd.Value = valueType
	}
	return &expd
}

// This function is called to make sense of variable declarations where we
// declare the type before the variable name
func (c *Compiler) GetDataType() ExpressionData {
	expd := new(ExpressionData)

	switch {

	case c.Check(TOKEN_TYPE_INTEGER):
		c.Advance()
		expd = c.CheckArrayType(VAL_INTEGER)
	case c.Check(TOKEN_TYPE_BOOL):
		c.Advance()
		expd = c.CheckArrayType(VAL_BOOL)
	case c.Check(TOKEN_TYPE_FLOAT):
		c.Advance()
		expd = c.CheckArrayType(VAL_FLOAT)
	case c.Check(TOKEN_TYPE_STRING):
		c.Advance()
		expd = c.CheckArrayType(VAL_STRING)
	case c.Check(TOKEN_FUNC):
		c.Advance()
		expd.ObjType = VAR_FUNCTION
		expd.Value = VAL_FUNCTION
	case c.Check(TOKEN_CLASS):
		c.Advance()
		expd.ObjType = VAR_CLASS
		expd.Value = VAL_CLASS
	case c.Check(TOKEN_LIST_TYPE):
		c.Advance()
		expd.Value = VAL_STRING
		expd.ObjType = VAR_HASH
	case c.Check(TOKEN_ENUM):
		c.Advance()
		expd.Value = VAL_ENUM
		expd.ObjType = VAR_ENUM
	case c.Check(TOKEN_IDENTIFIER):
		// This could be a user defined type such as a class
		//tok := c.Parser.Current
		//idx, expData, varScope := c.ResolveVariable(tok)
		//if idx != -1 {
		//	tt := ValueType(100)
		//}

	default:
		{
			expd.ObjType = VAR_UNKNOWN
			expd.Value = VAL_NIL
		}
	}

	return *expd
}

func (c *Compiler) _array(canAssign bool) {
	//array()
}

func (c *Compiler) AddGlobal(varName string) int16 {
	for i := int16(0); i < GlobalCount; i++ {
		if GlobalVars[i].name == varName {
			c.Error(fmt.Sprintf("%s has already been defined"))
		}
	}
	GlobalVars[GlobalCount].name = varName
	GlobalCount++
	return GlobalCount - 1
}

func (c *Compiler) MatchTypes( data1 ExpressionData, data2 ExpressionData) bool {
	return !(data1.Value !=data2.Value || data1.ObjType != data2.ObjType ||	data1.Dimensions != data2.Dimensions) &&
	data2.Value != VAL_NIL
}

func (c *Compiler) DeclareGlobalVariable(varName string) {

	index := c.AddGlobal(varName)
	if c.Match(TOKEN_EQUAL) {
		// This is the value we're going to assign
		c.Expression()
		GlobalVars[index].ExprData = PopExpressionValue()

		c.EmitInstr(OP_SET_GLOBAL, index)
		c.WriteComment(fmt.Sprintf("Setting global variable %s at location %d",varName,index))
	} else {
		GlobalVars[index].ExprData = c.GetDataType()
	}

}

func (c *Compiler) DeclareLocalVariable(varName string) {

	// Check if this variable was already declared in the current scope
	for i := c.Current.LocalCount - 1; i >= 0; i-- {
		if c.Current.Locals[i].depth != -1 &&
			c.Current.Locals[i].depth < c.ScopeDepth {
			break
		}

		if c.IdentifiersEqual(varName, c.Current.Locals[i].name) &&
			c.Current.Locals[i].scopeId == ScopeId {
			c.Error(fmt.Sprintf("Variable with the name %s already declared in this scope.", varName))
		}
	}
	//opcode = OP_SET_LOCAL
	index := c.AddLocal(varName)
	c.Current.Locals[index].name = varName

	if c.Match(TOKEN_EQUAL) {
		// This is the value we're going to assign
		c.Expression()
		c.Current.Locals[index].ExprData = PopExpressionValue()

		c.EmitInstr(OP_SET_LOCAL, index)
	} else {
		c.Current.Locals[index].ExprData = c.GetDataType()
	}
}

func (c *Compiler) DeclareVariable() {

	c.Consume(TOKEN_IDENTIFIER, "Expect variable name")
	tok := c.Parser.Previous

	// Error if this variable collides with an existing native function name
	if ResolveNativeFunction(tok.ToString()) != nil {
		c.Error(fmt.Sprintf("'%s' is a reserved name", tok.ToString()))
	}

	//var scope VariableScope
	if c.ScopeDepth == 0 {
		//scope = GLOBAL
		c.DeclareGlobalVariable(tok.ToString())
	} else {
		//scope = LOCAL
		c.DeclareLocalVariable(tok.ToString())
	}

}

func (c *Compiler) IdentifiersEqual(a string, b string) bool {
	return a == b
}

// Checks to see if the current token type matches the given
// token type without actually consuming or advancing the token
// position with the scanner
func (c *Compiler) Check(t TokenType) bool {
	return c.Parser.Current.Type == t
}

// If the current token type doesn't match with what we
// expected, then return false, but if it does, we advance
// the pointer and return true. Good for cases where we don't
// need to consume the token .. we already know what it is by
// returning true
func (c *Compiler) Match(t TokenType) bool {
	if !c.Check(t) {
		return false
	} else {
		c.Advance()
		return true
	}
}

// Here we're sure we know what the next token needs to be.
// Anything else and we have to consider it an error
func (c *Compiler) Consume(t TokenType, message string) {
	if c.Parser.Current.Type == t {
		c.Advance()
		return
	}
	// Stop the presses ..! This is a problem
	c.ErrorAtCurrent(fmt.Sprintf("%s: Looking for Token %d but have %d", message, t, c.Parser.Current.Type))
}

/*
If the scanner hands us an error token, we need to actually tell the user.
That happens here:
*/
func (c *Compiler) ErrorAtCurrent(message string) {
	fmt.Printf("Line %d: %s\n", c.Parser.Current.Line, message)
}

/*
In the end - all the real error management happens here
*/
func (c *Compiler) ErrorAt(token *Token, message string) {
	// No point in tracking further error tracking if we're
	// already in error mode
	if c.Parser.PanicMode {
		return
	}
	// This tells the app that we're in error mode now,
	// but we keep evaluating code without actually generating byte code
	c.Parser.PanicMode = true

	fmt.Printf("[line %d] Error", token.Line+1)
	switch token.Type {
	case TOKEN_EOF:
		fmt.Printf(" at end")
	case TOKEN_ERROR:
		fmt.Printf(" ERROR")
	default:
		fmt.Printf(" at '%s'", token.Value)
	}
	fmt.Printf(": %s\n", message)
	c.Parser.HadError = true
}

/*
We pull the location out of the current token in order to tell the user where the error occurred
and forward it to errorAt(). More often, we’ll report an error at the location of the token we just consumed,
so we give the shorter name the actual error here
*/
func (c *Compiler) Error(message string) {
	c.ErrorAt(&c.Parser.Previous, message)
}

func (c *Compiler) EmitOp(opcode byte) {
	c.CurrentInstructions().WriteSimpleInstruction(opcode, c.Parser.Previous.Line)
}

func (c *Compiler) EmitOperand(val int16) {
	c.CurrentInstructions().AddOperand(val)
}

func (c *Compiler) EmitOperand32(val int32) {
	c.CurrentInstructions().AddOperand32(val)
}

func (c *Compiler) EmitOperandByte(val byte) {
	c.CurrentInstructions().AddByteOperand(val)
}

func (c *Compiler) EmitPushInteger(val int16) {
	switch val {
	case 0:
		c.EmitOp(OP_PUSH_0)
	case 1:
		c.EmitOp(OP_PUSH_1)
	default:
		c.EmitInstr(OP_PUSH, val)
	}
}

func (c *Compiler) EmitInstr(opcode byte, operand int16) {
	c.CurrentInstructions().WriteInstruction(opcode, operand, c.Parser.Previous.Line)
}

func (c *Compiler) EmitSingleByteInstr(opcode byte, operand byte) {
	c.CurrentInstructions().WriteSingleByteInstruction(opcode, operand, c.Parser.Previous.Line)
}

func (c *Compiler) EmitReturn() {
	c.EmitOp(OP_NIL)
	c.EmitOp(OP_RETURN)
}

func (c *Compiler) ReturnStatement() {
	if c.Match(TOKEN_CR) {
		c.EmitReturn()
	} else {
		c.Expression()
		c.EmitOp(OP_RETURN)
	}
}

func (c *Compiler) ParsePrecedence(precedence Precedence) {
	// This loads the prefix rule which either contains a value such as
	// a variable or literal or a prefix that affects the next value
	c.Advance()
	prefixRule := c.GetRule(c.Parser.Previous.Type).prefix

	// This is an error in that an expression needs to at least begin
	// with a prefix rule
	if prefixRule == nil {
		fmt.Printf("RULE: %d has no prefix \n", c.Parser.Previous.Type)
		panic("Expect expression")
		return
	}

	canAssign := precedence <= PREC_ASSIGNMENT
	prefixRule(canAssign)

	for precedence <= c.GetRule(c.Parser.Current.Type).Prec {
		c.Advance()

		infixRule := c.GetRule(c.Parser.Previous.Type).infix
		if infixRule != nil {
			infixRule(canAssign)
		}

		postfixRule := c.GetRule(c.Parser.Previous.Type).postfix
		if postfixRule != nil {
			postfixRule(canAssign)
		}
	}

	if canAssign && c.Match(TOKEN_EQUAL) {
		str := fmt.Sprintf("Invalid assignment target.")
		c.ErrorAtCurrent(str)
	}
}

func (c *Compiler) Grouping(canAssign bool) {
	c.Expression()
	c.Consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression")
}
func (c *Compiler) Call(canAssign bool) {
	argumentCount := c.GetArguments()

	switch argumentCount {
	case 0:
		c.EmitOp(OP_CALL_0)
	case 1:
		c.EmitOp(OP_CALL_1)
	case 2:
		c.EmitOp(OP_CALL_2)
	case 3:
		c.EmitOp(OP_CALL_0)
	default:
		c.EmitInstr(OP_CALL, argumentCount)
	}

	c.WriteComment(fmt.Sprintf("Function call with %d arguments", argumentCount))
}

func (c *Compiler) CallMethod(constantIndex int16) {

	argumentCount := c.GetArguments() + 1
	c.EmitInstr(OP_CALL_METHOD, constantIndex)
	c.EmitOperand(argumentCount)
	c.EmitOp(OP_POP)
}

func (c *Compiler) GetArguments() int16 {
	argCount := int16(0)
	if !c.Check(TOKEN_RIGHT_PAREN) {
		for {
			c.Expression()
			if argCount == 255 {
				c.Error("Cannot have more than 255 arguments.")
			}
			argCount++

			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
	}

	c.Consume(TOKEN_RIGHT_PAREN, "Expect ')' after arguments.")
	return argCount
}

func (c *Compiler) Dollar(canAssign bool) {
	c.Consume(TOKEN_IDENTIFIER,"Expect key name after '$'")
	keyVal := c.Parser.Previous.ToString()
	idx := c.MakeConstant(ObjString(keyVal))
	c.EmitInstr(OP_HKEY,idx)
	c.WriteComment(fmt.Sprintf("Getting list value key %s",keyVal))
}

func (c *Compiler) New(canAssign bool) {

	var valType ValueType

	// This is an array
	switch {
	case c.Match(TOKEN_TYPE_INTEGER):
		valType = VAL_INTEGER
	case c.Match(TOKEN_TYPE_FLOAT):
		valType = VAL_FLOAT
	case c.Match(TOKEN_LIST_TYPE):
		// Defer to the new list expression
		c.NewList(canAssign)
		valType = VAL_LIST
	default:
		c.Expression()
		c.EmitOp(OP_OBJ_INSTANCE)
		PushExpressionValue(ExpressionData{
			Value: VAL_OBJECT,
			ObjType: VAR_OBJECT,
		})
		return
	}
	// Loop in order to handle multi-dimensional arrays
	dims := int16(0)
	c.Consume(TOKEN_LEFT_BRACKET, "Expect '[' after new array declaration")
	for {
		c.Expression()
		dims++
		if !c.Match(TOKEN_COMMA) {
			// Nor more dimensions
			break
		}
	}
	c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after new array declaration")
	c.EmitInstr(OP_PUSH,int16(valType))
	c.EmitInstr(OP_MAKE_ARRAY, int16(dims))
	PushExpressionValue(ExpressionData{
		Value: valType,
		ObjType: VAR_ARRAY,
	})

}

func (c *Compiler) Postary(canAssign bool) {
	operatorType := c.Parser.Previous.Type
	// Emit the operator instruction.
	switch operatorType {
	case TOKEN_PLUS_PLUS:

	case TOKEN_MINUS_MINUS:
		c.EmitOp(OP_DECREMENT)
	}
}

func (c *Compiler) Unary(canAssign bool) {
	operatorType := c.Parser.Previous.Type

	// Compile the operand.
	c.ParsePrecedence(PREC_UNARY)

	valtype := c.GetDataType().Value
	// Emit the operator instruction.
	switch operatorType {
	case TOKEN_BANG:
		c.EmitOp(OP_NOT)
	case TOKEN_MINUS:
		switch valtype {
		case VAL_INTEGER:
			c.EmitOp(OP_INEGATE)
		case VAL_FLOAT:
			c.EmitOp(OP_FNEGATE)
		}

	case TOKEN_PLUS_PLUS:
		c.EmitOp(OP_PREINCREMENT)
	case TOKEN_MINUS_MINUS:
		c.EmitOp(OP_PREDECREMENT)
	}
}

func (c *Compiler) Binary(canAssign bool) {
	// This the operator that made us call this
	// function in the first place
	operatorType := c.Parser.Previous.Type

	// Compile the right operand
	rule := c.GetRule(operatorType)
	rprec := rule.Prec + 1
	c.ParsePrecedence(rprec)

	data := PopExpressionValue()

	switch operatorType {
	case TOKEN_BANG_EQUAL:
		c.EmitOp(OP_NOT_EQUAL)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_EQUAL_EQUAL:
		c.EmitOp(OP_EQUAL)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_GREATER:
		c.EmitOp(OP_GREATER)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_GREATER_EQUAL:
		c.EmitOp(OP_GREATER_EQUAL)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_LESS:
		c.EmitOp(OP_LESS)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_LESS_EQUAL:
		c.EmitOp(OP_LESS_EQUAL)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: data.ObjType})
	case TOKEN_PLUS:
		switch data.Value {
		case VAL_INTEGER:
			c.EmitOp(OP_IADD)
			PushExpressionValue(data)
		case VAL_FLOAT:
			c.EmitOp(OP_FADD)
			PushExpressionValue(data)
		case VAL_STRING:
			c.EmitOp(OP_SADD)
			PushExpressionValue(data)
		}
	case TOKEN_MINUS:
		if data.Value == VAL_INTEGER {
			c.EmitOp(OP_ISUBTRACT)
			PushExpressionValue(data)
		} else if data.Value == VAL_FLOAT {
			c.EmitOp(OP_FSUBTRACT)
			PushExpressionValue(data)
		}
	case TOKEN_STAR:
		if data.Value == VAL_INTEGER {
			c.EmitOp(OP_IMULTIPLY)
			PushExpressionValue(data)
		} else if data.Value == VAL_FLOAT {
			c.EmitOp(OP_FMULTIPLY)
			PushExpressionValue(data)
		} else {
			//c.EmitOp(OP_IMULTIPLY)
			//PushExpressionValue(data)
			//fmt.Println("Error finding property type")
			c.Error(fmt.Sprintf("Can't multiply! Type: %v", data.Value))
		}
	case TOKEN_SLASH:
		if data.Value == VAL_INTEGER {
			c.EmitOp(OP_IDIVIDE)
			PushExpressionValue(data)
		} else if data.Value == VAL_FLOAT {
			c.EmitOp(OP_FDIVIDE)
			PushExpressionValue(data)
		}
	case TOKEN_PLUS_PLUS:
		c.EmitOp(OP_INCREMENT)
		PushExpressionValue(data)
	case TOKEN_HAT:
		if data.Value == VAL_INTEGER {
			c.EmitOp(OP_IEXP)
			PushExpressionValue(data)
		} else {
			c.Error("Exponents can only be defined on integers")
		}
	case TOKEN_TO:
		if data.Value == VAL_INTEGER {
			c.EmitOp(OP_IRANGE)
			PushExpressionValue(data)
		} else {
			c.Error("Exponents can only be defined on integers")
		}
	default:
		return
	}
}

func (c *Compiler) FindPropertyType(class *ClassVar, propertyName string) ExpressionData {
	for i := int16(0); i < class.PropertyCount; i++ {
		if class.Properties[i].Name == propertyName {
			return class.Properties[i].ExprData
		}
	}
	return ExpressionData{VAL_NIL, VAR_UNKNOWN, 0}
}

func (c *Compiler) NewList(canAssign bool) {
	// If this is on the right side
	if c.Match(TOKEN_LEFT_BRACKET) {
		keyType := c.GetDataType()
		c.Consume(TOKEN_COMMA,"Expect ',' between data types in list allocation")
		valType := c.GetDataType()
		c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after list allocation")

		c.EmitPushInteger(int16(VarType(valType.ObjType)))
		c.EmitPushInteger(int16(ValueType(valType.Value)))
		c.EmitPushInteger(int16(ValueType(keyType.Value)))

		c.EmitOp(OP_MAKE_LIST)
		PushExpressionValue(ExpressionData{keyType.Value, valType.ObjType ,1})
	}
	// Left side, do nothing

}

func (c *Compiler) List(canAssign bool) {
	keys := int16(0)
	var keyType ValueType
	var dType ValueType

	for {
		// Key
		c.Expression()
		expVal := PopExpressionValue().Value
		if keys == 0 {
			keyType = expVal
		}
		if keyType != expVal {
			c.Error(fmt.Sprintf("Key %d is incompatible with first key type %s",
				keys, ValueTypeLabel[dType]))
		}

		c.Consume(TOKEN_COLON, "Expect ':' after key definition")
		// Value
		c.Expression()
		expVal = PopExpressionValue().Value
		if keys == 0 {
			dType = expVal
		}
		if dType != expVal {
			c.Error(fmt.Sprintf("Value %d is incompatible with first value type %s",
				keys, ValueTypeLabel[dType]))
		}

		keys++
		if !c.Match(TOKEN_COMMA) {
			break
		}
	}
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' after list definition")
	c.EmitInstr(OP_PUSH, keys)
	c.EmitSingleByteInstr(OP_LIST, byte(keyType))
	PushExpressionValue(ExpressionData{
		Value:   VAL_LIST,
		ObjType: VAR_HASH,
	})

	c.WriteComment(fmt.Sprintf("List with %d elements type %s=%s", keys, ValueTypeLabel[keyType], ValueTypeLabel[dType]))
}

func (c *Compiler) Enum(canAssign bool) {
	elements := uint8(0)
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' after 'enum'")

	for {
		c.Consume(TOKEN_IDENTIFIER, "Expect enum element")
		name := c.Parser.Previous.ToString()
		idx := c.MakeConstant(ObjString(name))
		c.EmitInstr(OP_PUSH, idx)
		elements++
		if elements > 255 {
			c.Error("Enums cannot have more than 255 elements")
		}
		if !c.Match(TOKEN_COMMA) {
			break
		}

	}
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' to close enum definition")
	c.EmitInstr(OP_ENUM, int16(elements))
	PushExpressionValue(ExpressionData{
		Value:   VAL_ENUM,
		ObjType: VAR_ENUM,
	})
}

func (c *Compiler) Array(canAssign bool) {
	// Find how many items are in this array
	elements := int16(0)
	var dType ValueType

	dimCount := int16(0)

	// Pushing a multidimensional array
	if c.Match(TOKEN_LEFT_BRACKET) {
		for {
			dimCount++
			c.Expression()

			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.Consume(TOKEN_RIGHT_BRACKET, "Expect right bracket after dimensional sizing")
	} else {
		dimCount++
	}
	c.EmitPushInteger(dimCount)
	c.WriteComment(fmt.Sprintf("Dimensions of array type: %d",dimCount))

	for {

		c.Expression()
		valType := PopExpressionValue().Value
		if elements == 0 {
			dType = valType
		}
		if dType != valType {
			c.Error(fmt.Sprintf("Element %d is incompatible with first element type %d",
				elements, dType))
		}
		elements++
		if !c.Match(TOKEN_COMMA) {
			break
		}
	}

	c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after array definition")

	c.EmitInstr(OP_PUSH, elements)
	c.WriteComment(fmt.Sprintf("Push number of elements %d ",elements))

	c.EmitInstr(OP_ARRAY, int16(dType))
	c.WriteComment(fmt.Sprintf("Array of type %s ",ValueTypeLabel[dType]))

	PushExpressionValue(ExpressionData{
		Value:   dType,
		ObjType: VAR_ARRAY,
		Dimensions: int(dimCount),
	})

}
func (c *Compiler) Index(canAssign bool) {
	c.Expression()
	expData := PopExpressionValue()
	expData.ObjType = VAR_SCALAR
	expData.Dimensions = 0
	PushExpressionValue(expData)
	c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after index reference")
	c.EmitInstr(OP_AINDEX, int16(1))
}

// When evauating a compound object, one optimizaion is to pass the
// reference to the parent objects witout having to emit an OP code
// specifically to load it. Instead, the best thing is to keep a reference
// to it handy until we arrive at the final propery/method/variable. At that point,
// we'll simply "invoke" that variable along with the parent references all in one
// VM cycle. The parent references then will get "shlved" here for future use
type ShelvedRef struct {
	ReferenceId int16
	*VariableScope
	*ExpressionData
}
var shelved = make([]ShelvedRef,64)
func (c *Compiler) ShelveReference() {

}

func (c *Compiler) CompoundVariable(tok *Token) *ExpressionData {

	name := tok.ToString()

	if name == "this" {
		c.EmitOp(OP_GET_LOCAL_0)
		return &ExpressionData{
			Value:   VAL_OBJECT,
			ObjType: VAR_OBJECT,
		}
	}

	idx, expData := c.ResolveLocal(c.Current, name)

	// It's a local
	if idx != -1 {
		switch expData.ObjType {
		// Is it a class?
		case VAR_OBJECT:
			// Then treat this as a Class
			c.EmitInstr(OP_GET_LOCAL, idx)
			return &ExpressionData{
				Value:   VAL_OBJECT,
				ObjType: VAR_OBJECT,
			}
		case VAR_ENUM:
			c.EmitInstr(OP_GET_LOCAL, idx)
			return &ExpressionData{
				Value:   VAL_ENUM,
				ObjType: VAR_ENUM,
			}
		}

	}

	// It's a global
	idx, expData = c.ResolveGlobal(tok)
	if idx != -1 {
		switch expData.ObjType {
		case VAR_OBJECT:
			// Treat this as a Class
			c.EmitInstr(OP_GET_GLOBAL, idx)
			return &ExpressionData{
				Value:   VAL_OBJECT,
				ObjType: VAR_OBJECT,
			}
		case VAR_ENUM:
			c.EmitInstr(OP_GET_GLOBAL, idx)
			return &ExpressionData{
				Value:   VAL_ENUM,
				ObjType: VAR_ENUM,
			}
		}
		fmt.Printf("%s Class type: %s\n", name, VarTypeLabel[expData.ObjType])
	}
	fmt.Printf("Compound variable '%s' not found\n", name)
	return nil
}

func (c *Compiler) Variable(canAssign bool) {

	tok := &c.Parser.Previous // Variable token

	var expData *ExpressionData // Value and type data

	var foundCompoundObject bool

	// In the first pass, we check to see if it's a compound variable
	for c.Check(TOKEN_DOT) {
		foundCompoundObject = true
		// It is and so now we check to see what kind it is
		expData = c.CompoundVariable(tok)
		// Now we get the property
		c.Consume(TOKEN_DOT, "Expect '.' after object name")
		c.Consume(TOKEN_IDENTIFIER, "Expect name after '.'")

		tok = &c.Parser.Previous
		idx := c.MakeConstant(ObjString(tok.ToString()))

		switch expData.ObjType {
		case VAR_OBJECT:

			if c.Match(TOKEN_LEFT_PAREN) {
				// This is a method
				args := c.GetArguments()
				c.EmitInstr(OP_CALL_METHOD, idx)
				c.EmitOperand(args)
			} else {

				// It's a class, so let's manage that
				if c.Match(TOKEN_EQUAL) {
					c.Expression()
					c.EmitInstr(OP_SET_PROPERTY, idx)
				} else {
					c.EmitInstr(OP_GET_PROPERTY, idx)
				}
			}
		case VAR_ENUM:
			c.EmitInstr(OP_ENUM_TAG, idx)
			PushExpressionValue(ExpressionData{VAL_ENUM,VAR_ENUM,1})
		default:
			// Uh oh ..
			c.Error(fmt.Sprintf("Compound variable %s of type %s should not have a dot after it", tok.ToString(), VarTypeLabel[expData.ObjType]))
		}
	}

	if !foundCompoundObject {
		c.NamedVariable(canAssign)
	}

	// Next we check to see if this variable has further indications such as
	// array of hashes or hash of arrays, etc.
	if c.Check(TOKEN_LEFT_BRACKET) || c.Check(TOKEN_DOLLAR) || c.Check(TOKEN_DOT) {
		c.Expression()
	}

}

func (c *Compiler) String(canAssign bool) {
	value := c.Parser.Previous.ToString()
	// Remove the quotes
	idx := c.MakeConstant(ObjString(value[1:(len(value) - 1)]))
	c.EmitInstr(OP_SCONST, idx)

	if len(value) > 10 {
		value = value[0:10]
	}

	c.WriteComment(fmt.Sprintf("Value %s at constant index %d", value, idx))
	PushExpressionValue(ExpressionData{Value: VAL_STRING, ObjType: VAR_SCALAR})
}
func (c *Compiler) Integer(canAssign bool) {
	value, _ := strconv.ParseInt(string(c.Parser.Previous.Value), 10, 64)
	idx := c.MakeConstant(ObjInteger(value))
	c.EmitInstr(OP_ICONST, idx)
	c.WriteComment(fmt.Sprintf("Value %d at constant index %d", value, idx))
	PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: VAR_SCALAR})
}
func (c *Compiler) Float(canAssign bool) {
	value, _ := strconv.ParseFloat(string(c.Parser.Previous.Value), 64)
	idx := c.MakeConstant(ObjFloat(value))
	c.EmitInstr(OP_FCONST, idx)
	PushExpressionValue(ExpressionData{Value: VAL_FLOAT, ObjType: VAR_SCALAR})

}
func (c *Compiler) Browse(canAssign bool) {}
func (c *Compiler) and_(canAssign bool) {
	endJump := c.EmitJump(OP_JUMP_IF_FALSE)

	c.EmitOp(OP_POP)
	c.ParsePrecedence(PREC_AND)

	c.PatchJump(endJump)
}

func (c *Compiler) or_(canAssign bool) {
	elseJump := c.EmitJump(OP_JUMP_IF_FALSE)
	endJump := c.EmitJump(OP_JUMP)

	c.PatchJump(elseJump)
	c.EmitOp(OP_POP)

	c.ParsePrecedence(PREC_OR)
	c.PatchJump(endJump)
}

func (c *Compiler) Literal(canAssign bool) {
	fmt.Println("LITERAL")
}
func (c *Compiler) Boolean(canAssign bool) {
	value := strings.ToUpper(c.Parser.Previous.ToString())
	if value == "TRUE" {
		c.EmitOp(OP_TRUE)
	} else {
		c.EmitOp(OP_FALSE)
	}
	PushExpressionValue(ExpressionData{Value: VAL_BOOL, ObjType: VAR_SCALAR})
}
func (c *Compiler) SqlSelect(canAssign bool) {}

func (c *Compiler) Expression() {
	c.ParsePrecedence(PREC_ASSIGNMENT)
}

func (c *Compiler) ExpressionStatement() {
	c.Expression()
	// After the expression gets evaluated, we display it on the output device
	// That's what makes this a "statement" rather than an expression only
	//c.Consume(TOKEN_CR, "Expect 'CR' after expression.")
	//c.EmitOp(OP_POP)
	//c.WriteComment("Pop After expression statement")
	if c.Match(TOKEN_CR) {

	}
}

func (c *Compiler) Block() {
	for !c.Check(TOKEN_RIGHT_BRACE) && !c.Check(TOKEN_EOF) {
		c.Statement()
	}
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func (c *Compiler) BeginScope() {
	ScopeId++
	c.ScopeDepth++
}

func (c *Compiler) EndScope() {
	c.ScopeDepth--

	for c.Current.LocalCount > 0 &&
		c.Current.Locals[c.Current.LocalCount-1].depth > c.ScopeDepth {
		if c.Current.Locals[c.Current.LocalCount-1].isCaptured {
			c.EmitOp(OP_CLOSE_UPVALUE)
		} else {
			c.EmitOp(OP_POP)
			c.WriteComment("POP after end of scope")
		}
		c.Current.LocalCount--
	}
}

func (c *Compiler) EmitJump(opcode byte) int {
	c.EmitInstr(opcode, int16(9999))
	curPos := c.CurrentInstructions().CurrentBytePosition()
	c.WriteComment(fmt.Sprintf("Jump from %d", curPos))
	return c.CurrentInstructions().Count - 1
}

func (c *Compiler) PatchJump(instrNumber int) {

	jump := c.CurrentInstructions().JumpFrom(instrNumber)
	currentLocation := c.CurrentInstructions().NextBytePosition()
	byteJump := currentLocation - jump

	c.CurrentInstructions().SetOperand(instrNumber, int16(byteJump))
}

func (c *Compiler) PatchBreaks() {
	// Look for all the breaks in this loop
	bp := BreakPtr
	for i := 0; i < bp; i++ {
		if Breaks[i].CanPatch {
			startLoc := Breaks[i].StartLoc

			c.PatchJump(startLoc)
			Breaks[i].CanPatch = false
			BreakPtr--
		}
	}
}

func (c *Compiler) BreakStatement() {
	if PeekLoop() == LOOP_WHILE {
		Breaks[BreakPtr].StartLoc = c.EmitJump(OP_JUMP)
		Breaks[BreakPtr].CanPatch = true
		BreakPtr++
	} else if PeekLoop() == LOOP_FOR {
		c.EmitOp(OP_BREAK)
	}
}

func (c *Compiler) ContinueStatement() {
	if PeekLoop() == LOOP_WHILE {
		curLoc := c.CurrentInstructions().NextBytePosition() //+ 3
		start := StartLoop[StartPtr]
		offSet := start - curLoc+3
		c.EmitInstr(OP_JUMP, int16(offSet))
		c.WriteComment(fmt.Sprintf("Continue to %d from %d by offset %d", start, curLoc, offSet))
	} else if PeekLoop() == LOOP_FOR {
		c.EmitOp(OP_CONTINUE)
	}
}

func (c *Compiler) IfStatement() {
	// This part handles the logical 'if' condition
	c.Expression()

	thenJump := c.EmitJump(OP_JUMP_IF_FALSE)
	c.Statement()
	elseJump := c.EmitJump(OP_JUMP)

	c.PatchJump(thenJump)
	if c.Match(TOKEN_ELSE) {
		c.Statement()
	}
	c.PatchJump(elseJump)

}

func (c *Compiler) EmitLoop(offset int) {
	c.EmitInstr(OP_JUMP, int16(offset)-3)
	currByte := c.CurrentInstructions().CurrentBytePosition()
	c.WriteComment(fmt.Sprintf("Jump to %d", currByte+offset))
}

func (c *Compiler) ScanStatement() {
	c.BeginScope()
	c.Expression() // Array

	// Manages the target variable
	c.Consume(TOKEN_TO, "Expect 'to' after the object declaration")
	c.Consume(TOKEN_IDENTIFIER, "Expect variable name after 'to'")

	idx := c.AddLocal(c.Parser.Previous.ToString())
	c.EmitInstr(OP_PUSH, idx)
	c.WriteComment(fmt.Sprintf("Push variable index %d", idx))

	// Keeps track of the iterator
	reg := c.GetFreeRegister()
	c.EmitInstr(OP_PUSH, reg)
	c.WriteComment(fmt.Sprintf("Push register index %d", reg))

	scanJump := c.EmitJump(OP_SCAN)

	// Run the body of the code
	c.Evaluate()
	c.EmitOp(OP_CONTINUE)

	c.PatchJump(scanJump)

	c.EndScope()
	c.FreeRegister(reg)
}

func (c *Compiler) ForStatement() {
	PushLoop(LOOP_FOR)
	c.BeginScope()

	/*
		The variable we use to hold the loop initializer and iterator value is actually
		stored in a register and not a regular variable. We do this to make it easier to
		manipulate the value inside a VM function we use to execute the loop
	*/
	var varName string
	// This is the declaration portion
	if c.Match(TOKEN_IDENTIFIER) {
		// Get the name of the variable
		varName = c.Parser.Previous.ToString()
		// Push the initial value of the initializer on to the stack
		if c.Match(TOKEN_EQUAL) {
			c.Expression()
		}
	} else {
		c.ErrorAtCurrent("FOR initialized incorrectly")
	}

	// To value
	c.Consume(TOKEN_TO, "'to' is required after variable assignment")
	c.Expression()

	// Step
	if c.Match(TOKEN_STEP) {
		c.Expression()
	} else {
		c.EmitOp(OP_PUSH_1)
		PushExpressionValue(ExpressionData{Value: VAL_INTEGER, ObjType: VAR_SCALAR})
	}

	// Here is where we assign a variable name to the register
	ridInit := c.GetFreeRegister()
	namedRegisters[varName] = ridInit

	c.EmitInstr(OP_PUSH, ridInit)
	c.WriteComment(fmt.Sprintf("Index for register %d", ridInit))

	c.EmitInstr(OP_FOR_LOOP, 9999)
	c.WriteComment("Execute this many bytes in a loop")
	currInstr := c.CurrentInstructions().Count - 1
	start := c.CurrentInstructions().NextBytePosition()

	c.Evaluate()
	c.EmitOp(OP_CONTINUE)

	end := c.CurrentInstructions().NextBytePosition()
	// Backpatch the number of bytes the body of the loop takes up
	c.CurrentInstructions().OpCode[currInstr].Operand = Int16ToBytes(int16(end - start))

	// Free the register for future use
	c.FreeRegister(ridInit)
	delete(namedRegisters, varName)

	c.EndScope()
	PopLoop()
}

func (c *Compiler) WhileStatement() {

	PushLoop(LOOP_WHILE)

	start := c.CurrentInstructions().NextBytePosition()

	// Logical expression should return a boolean
	c.BeginScope()
	c.Expression()

	exitJump := c.EmitJump(OP_JUMP_IF_FALSE)
	//c.EmitOp(OP_POP)

	c.Statement()
	c.EndScope()

	end := c.CurrentInstructions().NextBytePosition()

	c.EmitLoop(-(end - start))
	c.PatchJump(exitJump)

	//c.EmitOp(OP_POP)
	c.PatchBreaks()

	//StartPtr--

	PopLoop()
}

func (c *Compiler) SwitchStatement() {
	c.Expression()
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' after SWITCH <expr>")
	// Clear out all CR after the brace
	for c.Match(TOKEN_CR) {
	}
	var leaveJump = make([]int, 1024)
	cases := 0

	// Get the expression statement we're interested in

	for c.Match(TOKEN_CR) {
	}

	idx := GlobalCount
	c.EmitInstr(OP_SET_GLOBAL, idx)
	GlobalCount++

	// Go over all the cases
	for c.Match(TOKEN_WHEN) || c.Match(TOKEN_DEFAULT) {
		cmd := c.Parser.Previous

		if cmd.Type == TOKEN_WHEN {
			// Get the value we're comparing against
			c.Expression()
			// If it didn't match ..
			c.EmitInstr(OP_GET_GLOBAL, idx)
			c.EmitOp(OP_EQUAL)
			exitJump := c.EmitJump(OP_JUMP_IF_FALSE)
			// If it matches we go ahead
			if c.Match(TOKEN_COLON) {
				for c.Match(TOKEN_CR) {
				}
				c.Evaluate()
				for c.Match(TOKEN_CR) {
				}
			}

			// .. and then leave the case statement
			leaveJump[cases] = c.EmitJump(OP_JUMP)
			cases++

			c.PatchJump(exitJump)
		} else if cmd.Type == TOKEN_DEFAULT {
			if c.Match(TOKEN_COLON) {
				for c.Match(TOKEN_CR) {
				}
				c.Evaluate()
				for c.Match(TOKEN_CR) {
				}
			}
			leaveJump[cases] = c.EmitJump(OP_JUMP)
			cases++
		}

	}
	// Ignore all CR's until the closing brace
	for c.Match(TOKEN_CR) {
	}
	// No matter what, we end up here
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' at end of CASE statement")
	for i := 0; i < cases; i++ {
		c.PatchJump(leaveJump[i])
	}
	c.EmitOp(OP_POP)
}

func (c *Compiler) CaseStatement() {
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' after CASE <expr>")
	// Clear out all CR after the brace
	for c.Match(TOKEN_CR) {
	}
	var leaveJump = make([]int, 1024)
	cases := 0
	// Go over all the cases
	for c.Match(TOKEN_WHEN) || c.Match(TOKEN_DEFAULT) {
		cmd := c.Parser.Previous

		if cmd.Type == TOKEN_WHEN {
			// Get the value we're comparing against
			c.Expression()
			// If it didn't match ..

			exitJump := c.EmitJump(OP_JUMP_IF_FALSE)
			// If it matches we go ahead
			if c.Match(TOKEN_COLON) {
				for c.Match(TOKEN_CR) {
				}
				c.Evaluate()
				for c.Match(TOKEN_CR) {
				}
			}

			// .. and then leave the case statement
			leaveJump[cases] = c.EmitJump(OP_JUMP)
			cases++

			c.PatchJump(exitJump)
		} else if cmd.Type == TOKEN_DEFAULT {
			if c.Match(TOKEN_COLON) {
				for c.Match(TOKEN_CR) {
				}
				c.Evaluate()
				for c.Match(TOKEN_CR) {
				}
			}
			leaveJump[cases] = c.EmitJump(OP_JUMP)
			cases++
		}

	}
	// Ignore all CR's until the closing brace
	for c.Match(TOKEN_CR) {
	}
	// No matter what, we end up here
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' at end of CASE statement")
	for i := 0; i < cases; i++ {
		c.PatchJump(leaveJump[i])
	}
	c.EmitOp(OP_POP)
}

func (c *Compiler) AddProperty(class *ClassVar, name string) {

	prop := &class.Properties[class.PropertyCount]

	prop.Name = name
	//prop.ObjType = expData.ObjType
	//prop.DataType = expData.Value

	prop.Index = class.PropertyCount
	prop.EnclosingClass = class

	idx := c.MakeConstant(ObjString(name))

	c.EmitInstr(OP_BIND_PROPERTY,idx)
	c.WriteComment(fmt.Sprintf("Property name %s index %d",prop.Name,prop.Index))

	class.PropertyCount++

	class.Class.FieldCount = class.PropertyCount - 1

}

func (c *Compiler) this_(canAssign bool) {
	c.Variable(false)
}

func (c *Compiler) GetAccessor() AccessorType {
	switch {
		case c.Match(TOKEN_PRIVATE): return PRIVATE
		case c.Match(TOKEN_PROTECTED): return PROTECTED
		default : return PUBLIC
	}
}

func (c *Compiler) Class(canAssign bool) {

	class := &ObjClass{
		Id:          ClassId,
		Class:       nil,
		Fields:      nil,
		FieldCount:  0,
		Methods:     nil,
		MethodCount: 0,
	}

	vclass := NewClassVar()
	vclass.Class = class

	ClassId++
	vclass.Id = ClassId

	c.EmitOp(OP_CLASS)

	c.Consume(TOKEN_LEFT_BRACE,"Expect '{' after class name")
	for {

		// Find out if it's public, protected, or private
		_ = c.GetAccessor()
		// Is it a property? if it is, it'll have a data type indicator. If not,
		// it means it's a method
		expData := c.GetDataType()
		// Either way, the next token needs to be the name
		c.Consume(TOKEN_IDENTIFIER,"Expect name of class component")
		compName := c.Parser.Previous.ToString()
		//idx := c.MakeConstant(ObjString(compName))
		if expData.ObjType == VAR_UNKNOWN {
			// It's a method .. so let's make one
			c.Procedure(TYPE_METHOD)
			c.AddProperty(&vclass, compName)
		} else {
			c.EmitOp(OP_NIL)
			c.AddProperty(&vclass, compName)
		}

		// Ok, we're done
		if c.Match(TOKEN_RIGHT_BRACE) {
			break
		}
	}

	ClassId--
}

func (c *Compiler) Method(canAssign bool) {
	c.Procedure(TYPE_METHOD)
}

func (c *Compiler) Function(canAssign bool) {
	c.Procedure(TYPE_FUNCTION)
}

func (c *Compiler) CallNative(nativeFunction *ObjNative) {
	c.Consume(TOKEN_LEFT_PAREN, "Expect '(' before native call")
	argumentCount := c.GetArguments()
	idx := c.MakeConstant(nativeFunction)
	c.EmitInstr(OP_CALL_NATIVE, idx)
	c.EmitOperand(argumentCount)
	PushExpressionValue(nativeFunction.ReturnType)
}

func (c *Compiler) Procedure(functionType FunctionType) {

	// Create the function object we're going to fill
	fn := &FunctionVar{
		paramCount: 0,
		instr:      NewInstructions(),
		returnType: VAL_NIL,
		Locals:     make([]Local, 65000),
		Upvalues:   make([]Upvalue, 65000),
	}

	// Set the current function as the enclosing function of this new function
	fn.Enclosing = c.Current
	// Make the current function into the new function
	c.Current = fn

	c.BeginScope()

	paramCount := int16(0)

	c.Current.Locals[c.Current.LocalCount].depth = 0
	c.Current.Locals[c.Current.LocalCount].isCaptured = false

	// if it's a method, we add the current class as a parameter
	if functionType == TYPE_METHOD {
		c.Current.Locals[c.Current.LocalCount].name = "this"
		c.Current.Locals[c.Current.LocalCount].ExprData.Value = VAL_CLASS
		c.Current.Locals[c.Current.LocalCount].ExprData.ObjType = VAR_CLASS
		c.Current.Locals[c.Current.LocalCount].Class = CurrentClass
		paramCount++
	}

	// Set up the locals for this function
	c.Current.LocalCount++

	// Parenthesis and parameter definition

	c.Consume(TOKEN_LEFT_PAREN, "Expect '(' after function definition.")
	// Here we just count the parameters

	for !c.Check(TOKEN_RIGHT_PAREN) {
		for {
			paramCount++
			if paramCount > 1024 {
				c.ErrorAtCurrent("Cannot have more than 1024 parameters.")
			}
			c.DefineParameter()
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
	}

	c.Current.paramCount = paramCount
	c.Consume(TOKEN_RIGHT_PAREN, "Expect ')' after parameters.")

	// If there is a return value, then declare it here
	isReturnValue := false
	c.Current.returnType = c.GetDataType().Value
	if c.Current.returnType != VAL_NIL {
		isReturnValue = true
	}

	// Body of the function
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' before function body.")
	c.Block()
	// If this function returns nothing, then return nil
	if !isReturnValue {
		c.EmitOp(OP_NIL)
		c.WriteComment("In lieu of explicit return value")
	}
	c.EmitReturn()
	c.EndScope()
	// Display
	if c.DebugMode {
		if functionType == TYPE_FUNCTION {
			fmt.Print("=== Function ===\n")
		} else {
			fmt.Print("=== METHOD ===\n")
		}
		fmt.Printf("Parameters: %d\n", c.Current.paramCount)
		fmt.Printf("Closures: %d\n", c.Current.UpvalueCount)
		c.Current.instr.Display()
	}
	// Return back to the calling function
	prev := c.Current
	c.Current = c.Current.Enclosing

	// The function is created by now, so lets turn it into a chunk
	// and push the value on to the stack
	idx := c.MakeConstant(prev.ConvertToObj())

	// Pop out of this function definition
	c.EmitInstr(OP_CLOSURE, idx)
	c.WriteComment(fmt.Sprintf("Closure index %d scope %d upvalues: %d", idx, c.ScopeDepth, c.Current.UpvalueCount))

	for i := int16(0); i < prev.UpvalueCount; i++ {
		b1 := byte(0)
		if prev.Upvalues[i].IsLocal {
			b1 = 1
		}
		c.EmitOperandByte(b1)
		c.EmitOperand(prev.Upvalues[i].Index)
	}

	PushExpressionValue(ExpressionData{
		Value:   VAL_FUNCTION,
		ObjType: VAR_FUNCTION,
	})

}

func (c *Compiler) CreateTable() {
	c.Consume(TOKEN_IDENTIFIER,"Expect table name after 'CREATE TABLE'")
	tblName := c.Parser.Previous.ToString()
	c.ClearCR()

	tbl := CreateTable(tblName)

	c.Consume(TOKEN_LEFT_PAREN,"Expect '(' after 'CREATE TABLE' statement")
	c.ClearCR()
	for {
		// Column name
		c.Consume(TOKEN_IDENTIFIER,"Expect column name")
		colName := c.Parser.Previous.ToString()

		// Column type
		dType := c.GetDataType()

		tbl.AddColumn(colName, dType.Value)

		c.ClearCR()
		if !c.Match(TOKEN_COMMA) {
			break
		}
	}
	c.Consume(TOKEN_RIGHT_PAREN,"Expect ')' after complete 'CREATE TABLE' statement")

	idx := c.MakeConstant(&tbl)

	c.EmitInstr(OP_CREATE_TABLE,idx)

}

func (c *Compiler) CreateStatement() {
	switch {
		case c.Match(TOKEN_TABLE) : c.CreateTable()
		case c.Match(TOKEN_INDEX) :
	}
}

func (c *Compiler) InsertStatement() {
	switch {
	case c.Match(TOKEN_INTO) : c.InsertInto()
	}
}

func (c *Compiler) InsertInto() {
	c.Consume(TOKEN_IDENTIFIER,"Expect table name after 'INSERT INTO'")
	tblName := c.Parser.Previous.ToString()
	tblIdx := c.MakeConstant(ObjString(tblName))

	colCount := 0

	if c.Match(TOKEN_LEFT_PAREN) {
		// We need to specify which columns
		for {
			c.Consume(TOKEN_IDENTIFIER,"Expect column names")
			colName := c.Parser.Previous.ToString()
			c.EmitInstr(OP_PUSH, c.MakeConstant(ObjString(colName)))
			colCount++
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.Consume(TOKEN_RIGHT_PAREN, "Expect ') after column list'")
	} else {
		colCount = 0
		// We use all columns
	}

	valCount := 0
	if c.Match(TOKEN_VALUES) {
		c.Consume(TOKEN_LEFT_PAREN,"Expect '(' after 'VALUES'")
		for {
			c.Expression()
			valCount++
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.Consume(TOKEN_RIGHT_PAREN, "Expect ') after values list'")
	}
	
	c.EmitInstr(OP_INSERT, tblIdx)
	c.EmitOperand(int16(colCount))

}

func (c *Compiler) SelectStatement() {

	fldCount := int16(0)
	if c.Match(TOKEN_STAR) || c.Match(TOKEN_ALL) {
		// All columns
	} else {
		// Get the columns as named
		for {
			c.Consume(TOKEN_IDENTIFIER,"Expect field name")
			fldCount++
			c.EmitInstr(OP_PUSH,c.MakeConstant(ObjString(c.Parser.Previous.ToString())))
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.EmitInstr(OP_PUSH,fldCount)
	}
	if c.Match(TOKEN_FROM) {
		c.Consume(TOKEN_IDENTIFIER,"Expect table name")
		c.EmitInstr(OP_PUSH,c.MakeConstant(ObjString(c.Parser.Previous.ToString())))
		c.EmitOp(OP_SQL_SELECT)
	}

}

func (c *Compiler) Allocate() {
	c.New(true)
}

func (c *Compiler) Statement() {

	switch {
		case c.Match(TOKEN_VAR): 		c.DeclareVariable()
		case c.Match(TOKEN_NEW):        c.Allocate()
		case c.Match(TOKEN_IF):			c.IfStatement()
		case c.Match(TOKEN_RETURN):		c.ReturnStatement()
		case c.Match(TOKEN_SCAN): 		c.ScanStatement()
		case c.Match(TOKEN_FOR): 		c.ForStatement()
		case c.Match(TOKEN_WHILE): 		c.WhileStatement()
		case c.Match(TOKEN_SWITCH):		c.SwitchStatement()
		case c.Match(TOKEN_CASE): 		c.CaseStatement()
		case c.Match(TOKEN_LEFT_BRACE):
			c.BeginScope()
			c.Block()
			c.EndScope()
		case c.Match(TOKEN_BREAK): 		c.BreakStatement()
		case c.Match(TOKEN_CONTINUE): 	c.ContinueStatement()
		case c.Match(TOKEN_CR):
		case c.Match(TOKEN_CREATE): 	c.CreateStatement()
		case c.Match(TOKEN_INSERT): 	c.InsertStatement()
		case c.Match(TOKEN_SELECT): 	c.SelectStatement()
		default: c.ExpressionStatement()
	}
}

func (c *Compiler) Evaluate() {
	c.Statement()
}

func (c *Compiler) WriteComment(comment string) {
	c.CurrentInstructions().WriteComment(comment)
}

func (c *Compiler) PushType(valType ValueType) {
	if c.Match(TOKEN_LEFT_BRACKET) {
		dims := int16(0)
		for {
			dims++
			c.Expression()
			if !c.Match(TOKEN_COMMA) {
				break
			}
		}
		c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after array type initializer")

		c.EmitInstr(OP_PUSH,int16(valType))
		c.WriteComment(fmt.Sprintf("Specifying value type as %s", ValueTypeLabel[valType]))

		c.EmitInstr(OP_MAKE_ARRAY,dims)
		c.WriteComment(fmt.Sprintf("Make array with %d dimensions", dims))

		PushExpressionValue(ExpressionData{
			Value:      valType,
			ObjType:    VAR_ARRAY,
			Dimensions: int(dims),
		})
		// If there's a left brace, then it means we're also filling the value with data
		if c.Match(TOKEN_LEFT_BRACE) {
			c.Expression()
			c.Consume(TOKEN_RIGHT_BRACE, "Right brace expected after array expression")
		}
	} else {
		PushExpressionValue(ExpressionData{valType, VAR_SCALAR,1})
	}
}

func (c *Compiler) IntegerType(canAssign bool) {
	c.PushType(VAL_INTEGER)
}

func (c *Compiler) FloatType(canAssign bool) {
	c.PushType(VAL_FLOAT)
}

func (c *Compiler) BoolType(canAssign bool) {
	c.PushType(VAL_BOOL)
}

func (c *Compiler) StringType(canAssign bool) {
	c.PushType(VAL_STRING)
}

func GetTokenType(tokType TokenType) ValueType {
	switch tokType {
	case TOKEN_INTEGER:
		return VAL_INTEGER
	case TOKEN_DECIMAL:
		return VAL_FLOAT
	case TOKEN_TYPE_BOOL:
		return VAL_BOOL
	case TOKEN_STRING:
		return VAL_STRING
	default:
		return VAL_NIL
	}
}
