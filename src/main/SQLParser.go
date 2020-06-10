package main

//import ("fmt")

var SqlCmd string 

func (c *Compiler) SetSqlMode(mode bool) {
	c.Parser.TokenScanner.SQLMode = mode
}

func (c *Compiler) CreateTable() {
	SqlCmd = "CREATE TABLE "
	for !c.Match(TOKEN_SEMICOLON) {
		c.Advance()
		SqlCmd += c.Parser.Previous.ToString() + "\n"
	}
	idx := c.MakeConstant(ObjString(SqlCmd))
	//fmt.Println(SqlCmd)
	c.EmitInstr(OP_CREATE_TABLE,idx)
}

func (c *Compiler) SelectStatement() {
	SqlCmd = "SELECT "
	varCount := int16(0)
	for !c.Match(TOKEN_SEMICOLON) {
		if c.Match(TOKEN_DOLLAR) {
			varCount++
			// Get the variable on to the stack
			c.Advance()
			//tok := c.Parser.Previous
			c.NamedVariable(true)
			SqlCmd += "%v"
		}
		c.Advance()
		SqlCmd += c.Parser.Previous.ToString() + " "
	}
	idx := c.MakeConstant(ObjString(SqlCmd))
	//fmt.Println(SqlCmd)
	c.EmitInstr(OP_PUSH, varCount)
	c.EmitInstr(OP_SQL_SELECT,idx)
	//c.EmitOp(OP_DISPLAY_TABLE)
}

func (c *Compiler) CreateStatement() {
	switch {
	case c.Match(TOKEN_TABLE) : c.CreateTable()
	case c.Match(TOKEN_INDEX) :
	}
}

func (c *Compiler) InsertStatement() {
	SqlCmd = "INSERT "
	varCount := int16(0)
	for !c.Match(TOKEN_SEMICOLON) {
		if c.Match(TOKEN_DOLLAR) {
			varCount++
			// Get the variable on to the stack
			c.Advance()
			//tok := c.Parser.Previous
			c.NamedVariable(true)
			SqlCmd += "%v"
		}
		c.Advance()

		SqlCmd += c.Parser.Previous.ToString() + " "
	}
	SqlCmd+="\n"
	idx := c.MakeConstant(ObjString(SqlCmd))
	c.EmitInstr(OP_PUSH, varCount)
	c.EmitInstr(OP_INSERT, idx)
}

