package main

import (
	"fmt"
)

var dfSequence int64

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

	fldCount := int16(v.Pop().(ObjInteger)) // Number of fields coming

	// Array that will hold the fields we're interested in
	cols := make([]string, fldCount)
	// Store the fields by name in the array
	for i := fldCount - 1; i >= 0; i-- {
		cols[i] = string(v.ReadConstant(int16(v.Pop().(ObjInteger))).(ObjString))
	}

	// Get a new, empty, result set
	rs := v.NewResultSet()
	// Get the names of the columns to correspond with the expected result set
	rs.ColNames = cols
	// Copy the column definitions from the source data to the result set
	for c := int16(0); c < fldCount; c++ {
		rs.Columns[cols[c]] = df.Columns[cols[c]]
		rs.Columns[cols[c]].Ordinal = c // The ordinal position is assigned according to the position of the resultset
	}
	// The new result set need to have the same number of columns as the query
	rs.ColumnCount = fldCount

	// Row values get assigned to the result set here
	vals := make([]Obj,fldCount) // Values are going to equal the number of columns (obviously)
	for i := int64(0); i < df.RowCount; i++ {
		// Get the value from the source table and copy it to the resultset
		for c := int16(0); c < fldCount; c++ {
			vals[c] = df.Columns[cols[c]].GetValue(i)
		}
		// Now that we've filled the columns with data for this row, add it to the result set
		rs.AddRow(cols,vals)
	}
	// This indicates that this table is usable
	rs.Defined = true
	v.Push(rs)
}


func (v *VM) NewResultSet() *ObjDataFrame {

	dfName := fmt.Sprintf("#df%08d",dfSequence)

	df := CreateTable(dfName)

	dfSequence++
	return &df
}

