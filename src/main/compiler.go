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

// Classes
var ClassBag = make(map[string]*ClassVar)

// Keeps track of break and continue instruction locations
type Break struct {
	StartLoc int  // Continue will bump up to here
	CanPatch bool // Flag to indicate if this is waiting to be patched
}

var Breaks = make([]Break, 255)
var BreakPtr int

var StartLoop = make([]int, 255)
var StartPtr int

// We use this to assign a unique if to every scope defined in the application
// so that we can know later on that variables declared in a given scope are different
// to variables of the same name are indeed different despite being in the same depth
var ScopeId = -1

// Keeps track of the last expression's return value
type ExpressionData struct {
	Value   ValueType
	ObjType VarType
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

type PropertyVar struct {
	EnclosingClass *ClassVar
	Access         AccessorType
	Name           string
	Index          int16
	DataType       ValueType
	ObjType        VarType
	HasValue       bool
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

// This is the global variable space
type Global struct {
	name          Token
	datatype      ValueType
	objtype       VarType
	IsInitialized bool
	Class         *ClassVar
}

var GlobalVars = make([]Global, 65000)
var GlobalCount = int16(0)

func AddGlobal(tok Token) int16 {

	vType := GetTokenType(tok.Type)

	GlobalVars[GlobalCount].name = tok
	GlobalVars[GlobalCount].datatype = vType
	GlobalVars[GlobalCount].objtype = VAR_UNKNOWN

	GlobalCount++
	return GlobalCount - 1
}

type Local struct {
	name          string
	depth         int
	isCaptured    bool
	dataType      ValueType
	scopeId       int
	objtype       VarType
	IsInitialized bool
	Class         *ClassVar
}

type Upvalue struct {
	Index    int16
	IsLocal  bool
	dataType ValueType
	Class    *ClassVar
}

var namedRegisters = make(map[string]int16)

type register struct {
	isUsed bool
}

var registers = make([]register, 256)

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
		if tok.Type == GlobalVars[i].name.Type && tok.ToString() == GlobalVars[i].name.ToString() {
			return i, &ExpressionData{
				Value:   GlobalVars[i].datatype,
				ObjType: GlobalVars[i].objtype,
			}
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

			return i, &ExpressionData{Value: fn.Locals[i].dataType, ObjType: fn.Locals[i].objtype}
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
	fn.Upvalues[upvCount].dataType = fn.Enclosing.Locals[index].dataType
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
	// Needs to return an integer
	c.Expression()
	indexType := PopExpressionValue().Value
	if indexType != VAL_INTEGER {
		c.Error("Array index must be an integer value")
	}

	c.Consume(TOKEN_RIGHT_BRACKET, "Expect ']' after array expression")

	idx, _, varscope := c.ResolveVariable(tok)
	if c.Match(TOKEN_EQUAL) {
		c.Expression()
		if varscope == GLOBAL {
			c.EmitInstr(OP_SET_AGLOBAL, idx)
		} else {
			c.EmitInstr(OP_SET_ALOCAL, idx)
		}
		c.WriteComment(fmt.Sprintf("Array name %s Index %d", tok.ToString(), idx))
	} else {
		if varscope == GLOBAL {
			c.EmitInstr(OP_GET_AGLOBAL, idx)
		} else {
			c.EmitInstr(OP_GET_ALOCAL, idx)
		}
		c.WriteComment(fmt.Sprintf("Array name %s Index %d", tok.ToString(), idx))
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
	if c.Match(TOKEN_LEFT_BRACKET) {
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
		//if isHasOperand {
		c.EmitInstr(setOp, idx)
		//} else {
		//	c.EmitOp(setOp)
		//}
		data := PopExpressionValue()

		valType = data.Value
		objType = data.ObjType

		if isGlobal {

			GlobalVars[idx].IsInitialized = true
			GlobalVars[idx].datatype = valType
			GlobalVars[idx].objtype = objType

			if objType == VAR_CLASS {

				GlobalVars[idx].Class = CurrentClass
			}

		} else if isLocal {
			if c.Current.Locals[idx].IsInitialized && valType != c.Current.Locals[idx].dataType {
				c.Error("Cannot assign incompatible variable")
			}

			c.Current.Locals[idx].IsInitialized = true
			c.Current.Locals[idx].dataType = valType
			c.Current.Locals[idx].objtype = objType

			if objType == VAR_CLASS {
				c.Current.Locals[idx].Class = CurrentClass
			}

		} else if isUpvalue {
			c.Current.Upvalues[idx].dataType = valType
			if objType == VAR_CLASS {
				c.Current.Upvalues[idx].Class = CurrentClass
			}
		}

		c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[setOp], tok.ToString(), idx, valType))
	} else {
		if isHasOperand {
			c.EmitInstr(getOp, idx)
		} else {
			c.EmitOp(getOp)
		}
		if isGlobal {
			valType = GlobalVars[idx].datatype
			objType = GlobalVars[idx].objtype
			if objType == VAR_CLASS {
				CurrentClass = GlobalVars[idx].Class
			}

		} else if isLocal {
			valType = c.Current.Locals[idx].dataType
			objType = c.Current.Locals[idx].objtype
			if objType == VAR_CLASS {
				CurrentClass = c.Current.Locals[idx].Class
			}
		} else if isUpvalue {
			valType = c.Current.Upvalues[idx].dataType
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
	c.Current.Locals[index].dataType = valType

}

func (c *Compiler) GetDataType() ExpressionData {

	expd := new(ExpressionData)

	if c.Check(TOKEN_TYPE_INTEGER) {
		expd.Value = VAL_INTEGER
		expd.ObjType = VAR_SCALAR
	} else if c.Check(TOKEN_TYPE_FLOAT) {
		expd.Value = VAL_FLOAT
		expd.ObjType = VAR_SCALAR
	} else if c.Check(TOKEN_TYPE_STRING) {
		expd.Value = VAL_STRING
		expd.ObjType = VAR_SCALAR
	} else if c.Check(TOKEN_FUNC) {
		expd.Value = VAL_FUNCTION
		expd.ObjType = VAR_FUNCTION
	} else if c.Check(TOKEN_CLASS) {
		expd.Value = VAL_CLASS
		expd.ObjType = VAR_CLASS
	} else {
		//c.Advance()
		expd.Value = VAL_NIL
		expd.ObjType = VAR_SCALAR
	}
	return *expd
}

func (c *Compiler) _array(canAssign bool) {
	//array()
}

func (c *Compiler) DeclareVariable() {

	c.Consume(TOKEN_IDENTIFIER, "Expect variable name")

	// Store the token here
	tok := c.Parser.Previous

	// Error if this variable collides with an existing native function name
	if ResolveNativeFunction(tok.ToString()) != nil {
		c.Error(fmt.Sprintf("'%s' is a reserved name", tok.ToString()))
	}

	var index int16
	var opcode byte
	var valType ValueType

	// At scopedepth 0 - it's a global
	if c.ScopeDepth == 0 {
		// So we add it to the global store
		index = AddGlobal(tok)
		opcode = OP_SET_GLOBAL
	} else {

		// Otherwise, it's a local variable
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
		opcode = OP_SET_LOCAL
		index = c.AddLocal(tok.ToString())
	}
	// If there is an assigment operator after this, then we pop the rvalue
	// on the stack as well
	var data ExpressionData
	if c.Match(TOKEN_EQUAL) {

		c.Expression()
		data = PopExpressionValue()
		valType = data.Value

	} else {

		c.EmitOp(OP_NIL)
		c.WriteComment("No equality token after variable declaration")
		valType = VAL_NIL

	}
	c.EmitInstr(opcode, index)
	PushExpressionValue(data)
	c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[opcode], tok.ToString(), index, valType))

	switch opcode {
	case OP_SET_GLOBAL:

		if GlobalVars[index].IsInitialized && valType != GlobalVars[index].datatype {
			c.Error("Cannot assign incompatible variable")
		}

		GlobalVars[index].IsInitialized = true
		GlobalVars[index].datatype = valType
		GlobalVars[index].objtype = data.ObjType

		if data.ObjType == VAR_CLASS {
			GlobalVars[index].Class = CurrentClass
		}

	case OP_SET_LOCAL:

		if c.Current.Locals[index].IsInitialized && valType != c.Current.Locals[index].dataType {
			c.Error("Cannot assign incompatible variable")
		}
		c.Current.Locals[index].IsInitialized = true
		c.Current.Locals[index].dataType = valType
		c.Current.Locals[index].objtype = data.ObjType
		if data.ObjType == VAR_CLASS {
			c.Current.Locals[index].Class = CurrentClass
		}

	}

	c.Match(TOKEN_CR) // Remove any CR after the declaration

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

	fmt.Printf("[line %d] Error", token.Line)
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
		c.Consume(TOKEN_CR, "Expect CR after return value")
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

	//rtype := c.CurrentInstructions().GetType(0)

	//if ltype != rtype {
	//	c.ErrorAtCurrent("Types don't match")
	//	return
	//}
	// This is the value for both types

	// Emit the operator instruction.
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
	default:
		return
	}
}

