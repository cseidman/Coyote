package main

import (
)


/*
This runs the OP_SQL_SELECT opcode
By the time thie function is called, the stack should look like this:
-- Top of the stack --
0: 		Constant index of the table name
1:		Number of requested columns
1-n: 	Constant indexes for the column names

*/
func (v *VM) RunSqlQuery() {

}

