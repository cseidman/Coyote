package compiler

import (
	. "../common"
	. "../instructions"
	. "../parser"
	. "../token"
	. "../value"
	"bytes"
	"fmt"
	"strconv"
)

// Keeps track of break and continue instruction locations
type Break struct {
	StartLoc int  // Continue will bump up to here
	CanPatch bool // Flag to indicate if this is waiting to be patched
}

var Breaks = make([]Break, 255)
var BreakPtr int

var StartLoop = make([]int, 255)
var StartPtr int

// This is the global variable space
type Global struct {
	name     Token
	datatype ValueType
}

var GlobalVars []Global
var GlobalCount int16

func LoadGlobals() {
	GlobalVars = make([]Global, 65000)
	GlobalCount = 0
}

func AddGlobal(tok Token) int16 {

	vType := GetTokenType(tok.Type)

	GlobalVars[GlobalCount] = Global{tok, vType}
	GlobalCount++
	return GlobalCount - 1
}

type Local struct {
	name       Token
	depth      int
	isCaptured bool
	dataType   ValueType
}

type Upvalue struct {
	Index    int
	IsLocal  bool
	dataType ValueType
}

var namedRegisters = make(map[string]int16)

type register struct {
	isUsed bool
}

type Compiler struct {
	Parser              *Parser
	CurrentInstructions *Instructions
	Rules               []ParseRule

	locals     []Local
	LocalCount int16

	registers []register

	Upvalues []Upvalue

	ScopeDepth int
}

func Compile(source *string) Instructions {

	// First parse the source
	parser := NewParser(source)
	instr := NewInstructions()
	compiler := Compiler{
		Parser:              &parser,
		CurrentInstructions: &instr,
		LocalCount:          0,
		ScopeDepth:          0,
		locals:              make([]Local, 16000),
		Upvalues:            make([]Upvalue, 16000),
		registers:           make([]register, 256),
	}

	compiler.LocalCount++
	compiler.locals[compiler.LocalCount].depth = 0
	compiler.locals[compiler.LocalCount].isCaptured = false
	compiler.locals[compiler.LocalCount].name.Length = 0
	compiler.locals[compiler.LocalCount].name.Value = []byte(nil)

	compiler.LoadRules()
	compiler.Advance()
	for !compiler.Match(TOKEN_EOF) {
		compiler.Evaluate()
	}
	compiler.ReturnStatement()

	instr.Display()

	return instr
}

func (c *Compiler) FreeRegister(location int16) {
	c.registers[location].isUsed = false
}