func (c *Compiler) FindPropertyType(class *ClassVar, propertyName string) ValueType {
	for i := int16(0); i < class.PropertyCount; i++ {
		if class.Properties[i].Name == propertyName {
			return class.Properties[i].DataType
		}
	}
	return VAL_NIL
}

func (c *Compiler) Dot(canAssign bool) {
	c.Consume(TOKEN_IDENTIFIER, "Expect property or method name after '.'.")
	name := c.Parser.Previous.ToString()

	isMethod := c.Match(TOKEN_LEFT_PAREN) // Is this is a method?

	idx := c.MakeConstant(ObjString(name))
	vType := c.FindPropertyType(CurrentClass, name)

	if isMethod {
		c.CallMethod(idx)
	} else {

		var getOp byte
		var setOp byte

		getOp = OP_GET_PROPERTY
		setOp = OP_SET_PROPERTY

		if canAssign && c.Match(TOKEN_EQUAL) {
			c.Expression()
			c.EmitInstr(setOp, idx)
		} else {
			c.EmitInstr(getOp, idx)
			c.WriteComment(fmt.Sprintf("Name '%s' of type %s", name, ValueTypeLabel[vType]))
			PushExpressionValue(ExpressionData{
				Value:   vType,
				ObjType: VAR_SCALAR,
			})
		}
	}
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
	c.ClearCR()
	for {
		c.Consume(TOKEN_IDENTIFIER, "Expect enum element")
		name := c.Parser.Previous.ToString()
		idx := c.MakeConstant(ObjString(name))
		c.EmitInstr(OP_PUSH, idx)
		elements++
		if elements > 255 {
			c.Error("Enums cannot have more than 255 elements")
		}
		c.ClearCR()
		if !c.Match(TOKEN_COMMA) {
			break
		}

		c.ClearCR()
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

	for {
		c.Expression()
		valType := c.GetDataType().Value
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
	c.EmitInstr(OP_ARRAY, int16(dType))
	PushExpressionValue(ExpressionData{
		Value:   dType,
		ObjType: VAR_ARRAY,
	})

}
func (c *Compiler) Index(canAssign bool) {}

func (c *Compiler) CompoundVariable(tok *Token) *ExpressionData {

	name := tok.ToString()

	if name == "this" {
		c.EmitOp(OP_GET_LOCAL_0)
		return &ExpressionData{
			Value:   VAL_CLASS,
			ObjType: VAR_CLASS,
		}
	}

	idx, expData := c.ResolveLocal(c.Current, name)

	// It's a local
	if idx != -1 {
		switch expData.ObjType {
		// Is it a class?
		case VAR_CLASS:
			// Then treat this as a Class
			c.EmitInstr(OP_GET_LOCAL, idx)
			return &ExpressionData{
				Value:   VAL_CLASS,
				ObjType: VAR_CLASS,
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
		case VAR_CLASS:
			// Treat this as a Class
			c.EmitInstr(OP_GET_GLOBAL, idx)
			return &ExpressionData{
				Value:   VAL_CLASS,
				ObjType: VAR_CLASS,
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

	tok := &c.Parser.Previous

	var expData *ExpressionData

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
		case VAR_CLASS:

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
		default:
			// Uh oh ..
			c.Error(fmt.Sprintf("Variable %s of type %s should not have a dot after it", tok.ToString(), VarTypeLabel[expData.ObjType]))
		}

	}

	if !foundCompoundObject {
		c.NamedVariable(canAssign)
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
	c.Consume(TOKEN_CR, "Expect 'CR' after expression.")
	c.EmitOp(OP_POP)
	c.WriteComment("Pop After expression statement")
}

func (c *Compiler) Block() {
	for !c.Check(TOKEN_RIGHT_BRACE) && !c.Check(TOKEN_EOF) {
		c.Evaluate()
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
		curLoc := c.CurrentInstructions().NextBytePosition() + 3
		start := StartLoop[StartPtr]
		offSet := start - curLoc
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
	//StartPtr++
	//StartLoop[StartPtr] = start

	// Logical expression should return a boolean
	c.BeginScope()
	c.Expression()

	exitJump := c.EmitJump(OP_JUMP_IF_FALSE)
	c.EmitOp(OP_POP)

	c.Evaluate()
	c.EndScope()

	end := c.CurrentInstructions().NextBytePosition()

	c.EmitLoop(-(end - start))
	c.PatchJump(exitJump)

	c.EmitOp(OP_POP)
	c.PatchBreaks()

	//StartPtr--

	PopLoop()
}

func (c *Compiler) SwitchStatement() {
	c.Expression()
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' after CASE <expr>")
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

func (c *Compiler) CreateClassComponent(class *ClassVar, tType TokenType) {
	c.Consume(TOKEN_IDENTIFIER, "Expect property or method name after access indicator")

	// Name of the property
	pName := c.Parser.Previous.ToString()
	// Use this to set the property
	idx := c.MakeConstant(ObjString(pName))
	c.AddProperty(class, pName)

	if c.Match(TOKEN_EQUAL) {
		if c.Match(TOKEN_METHOD) {
			c.Method(true)
			c.EmitInstr(OP_SET_PROPERTY, idx)
		} else {

		}
	}

}

func (c *Compiler) AddProperty(class *ClassVar, name string) {

	prop := &class.Properties[class.PropertyCount]

	prop.Name = name
	if c.Match(TOKEN_COLON) {
		expData := c.GetDataType()
		prop.ObjType = expData.ObjType
		prop.DataType = expData.Value

	} else {
		// This is a function
		//expData := c.GetDataType()
		prop.ObjType = VAR_FUNCTION
		prop.DataType = VAL_CLOSURE
	}

	prop.Index = class.PropertyCount
	prop.EnclosingClass = class

	class.PropertyCount++

	class.Class.FieldCount = class.PropertyCount - 1

}

func (c *Compiler) ClassComponents(class *ClassVar) {
	// Default Access indicator is public
	accessType := TOKEN_PUBLIC

	// If there's an explicit indicator, we consume it here
	if c.Match(TOKEN_PUBLIC) ||
		c.Match(TOKEN_PRIVATE) ||
		c.Match(TOKEN_PROTECTED) {

		accessType = c.Parser.Previous.Type
	}

	c.CreateClassComponent(class, accessType)
}

func (c *Compiler) this_(canAssign bool) {
	c.Variable(false)
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

	CurrentClass = &vclass
	c.Consume(TOKEN_LEFT_BRACE, "Expect '{' before class body.")
	c.ClearCR()
	// Keep going until we identified all properties and methods
	for !c.Match(TOKEN_RIGHT_BRACE) && !c.Match(TOKEN_EOF) {
		CurrentClass = &vclass
		c.ClassComponents(&vclass)
		c.ClearCR()
	}
	c.ClearCR()
	PushExpressionValue(ExpressionData{
		Value:   VAL_CLASS,
		ObjType: VAR_CLASS,
	})
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
		c.Current.Locals[c.Current.LocalCount].dataType = VAL_CLASS
		c.Current.Locals[c.Current.LocalCount].objtype = VAR_CLASS
		c.Current.Locals[c.Current.LocalCount].Class = CurrentClass
		paramCount++
	}

	// Set up the locals for this function
	c.Current.LocalCount++

	// Parenthesis and parameter definition

	c.Consume(TOKEN_LEFT_PAREN, "Expect '(' after function definition.")
	// Here we just count the parameters

	if !c.Check(TOKEN_RIGHT_PAREN) {
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
	c.ReturnStatement()
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
		Value:   prev.returnType,
		ObjType: VAR_FUNCTION,
	})

}

func (c *Compiler) Statement() {
	if c.Match(TOKEN_VAR) {
		c.DeclareVariable()
	} else if c.Match(TOKEN_UPDATE) {
		//UpdateStatement()
	} else if c.Match(TOKEN_INSERT) {
		//c.InsertStatement()
	} else if c.Match(TOKEN_CREATE) {
		//c.CreateStatement()
	} else if c.Match(TOKEN_INCLUDE) {
		//c.IncludeStatement()
	} else if c.Match(TOKEN_IF) {
		c.IfStatement()
	} else if c.Match(TOKEN_RETURN) {
		c.ReturnStatement()
	} else if c.Match(TOKEN_SCAN) {
		c.ScanStatement()
	} else if c.Match(TOKEN_FOR) {
		c.ForStatement()
	} else if c.Match(TOKEN_WHILE) {
		c.WhileStatement()
	} else if c.Match(TOKEN_SWITCH) {
		c.SwitchStatement()
	} else if c.Match(TOKEN_CASE) {
		c.CaseStatement()
	} else if c.Match(TOKEN_LEFT_BRACE) {
		c.BeginScope()
		c.Block()
		c.EndScope()
	} else if c.Match(TOKEN_BREAK) {
		c.BreakStatement()
	} else if c.Match(TOKEN_CONTINUE) {
		c.ContinueStatement()
	} else if c.Match(TOKEN_CR) {
		// Do nothing
	} else {
		c.ExpressionStatement()
	}
}

func (c *Compiler) Evaluate() {
	c.Statement()
}

func (c *Compiler) WriteComment(comment string) {
	c.CurrentInstructions().WriteComment(comment)
}

func (c *Compiler) IntegerType(canAssign bool) {
	c.EmitInstr(OP_PUSH, int16(VAL_INTEGER))
	c.WriteComment("Push the integer type to the stack")
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
