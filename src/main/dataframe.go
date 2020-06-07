package main

import (
	"fmt"
	"database/sql"
)

type ObjDataFrame struct {
	DbRef *sql.DB
	Name string
	Columns []*sql.ColumnType
	Rows *sql.Rows
	ColNames []string
	ColumnCount int16
	CurrentRow *sql.Row
	RowNumber int64
}

// Interface functions
func (o ObjDataFrame) ShowValue() string {return o.Name}
func (o ObjDataFrame) Type() ValueType {return VAL_TABLE}
func (o ObjDataFrame) ToBytes() []byte {panic("implement me")}
func (o ObjDataFrame) ToValue() interface{} {return o}
func (o ObjDataFrame) Print() string {
	return "<table:"+o.Name+">"
}

func (o *ObjDataFrame) PrintHeader() {
	for c,_ := range o.ColNames {
		fmt.Printf("%s\t",o.ColNames[c])
	}
	fmt.Println()
}

func (o *ObjDataFrame) NextRow() bool {
	o.RowNumber++
	return o.Rows.Next()
}

func (o *ObjDataFrame) PrintRow() {
	res := make([]interface{},o.ColumnCount)

	for o.Rows.Next() {

		for i:= range o.Columns {
			res[i] = &res[i]
		}
		o.Rows.Scan(res...)
		for i := range o.Columns {
			fmt.Printf("%v\t",res[i])
		}
		fmt.Println()
	}
}

func (o *ObjDataFrame) PrintData(rows int64) {

	o.PrintHeader()

	res := make([]interface{},o.ColumnCount)
	for o.NextRow() {
		for i:= range o.Columns {
			res[i] = &res[i]
		}
		o.Rows.Scan(res...)
		for i := range o.Columns {
			fmt.Printf("%v\t",res[i])
		}
		fmt.Println()
		if o.RowNumber == rows && rows > 0 {
			break
		}
	}
}