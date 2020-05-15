package main

const ROW_DB_CAPACITY = 1000

type Column struct {
	Name string
	ValType ValueType
	Ordinal int16
	Elements int64
	BlockSize int16
	StoragePtr []DataStorage
}

type ObjDataFrame struct {
	Name string
	Columns map[string]*Column
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
		Defined:     false, // Is the table finished being deined?
	}
}

func (o *ObjDataFrame) AddColumn(name string, valType ValueType) {
	// Check that the column doesn't already exist
	if _,ok := o.Columns[name];ok {
		panic("Column already exists")
	}
	o.Columns[o.ColumnCount] = Column{
		Name:       name,
		ValType:    valType,
		Ordinal:    o.ColumnCount+1,
		Elements:   0,
		StoragePtr: nil,
	}
	o.ColumnCount++
}

func (o *ObjDataFrame) AddRow(names []string, values []Obj) {
	if o.RowCapacity == o.RowCount {
		o.AllocateRows()
	}
	for i:=0;i<len(name);i++ {
		o.Columns[names[i]].StoragePtr[o.RowCount] = values[i].ToBytes()
	}
	o.RowCount++
}

func (o *ObjDataFrame) AllocateRows() {
	newData := make([]DataStorage,ROW_DB_CAPACITY)
	for k,_ := range o.Columns {
		o.Columns[k] = append(o.Columns[k],newData)
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

