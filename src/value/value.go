package value

import (
	. "../chunk"
	. "../common"
	"fmt"
	"hash/fnv"
)

// Internal types
type ObjColumnDef struct {
	TableName string
	ColumnName string
	Alias string
	Ordinal byte
}

// Interface types
func (c *ObjColumnDef) ShowValue() string { return fmt.Sprintf("%s", c.ColumnName) }
func (c *ObjColumnDef) Type() ValueType   { return VAL_COLUMN_DEF }
func (c *ObjColumnDef) Add(col *ObjColumnDef) *ObjColumnDef {return nil}

// User types
type ObjInteger struct{ Value int64 }
type ObjFloat struct{ Value float64 }
type ObjString struct{ Value string }
type ObjBool struct{ Value bool }
type ObjByte struct{ Value byte }
type NULL struct{}
type ObjArray struct {
	ElementCount int
	ElementTypes ValueType
	Elements     []Obj
}

var FunctionId int

type ObjFunction struct {
	Arity        int
	FChunk       Chunk
	Name         string
	UpvalueCount int
	FuncType     FunctionType
	Id           int
}

type NativeFn func(argCounts int, stackPos int) Obj

type ObjNative struct {
	Function NativeFn
}

var ClosureId int

type ObjClosure struct {
	Function     *ObjFunction
	Upvalues     []*ObjUpvalue
	UpvalueCount int
	Id           int
}

var ObjLocation int // Global increment value

type ObjUpvalue struct {
	Reference *Obj
	Next      *ObjUpvalue
	Closed    Obj
	Location  int // Increments at every creation
}

type ObjList struct {
	KeyType      ValueType
	ElementCount int
	List         map[HashKey]Obj
}

type ObjClass struct {
	Name string
	Methods map[string]Obj
}

type ObjInstance struct {
	Class  *ObjClass
	Fields map[string]Obj
}

// Class functions
func NewClass(className string) *ObjClass {
	class := new(ObjClass)
	class.Name = className
	class.Methods = make(map[string]Obj)
	return class
}

func (c *ObjClass) ShowValue() string { return fmt.Sprintf("%s", c.Name) }
func (c *ObjClass) Type() ValueType   { return VAL_CLASS }
func (c *ObjClass) Add(*ObjClass) *ObjClass {return nil}

// Class instance functions
func NewInstance(class *ObjClass) *ObjInstance {
	instance := new(ObjInstance)
	instance.Class = class
	instance.Fields = make(map[string]Obj)
	return instance
}

func (i *ObjInstance) ShowValue() string { return fmt.Sprintf("Instance of %s", i.Class.Name) }
func (i *ObjInstance) Type() ValueType   { return VAL_INSTANCE }
func (i *ObjInstance) Add(*ObjInstance) *ObjInstance {return nil}

func (i *ObjInstance) GetFieldValue(fieldName string) *Obj {
	if val, ok := i.Fields[fieldName]; ok {
		return &val
	} else {
		return nil
	}
}

func (i *ObjInstance) SetFieldValue(fieldName string, val *Obj) {
	i.Fields[fieldName] = *val
}

// List functions

// This is to help when we translate content into hashkeys
type HashKey struct {
	Type      ValueType
	HashValue uint64
}

type HKey interface {
	HashValue() HashKey
}

func (l *ObjList) ShowValue() string {
	return fmt.Sprintf("%s", "List")
}
func (l *ObjList) Type() ValueType { return VAL_LIST }
func (l *ObjList) Add(*ObjList) *ObjList {return nil}

func (l *ObjList) Init(keyType ValueType, elementCount int) {
	l.ElementCount = elementCount
	l.KeyType = keyType
	l.List = make(map[HashKey]Obj)
}

func (l *ObjList) GetValue(obj Obj) Obj {
	return l.List[obj.(HKey).HashValue()]
}

func (l *ObjList) AddNew(key Obj, val Obj) {
	l.List[key.(HKey).HashValue()] = val
}

// Upvalue functions
func (u *ObjUpvalue) ShowValue() string {
	return fmt.Sprintf("%s", "Upvalue")
}
func (u *ObjUpvalue) Type() ValueType { return VAL_UPVALUE }
func (u *ObjUpvalue) Add() *ObjUpvalue {return nil}

func NewUpvalue(slot *Obj) *ObjUpvalue {
	upvalue := new(ObjUpvalue)
	upvalue.Reference = slot
	//upvalue.Closed = new(NULL)
	upvalue.Next = nil
	upvalue.Location = ObjLocation

	ObjLocation++

	return upvalue
}

// Closure functions
func (c *ObjClosure) ShowValue() string { return fmt.Sprintf("%s", c.Function.Name) }
func (c *ObjClosure) Type() ValueType   { return VAL_CLOSURE }

