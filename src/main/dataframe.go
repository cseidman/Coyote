package main

import (
	"fmt"
)

const ROW_DB_GROWTH = 1000 // Increment new rows bu this number

// Data Object Description
type DBObjectEntry struct {

}
// Data Columsn description
type DBCol struct {

}

// Data dictionary
type DataDict struct {
	DBId int64 // Database id
	DBName string // Database name


	DBObjCount int

}

type Column struct {
	Name string
	ValType ValueType
	Ordinal int16
	Elements int64
	BlockSize int16
	StoragePtr []DataStorage
}

func (c *Column) GetValue(row int64) Obj {
	switch c.ValType {
		case VAL_STRING:
			return ObjString(string(c.StoragePtr[row]))
	}
	return nil
}

type ObjDataFrame struct {
	Name string
	Columns map[string]*Column
	ColNames []string
	ColumnCount int16
	RowCount int64
	Defined bool
	RowCapacity int64
}

// Interface functions
func (o ObjDataFrame) ShowValue() string {return o.Name}
func (o ObjDataFrame) Type() ValueType {return VAL_TABLE}
func (o ObjDataFrame) ToBytes() []byte {panic("implement me")}
func (o ObjDataFrame) ToValue() interface{} {return o}
func (o ObjDataFrame) Print() string {
	str := ""
	for col := range o.ColNames {
		str += fmt.Sprintf("%s\t",col)
	}
	str+="\n"
	for r:=int64(0);r<o.RowCount;r++ {
		for col := range o.ColNames {
			str += fmt.Sprintf("%s\t", o.Columns[o.ColNames[col]].GetValue(r).Print())
		}
	}

	return str
}

// Data Storage
type DataStorage []byte

func CreateTable(name string) ObjDataFrame {
	// Start off by creating a barebones table with the expectation that we're
	// going to add columns later
	return ObjDataFrame{
		Name:        name,
		Columns:     make(map[string]*Column),
		ColumnCount: 0,
		Defined:     false, // Is the table finished being deined?
		ColNames:    make([]string,0),
	}
}

func (o *ObjDataFrame) AddColumn(name string, valType ValueType) {
	// Check that the column doesn't already exist
	if _,ok := o.Columns[name];ok {
		panic("Column already exists")
	}
	o.Columns[name] = &Column{
		Name:       name,
		ValType:    valType,
		Ordinal:    o.ColumnCount+1,
		Elements:   0,
		StoragePtr: nil,
	}
	o.ColNames = append(o.ColNames,name)
	o.ColumnCount++
}

func (o *ObjDataFrame) AddRow(names []string, values []Obj) {
	if o.RowCapacity == o.RowCount {
		o.AllocateRows()
	}
	for i:=0;i<len(names);i++ {
		o.Columns[names[i]].StoragePtr[o.RowCount] = values[i].ToBytes()
	}
	o.RowCount++
}

func (o *ObjDataFrame) AllocateRows() {
	newData := make([]DataStorage,ROW_DB_GROWTH)
	for k,_ := range o.Columns {
		o.Columns[k].StoragePtr = append(o.Columns[k].StoragePtr,newData...)
	}
}

func (o *ObjDataFrame) QueryBuilder() {

	// 1) First POP the data sources off the stack
	// 2) Assign them each an alias if they don't have one
	// 3) Make sure the aliases aren't repeated

	// 1) Evaluate JOIN conditions

	// 1) Get a list of all the columns in all the data sources
	// 2) Match the columns in the SELECT statement with the ones in the above list

	// Get the filer conditions

}

