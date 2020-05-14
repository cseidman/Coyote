package main

type Column struct {
	Name string
	ValType ValueType
	Ordinal int16
	Elements int64
	BlockSize int16
	StoragePtr *DataStorage
}

type ObjDataFrame struct {
	Name string
	Columns map[string]*Column
	ColumnCount int16
	Defined bool
}

// Interface functions
func (o ObjDataFrame) ShowValue() string {return o.Name}
func (o ObjDataFrame) Type() ValueType {return VAL_TABLE}
func (o ObjDataFrame) ToBytes() []byte {panic("implement me")}
func (o ObjDataFrame) ToValue() interface{} {return o}
func (o ObjDataFrame) Print() string {return "<table>"}

// Data Storage
type DataStorage []byte

func CreateTable(name string) ObjDataFrame {
	// Start off by creating a barebones table with the expectation that we're
	// going to add columns later
	return ObjDataFrame{
		Name:        name,
		Columns:     make(map[string]*Column),
		ColumnCount: 0,
		Defined:     false,
	}
}

func (o *ObjDataFrame) AddColumn(name string, valType ValueType) {
	// Check that the column doesn't already exist
	if _,ok := o.Columns[name];ok {
		panic("Column already exists")
	}
	o.ColumnCount++
	o.Columns[o.ColumnCount] = Column{
		Name:       name,
		ValType:    valType,
		Ordinal:    o.ColumnCount+1,
		Elements:   0,
		StoragePtr: nil,
	}
}

func (o *ObjDataFrame) AddRow(names []string, values []Obj) {
	for i:=0;i<len(name);i++ {
		o.Columns[names[i]].StoragePtr = values[i].ToBytes()
	}
}

// Specific functions