func NewClosure(function *ObjFunction) *ObjClosure {
	// Make an array of upvalues of the same size as the number of
	// upvalues in the enclosed function
	upvalues := make([]*ObjUpvalue, function.UpvalueCount)

	// Set them all to nil (so they don't get garbage collected yet))
	for i := 0; i < function.UpvalueCount; i++ {
		upvalues[i] = nil
	}
	// Increment the ID
	ClosureId++

	// Return the closure with a pointer to the function,
	// the array of upvalues and the count
	return &ObjClosure{
		Function:     function,
		Upvalues:     upvalues,
		UpvalueCount: function.UpvalueCount,
		Id:           ClosureId - 1,
	}
}

// Function functions
func (f *ObjFunction) ShowValue() string { return fmt.Sprintf("%s", f.Name) }
func (f *ObjFunction) Type() ValueType   { return VAL_FUNCTION }
/*
func NewFunction() *ObjFunction {

	chunk := New
	InitChunk(chunk)

	function := &ObjFunction{
		Arity:        0,
		Name:         "",
		FChunk:       *chunk,
		UpvalueCount: 0,
		Id:           FunctionId,
	}

	FunctionId++

	return function
}
*/
// Native functions
func (n *ObjNative) ShowValue() string { return fmt.Sprintf("%s", "<native fn>") }
func (n *ObjNative) Type() ValueType   { return VAL_NATIVE }

func NewNative(function NativeFn) *ObjNative {
	native := new(ObjNative)
	native.Function = function
	return native
}

// NULL functions
func (n *NULL) ShowValue() string { return fmt.Sprintf("%s", "null") }
func (n *NULL) Type() ValueType   { return VAL_NIL }

// Integer functions
func (i *ObjInteger) ShowValue() string { return fmt.Sprintf("%d", i.Value) }
func (i *ObjInteger) Type() ValueType   { return VAL_INTEGER }
func (i *ObjInteger) HashValue() HashKey {
	return HashKey{
		Type:      VAL_INTEGER,
		HashValue: uint64(i.Value),
	}
}
func (i *ObjInteger) Add(oInt *ObjInteger) *ObjInteger {
	return &ObjInteger{i.Value + oInt.Value}
}

// Float functions
func (f *ObjFloat) ShowValue() string { return fmt.Sprintf("%f", f.Value) }
func (f *ObjFloat) Type() ValueType   { return VAL_FLOAT }
func (f *ObjFloat) HashValue() HashKey {
	return HashKey{
		Type:      VAL_FLOAT,
		HashValue: uint64(f.Value),
	}
}
func (f *ObjFloat) Add(oFloat *ObjFloat) *ObjFloat {
	return &ObjFloat{f.Value + oFloat.Value}
}

// String functions
func (s *ObjString) ShowValue() string { return fmt.Sprintf("%s", s.Value) }
func (s *ObjString) Type() ValueType   { return VAL_STRING }
func (s *ObjString) Add(*ObjString) *ObjString   { return nil }

func (s *ObjString) HashValue() HashKey {
	bw := fnv.New64a()
	_, _ = bw.Write([]byte(s.Value))
	return HashKey{
		Type:      VAL_STRING,
		HashValue: bw.Sum64(),
	}
}

// Byte functions
func (b *ObjByte) ShowValue() string { return fmt.Sprintf("%d", b.Value) }
func (b *ObjByte) Type() ValueType   { return VAL_BYTE }
func (b *ObjByte) Add(bval *ObjByte) *ObjByte {return &ObjByte{b.Value + bval.Value} }

// Bool functions
func (b *ObjBool) ShowValue() string {
	if b.Value {
		return fmt.Sprintf("%s", "T")
	} else {
		return fmt.Sprintf("%s", "F")
	}
}
func (b *ObjBool) Type() ValueType { return VAL_BOOL }
func (b *ObjBool) Add(btype *ObjBool) *ObjBool {return &ObjBool{b.Value && btype.Value}}

func (b *ObjBool) HashValue() HashKey {
	var val uint64

	if b.Value {
		val = 1
	} else {
		val = 0
	}
	return HashKey{
		Type:      VAL_BOOL,
		HashValue: val,
	}
}

// Array functions
func (a *ObjArray) ShowValue() string {
	strVal := "|"
	for i := a.ElementCount - 1; i >= 0; i-- {
		strVal = strVal + a.Elements[i].ShowValue() + "|"
	}
	return strVal
}
func (a *ObjArray) Type() ValueType { return VAL_ARRAY }
func (a *ObjArray) Add(*ObjArray) *ObjArray {return nil}

func (a *ObjArray) Init(v ValueType, e int) {
	a.ElementCount = e
	a.ElementTypes = v
	a.Elements = make([]Obj, e)
}
func (a *ObjArray) GetElement(element int64) Obj {
	return a.Elements[element]
}

// Utility functions
func MakeStringObj(str string)*ObjString {
	obj := new(ObjString)
	obj.Value = str
	return obj
}