package main

import (
	"fmt"
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
	tableName := string(v.ReadConstant(int16(v.Pop().(ObjInteger))).(ObjString))
	df := v.CurrDb.GetTable(tableName)

	fldCount := int(v.Pop().(ObjInteger)) // Number of fields coming

	// Array that will hold the fields we're interested in
	cols := make([]string, fldCount)
	// Store the fields by name in the array
	for i := fldCount - 1; i >= 0; i-- {
		cols[i] = string(v.ReadConstant(int16(v.Pop().(ObjInteger))).(ObjString))
	}

	resultDf := ObjDataFrame{
		Name:        ,
		Columns:     nil,
		ColNames:    nil,
		ColumnCount: 0,
		RowCount:    0,
		Defined:     false,
		RowCapacity: 0,
	}
	
	for c := 0; c < fldCount; c++ {
		fmt.Printf("%s\t", df.Columns[cols[c]].Name)
	}
	fmt.Println()
	for c := 0; c < fldCount; c++ {
		fmt.Printf("%s\t", "-------------")
	}
	fmt.Println()
	for i := int64(0); i < df.RowCount; i++ {
		for c := 0; c < fldCount; c++ {
			fmt.Printf("%s\t", df.Columns[cols[c]].GetValue(i).ShowValue())
		}
		fmt.Println()
	}
}