func (c *Compiler) GetFreeRegister() int16 {
	for i := int16(0); i < 256; i++ {
		if !c.registers[i].isUsed {
			c.registers[i] = register{
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
		for i = 0; i < c.CurrentInstructions.ConstantsCount; i++ {
			// it's a string
			if c.CurrentInstructions.Constants[i].Type() == VAL_STRING {
				// If the strings match
				if c.CurrentInstructions.Constants[i].(*ObjString).Value == value.(*ObjString).Value {
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
	c.CurrentInstructions.Constants[c.CurrentInstructions.ConstantsCount] = val
	c.CurrentInstructions.ConstantsCount++
	return c.CurrentInstructions.ConstantsCount - 1
}

func (c *Compiler) Advance() {
	c.Parser.Prev2 = c.Parser.Previous
	c.Parser.Previous = c.Parser.Current
	c.Parser.Current = c.Parser.NextToken()
}

func (c *Compiler) ResolveGlobal(tok *Token) int16 {
	var i int16
	for i = 0; i < GlobalCount; i++ {
		if tok.Type == GlobalVars[i].name.Type && tok.ToString() == GlobalVars[i].name.ToString() {
			return i
		}
	}
	c.ErrorAtCurrent(fmt.Sprintf("Global variable '%s' not found", tok.ToString()))
	return -1
}

func (c *Compiler) ResolveLocal(tok *Token) int16 {
	for i := c.LocalCount - 1; i >= 0; i-- {
		if c.IdentifiersEqual(tok.Value, c.locals[i].name.Value) {
			if c.locals[i].depth == -1 {
				c.Error("Cannot read local variable in its own initializer.")
			}
			return i
		}
	}
	return -1
}

func (c *Compiler) AddLocal(name Token) int16 {

	if c.LocalCount == 16000 {
		c.Error("Too many local variables in function.")
		return -1
	}

	c.locals[c.LocalCount].name = name
	c.locals[c.LocalCount].depth = c.ScopeDepth
	c.locals[c.LocalCount].isCaptured = false

	c.LocalCount++
	return c.LocalCount - 1

}

func (c *Compiler) NamedVariable(canAssign bool) {

	tok := c.Parser.Previous
	var getOp byte
	var setOp byte
	var valType ValueType

	// If it's a local variable, we look for that before globals
	idx := c.ResolveLocal(&tok)
	// -1 means it wasn't found
	if idx != -1 {
		getOp = OP_GET_LOCAL
		setOp = OP_SET_LOCAL
	} else {
		// Is it in a register?
		ridx, ok := namedRegisters[tok.ToString()]
		if ok {
			idx = ridx
			getOp = OP_GET_REGISTER
			setOp = OP_SET_REGISTER
			valType = VAL_INTEGER
		} else {
			idx = c.ResolveGlobal(&tok)
			setOp = OP_SET_GLOBAL
			if idx != -1 {
				getOp = OP_GET_GLOBAL
			}
		}
	}

	if canAssign && c.Match(TOKEN_EQUAL) {
		c.Expression()
		c.EmitInstr(setOp, idx)
		valType = c.CurrentInstructions.GetType(0)
		if setOp == OP_SET_GLOBAL {
			GlobalVars[idx].datatype = valType
		} else if setOp == OP_SET_LOCAL {
			c.locals[idx].dataType = valType
		}
		c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[setOp], tok.ToString(), idx, valType))
	} else {
		c.EmitInstr(getOp, idx)
		if getOp == OP_GET_GLOBAL {
			valType = GlobalVars[idx].datatype
		} else if getOp == OP_GET_LOCAL {
			valType = c.locals[idx].dataType
		}
		c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[getOp], tok.ToString(), idx, valType))
	}

	c.SetType(valType)
}

func (c *Compiler) IdentifierConstant() int16 {
	return c.MakeConstant(&ObjString{Value: string(c.Parser.Previous.Value)})
}

func (c *Compiler) DeclareVariable() {

	c.Consume(TOKEN_IDENTIFIER, "Expect variable name")

	// Store the token here
	tok := c.Parser.Previous

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
		for i := c.LocalCount - 1; i >= 0; i-- {
			if c.locals[i].depth != -1 &&
				c.locals[i].depth < c.ScopeDepth {
				break
			}

			if c.IdentifiersEqual(tok.Value, c.locals[i].name.Value) {
				c.Error(fmt.Sprintf("Variable with the name %s already declared in this scope.", tok.ToString()))
			}
		}
		opcode = OP_SET_LOCAL
		index = c.AddLocal(tok)
	}
	// If there is an assigment operator after this, then we pop the rvalue
	// on the stack as well
	if c.Match(TOKEN_EQUAL) {
		c.Expression()
		valType = c.CurrentInstructions.GetType(0)

	} else {
		c.EmitOp(OP_NIL)
	}
	c.EmitInstr(opcode, index)
	c.SetType(valType)
	c.WriteComment(fmt.Sprintf("%s name %s at index %d type %d", OpLabel[opcode], tok.ToString(), index, valType))

	switch opcode {
	case OP_SET_GLOBAL:
		GlobalVars[index].datatype = valType
	case OP_SET_LOCAL:
		c.locals[index].dataType = valType
	}

	c.Match(TOKEN_CR) // Remove any CR after the declaration

}

func (c *Compiler) IdentifiersEqual(a []byte, b []byte) bool {
	return bytes.Equal(a, b)
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
	// This tells the compiler that we're in error mode now,
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

func (c *Compiler) SetType(valueType ValueType) {
	c.CurrentInstructions.AddType(valueType)
}

func (c *Compiler) EmitOp(opcode byte) {
	c.CurrentInstructions.WriteSimpleInstruction(opcode)
}

func (c *Compiler) EmitInstr(opcode byte, operand int16) {
	c.CurrentInstructions.WriteInstruction(opcode, operand)
}

func (c *Compiler) ReturnStatement() {
	c.EmitOp(OP_RETURN)
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
func (c *Compiler) Call(canAssign bool) {}
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

	// Emit the operator instruction.
	switch operatorType {
	case TOKEN_BANG:
		c.EmitOp(OP_NOT)
	case TOKEN_MINUS:
		c.EmitOp(OP_NEGATE)
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

	ltype := c.CurrentInstructions.GetType(1)
	rtype := c.CurrentInstructions.GetType(0)

	if ltype != rtype {
		c.ErrorAtCurrent("Types don't match")
		return
	}
	// This is the value for both types
	dType := ltype

	// Emit the operator instruction.
	switch operatorType {
	case TOKEN_BANG_EQUAL:
		c.EmitOp(OP_NOT_EQUAL)
		c.SetType(VAL_BOOL)
	case TOKEN_EQUAL_EQUAL:
		c.EmitOp(OP_EQUAL)
		c.SetType(VAL_BOOL)
	case TOKEN_GREATER:
		c.EmitOp(OP_GREATER)
		c.SetType(VAL_BOOL)
	case TOKEN_GREATER_EQUAL:
		c.EmitOp(OP_GREATER_EQUAL)
		c.SetType(VAL_BOOL)
	case TOKEN_LESS:
		c.EmitOp(OP_LESS)
		c.SetType(VAL_BOOL)
	case TOKEN_LESS_EQUAL:
		c.EmitOp(OP_LESS_EQUAL)
		c.SetType(VAL_BOOL)
	case TOKEN_PLUS:
		if dType == VAL_INTEGER {
			c.EmitOp(OP_IADD)
			c.SetType(VAL_INTEGER)
		} else if dType == VAL_FLOAT {
			c.EmitOp(OP_FADD)
			c.SetType(VAL_FLOAT)
		} else {
			fmt.Printf("Error in the data types %d:%d\n", ltype, rtype)
		}
	case TOKEN_MINUS:
		if dType == VAL_INTEGER {
			c.EmitOp(OP_ISUBTRACT)
			c.SetType(VAL_INTEGER)
		} else if dType == VAL_FLOAT {
			c.EmitOp(OP_FSUBTRACT)
			c.SetType(VAL_FLOAT)
		}
	case TOKEN_STAR:
		if dType == VAL_INTEGER {
			c.EmitOp(OP_IMULTIPLY)
			c.SetType(VAL_INTEGER)
		} else if dType == VAL_FLOAT {
			c.EmitOp(OP_FMULTIPLY)
			c.SetType(VAL_FLOAT)
		}
	case TOKEN_SLASH:
		if dType == VAL_INTEGER {
			c.EmitOp(OP_IDIVIDE)
			c.SetType(VAL_INTEGER)
		} else if dType == VAL_FLOAT {
			c.EmitOp(OP_FDIVIDE)
			c.SetType(VAL_FLOAT)
		}
	case TOKEN_PLUS_PLUS:
		c.EmitOp(OP_INCREMENT)
		if dType == VAL_INTEGER {
			c.SetType(VAL_INTEGER)
		} else if dType == VAL_FLOAT {
			c.SetType(VAL_FLOAT)
		}

	default:
		return
	}
}
func (c *Compiler) Dot(canAssign bool)   {}
func (c *Compiler) Array(canAssign bool) {}
func (c *Compiler) Index(canAssign bool) {}
func (c *Compiler) HMap(canAssign bool)  {}
func (c *Compiler) Variable(canAssign bool) {
	c.NamedVariable(canAssign)
}
func (c *Compiler) String(canAssign bool) {
	value := c.Parser.Previous.ToString()
	// Remove the quotes
	idx := c.MakeConstant(&ObjString{Value: value[1:(len(value) - 1)]})
	c.EmitInstr(OP_SCONST, idx)

	if len(value) > 10 {
		value = value[0:10]
	}

	c.WriteComment(fmt.Sprintf("Value %s at constant index %d", value, idx))
	c.SetType(VAL_STRING)
}
func (c *Compiler) Integer(canAssign bool) {
	value, _ := strconv.ParseInt(string(c.Parser.Previous.Value), 10, 64)
	idx := c.MakeConstant(&ObjInteger{Value: value})
	c.EmitInstr(OP_ICONST, idx)
	c.WriteComment(fmt.Sprintf("Value %d at constant index %d", value, idx))
	c.SetType(VAL_INTEGER)
}
func (c *Compiler) Float(canAssign bool) {
	value, _ := strconv.ParseFloat(string(c.Parser.Previous.Value), 64)
	idx := c.MakeConstant(&ObjFloat{Value: value})
	c.EmitInstr(OP_FCONST, idx)
	c.SetType(VAL_FLOAT)

}
func (c *Compiler) Browse(canAssign bool)    {}
func (c *Compiler) and_(canAssign bool)      {}
func (c *Compiler) or_(canAssign bool)       {}
func (c *Compiler) Literal(canAssign bool)   {}
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
}

func (c *Compiler) Block() {
	for !c.Check(TOKEN_RIGHT_BRACE) && !c.Check(TOKEN_EOF) {
		c.Evaluate()
	}
	c.Consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.")
}

func (c *Compiler) BeginScope() {
	c.ScopeDepth++
}

func (c *Compiler) EndScope() {
	c.ScopeDepth--

	for c.LocalCount > 0 &&
		c.locals[c.LocalCount-1].depth > c.ScopeDepth {
		if c.locals[c.LocalCount-1].isCaptured {
			c.EmitOp(OP_CLOSE_UPVALUE)
		} else {
			c.EmitOp(OP_POP)
			//c.WriteComment(currentChunk(), "POP after end of scope")
		}
		c.LocalCount--
	}
}

func (c *Compiler) EmitJump(opcode byte) int {
	c.EmitInstr(opcode, int16(9999))
	curPos := c.CurrentInstructions.CurrentBytePosition()
	c.WriteComment(fmt.Sprintf("Jump from %d", curPos))
	return c.CurrentInstructions.Count - 1
}

func (c *Compiler) PatchJump(instrNumber int) {

	jump := c.CurrentInstructions.JumpFrom(instrNumber)
	currentLocation := c.CurrentInstructions.NextBytePosition()
	byteJump := currentLocation - jump

	c.CurrentInstructions.SetOperand(instrNumber, int16(byteJump))
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
	Breaks[BreakPtr].StartLoc = c.EmitJump(OP_JUMP)
	Breaks[BreakPtr].CanPatch = true
	BreakPtr++
}

func (c *Compiler) ContinueStatement() {
	curLoc := c.CurrentInstructions.NextBytePosition() + 3
	start := StartLoop[StartPtr]
	offSet := start - curLoc
	c.EmitInstr(OP_JUMP, int16(offSet))
	c.WriteComment(fmt.Sprintf("Continue to %d from %d by offset %d", start, curLoc, offSet))
}

func (c *Compiler) IfStatement() {
	// This part handles the logical 'if' condition
	//c.Consume(TOKEN_LEFT_PAREN, "Expect '(' after 'if'")
	c.Expression()
	//c.Consume(TOKEN_RIGHT_PAREN, "Expect ')' to close 'if' condition")

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
	currByte := c.CurrentInstructions.CurrentBytePosition()
	c.WriteComment(fmt.Sprintf("Jump to %d", currByte+offset))
}

func (c *Compiler) ForStatement() {
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
		c.EmitInstr(OP_PUSH, 1)
		c.SetType(VAL_INTEGER)
	}

	// Here is where we assign a variable name to the register
	ridInit := c.GetFreeRegister()
	namedRegisters[varName] = ridInit

	c.EmitInstr(OP_PUSH, int16(ridInit))
	c.WriteComment(fmt.Sprintf("Index for register %d", ridInit))

	c.EmitInstr(OP_FOR_LOOP, 9999)
	c.WriteComment("Execute this many bytes in a loop")
	currInstr := c.CurrentInstructions.Count - 1
	start := c.CurrentInstructions.NextBytePosition()

	c.Evaluate()

	end := c.CurrentInstructions.NextBytePosition()
	// Backpatch the number of bytes the body of the loop takes up
	c.CurrentInstructions.OpCode[currInstr].Operand = int16(end - start)

	// Free the register for future use
	c.FreeRegister(ridInit)
	delete(namedRegisters, varName)

	c.EndScope()
}

func (c *Compiler) WhileStatement() {

	start := c.CurrentInstructions.NextBytePosition()
	StartPtr++
	StartLoop[StartPtr] = start

	// Logical expression should return a boolean
	c.BeginScope()
	c.Expression()

	exitJump := c.EmitJump(OP_JUMP_IF_FALSE)
	c.Evaluate()
	c.EndScope()

	end := c.CurrentInstructions.NextBytePosition()

	c.EmitLoop(-(end - start))
	c.PatchJump(exitJump)
	c.PatchBreaks()

	StartPtr--

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
	} else if c.Match(TOKEN_PRINT) {
		c.PrintStatement()
	} else if c.Match(TOKEN_IF) {
		c.IfStatement()
	} else if c.Match(TOKEN_RETURN) {
		c.ReturnStatement()
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

func (c *Compiler) PrintStatement() {
	c.Expression()
	c.Consume(TOKEN_CR, "Expect 'CR' after value.")
	c.EmitOp(OP_PRINT)
}

func (c *Compiler) Evaluate() {
	c.Statement()
}

func (c *Compiler) WriteComment(comment string) {
	c.CurrentInstructions.WriteComment(comment)
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
