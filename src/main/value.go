package main

import (
	"fmt"
	"hash/fnv"
)

type Obj interface {
	ShowValue() string
	Type() ValueType
	ToBytes() []byte
	ToValue() interface{}
}

type Iterator interface {
	Count() int
	Next() Obj
	First() Obj
	Current() Obj
	Position() int
}

// Internal types
type ObjColumnDef struct {
	TableName  string
	ColumnName string
	Alias      string
	Ordinal    byte
}

type oByte byte

func (o oByte) ShowValue() string    { return string(o) }
func (o oByte) Type() ValueType      { return VAL_BYTE }
func (o oByte) ToBytes() []byte      { return []byte{byte(o)} }
func (o oByte) ToValue() interface{} { return o }

// Interface types
func (c ObjColumnDef) ShowValue() string    { return fmt.Sprintf("%s", c.ColumnName) }
func (c ObjColumnDef) Type() ValueType      { return VAL_COLUMN_DEF }
func (c ObjColumnDef) ToBytes() []byte      { return nil }
func (c ObjColumnDef) ToValue() interface{} { return c.ShowValue() }

// User types
type ObjInteger int64
type ObjFloat float64 //struct{ Value float64 }
type ObjString string
type ObjBool struct{ Value bool }
type ObjByte struct{ Value byte }
type NULL struct{}
type ObjArray struct {
	ElementCount int
	ElementTypes ValueType
	Elements     []Obj
}

type ObjPointer struct {
}

type ObjMethod struct {
	Method *ObjClosure
	Class  *ObjClass
}

// Interface types
func (c ObjMethod) ShowValue() string    { return fmt.Sprintf("%s", "<Method>") }
func (c ObjMethod) Type() ValueType      { return VAL_METHOD }
func (c ObjMethod) ToBytes() []byte      { return nil }
func (c ObjMethod) ToValue() interface{} { return c.ShowValue() }

type ObjProperty struct {
	Property Obj
	Class    *ObjClass
}

// Interface types
func (c ObjProperty) ShowValue() string    { return fmt.Sprintf("%s", "<Property>") }
func (c ObjProperty) Type() ValueType      { return c.Property.Type() }
func (c ObjProperty) ToBytes() []byte      { return nil }
func (c ObjProperty) ToValue() interface{} { return c.ShowValue() }

var FunctionId int

type ObjFunction struct {
	Arity        int16
	Code         *Chunk
	UpvalueCount int
	Upvalues     []ObjUpvalue
	FuncType     FunctionType
	Id           int
}

type NativeFn func(vm *VM, argCounts int, stackPos int) Obj

type ObjNative struct {
	Function *NativeFn
}

var ClosureId int

type ObjClosure struct {
	Function     *ObjFunction
	Upvalues     []*ObjUpvalue
	UpvalueCount int16
	Id           int
}

var ObjLocation int // Global increment value

type ObjUpvalue struct {
	Reference *Obj
	Next      *ObjUpvalue
	Closed    Obj
	Location  int
}

type ObjList struct {
	KeyType      ValueType
	ElementCount int
	List         map[HashKey]Obj
}

var ClassId int

type ObjClass struct {
	Id    int
	Class *ObjClass

	Fields     map[string]Obj
	FieldCount int16

	Methods     []ObjClosure
	MethodCount int16
}

// To Encapulate new types
// todo Expand the concept of user defined datatypes
type ObjUdt struct {
	BaseType   ValueType
	BaseObject Obj
}

// Type functions

// Class functions
func NewClass(className string) *ObjClass {
	class := new(ObjClass)

	class.Id = ClassId
	ClassId++
	return class
}

func (c ObjClass) ShowValue() string    { return fmt.Sprintf("%s", "Class") }
func (c ObjClass) Type() ValueType      { return VAL_CLASS }
func (c ObjClass) ToBytes() []byte      { return nil }
func (c ObjClass) ToValue() interface{} { return c.ShowValue() }

// List functions

// This is to help when we translate content into hashkeys
type HashKey struct {
	Type      ValueType
	HashValue uint64
}

type HKey interface {
	HashValue() HashKey
}

func (l ObjList) ShowValue() string {
	s := ""
	for k, v := range l.List {
		s += fmt.Sprintf("(%v=%v)\n", k, v)
	}
	return s
}
func (l ObjList) Type() ValueType      { return VAL_LIST }
func (l ObjList) ToBytes() []byte      { return nil }
func (l ObjList) ToValue() interface{} { return l.List }

func (l *ObjList) Init(keyType ValueType, elementCount int) {
	l.ElementCount = elementCount
	l.KeyType = keyType
	l.List = make(map[HashKey]Obj)
}

func (l ObjList) GetValue(obj Obj) Obj {
	if obj.Type() == VAL_STRING {
		return l.List[obj.(ObjString).HashValue()]
	} else {
		return nil
	}
}

func (l ObjList) AddNew(key Obj, val Obj) {
	l.List[key.(ObjString).HashValue()] = val
}

func (l ObjList) SetValue(obj Obj, val Obj) {
	var hVal HashKey
	switch l.KeyType {
	case VAL_INTEGER:
		hVal = obj.(ObjInteger).HashValue()
	case VAL_STRING:
		hVal = obj.(ObjString).HashValue()
	}
	l.List[hVal] = val
}

// Upvalue functions
func (u ObjUpvalue) ShowValue() string {
	return fmt.Sprintf("%s", "Upvalue")
}
func (u ObjUpvalue) Type() ValueType      { return VAL_UPVALUE }
func (u ObjUpvalue) ToBytes() []byte      { return nil }
func (u ObjUpvalue) ToValue() interface{} { return u.ShowValue() }

func NewUpvalue(slot *Obj) *ObjUpvalue {
	upvalue := new(ObjUpvalue)
	upvalue.Reference = slot
	upvalue.Closed = new(NULL)
	upvalue.Next = nil
	upvalue.Location = ObjLocation

	ObjLocation++

	return upvalue
}

// Closure functions
func (c ObjClosure) ShowValue() string    { return fmt.Sprintf("%s", "<fn>") }
func (c ObjClosure) Type() ValueType      { return VAL_CLOSURE }
func (c ObjClosure) ToBytes() []byte      { return nil }
func (c ObjClosure) ToValue() interface{} { return c.ShowValue() }

func NewClosure(function *ObjFunction) *ObjClosure {
	// Make an array of upvalues of the same size as the number of
	// upvalues in the enclosed function
	upvalues := make([]*ObjUpvalue, function.UpvalueCount)
	// Increment the ID
	ClosureId++

	// Return the closure with a pointer to the function,
	// the array of upvalues and the count
	return &ObjClosure{
		Function:     function,
		Upvalues:     upvalues,
		UpvalueCount: int16(function.UpvalueCount),
		Id:           ClosureId - 1,
	}
}

// Function functions
func (f ObjFunction) ShowValue() string {
	return fmt.Sprintf("%s", "<fn>")
}
func (f ObjFunction) Type() ValueType      { return VAL_FUNCTION }
func (f ObjFunction) ToBytes() []byte      { return nil }
func (f ObjFunction) ToValue() interface{} { return f.ShowValue() }

func NewFunction() ObjFunction {

	code := NewChunk()

	function := ObjFunction{
		Arity:        0,
		Code:         &code,
		UpvalueCount: 0,
		Id:           FunctionId,
	}

	FunctionId++

	return function
}

// Native functions
func (n ObjNative) ShowValue() string {
	return fmt.Sprintf("%s", "<native fn>")
}
func (n ObjNative) Type() ValueType      { return VAL_NATIVE }
func (n ObjNative) ToBytes() []byte      { return nil }
func (n ObjNative) ToValue() interface{} { return n.ShowValue() }

func NewNative(function *NativeFn) *ObjNative {
	native := new(ObjNative)
	native.Function = function
	return native
}

// NULL functions
func (n NULL) ShowValue() string    { return fmt.Sprintf("%s", "null") }
func (n NULL) Type() ValueType      { return VAL_NIL }
func (n NULL) ToBytes() []byte      { return nil }
func (n NULL) ToValue() interface{} { return nil }

// Integer functions
func (i ObjInteger) ShowValue() string { return fmt.Sprintf("%d", i) }
func (i ObjInteger) Type() ValueType   { return VAL_INTEGER }
func (i ObjInteger) ToBytes() []byte   { return Int64ToBytes(int64(i)) }
func (i ObjInteger) ToValue() interface{} {
	return int64(i)
}

func (i ObjInteger) HashValue() HashKey {
	return HashKey{
		Type:      VAL_INTEGER,
		HashValue: uint64(i),
	}
}

// Float functions
func (f ObjFloat) ShowValue() string    { return fmt.Sprintf("%f", f) }
func (f ObjFloat) Type() ValueType      { return VAL_FLOAT }
func (f ObjFloat) ToBytes() []byte      { return Float64ToBytes(float64(f)) }
func (f ObjFloat) ToValue() interface{} { return f }

func (f *ObjFloat) HashValue() HashKey {
	return HashKey{
		Type:      VAL_FLOAT,
		HashValue: uint64(float64(*f)),
	}
}

// String functions
func (s ObjString) ShowValue() string    { return fmt.Sprintf("%s", s) }
func (s ObjString) Type() ValueType      { return VAL_STRING }
func (s ObjString) ToBytes() []byte      { return []byte(s) }
func (s ObjString) ToValue() interface{} { return s }

func (s ObjString) HashValue() HashKey {
	bw := fnv.New64a()
	_, _ = bw.Write([]byte(s))
	return HashKey{
		Type:      VAL_STRING,
		HashValue: bw.Sum64(),
	}
}

// Byte functions
func (b ObjByte) ShowValue() string    { return fmt.Sprintf("%d", b.Value) }
func (b ObjByte) Type() ValueType      { return VAL_BYTE }
func (b ObjByte) ToBytes() []byte      { return []byte{b.Value} }
func (b ObjByte) ToValue() interface{} { return b.Value }

// Bool functions
func (b ObjBool) ShowValue() string {
	if b.Value {
		return fmt.Sprintf("%s", "T")
	} else {
		return fmt.Sprintf("%s", "F")
	}
}
func (b ObjBool) Type() ValueType      { return VAL_BOOL }
func (b ObjBool) ToBytes() []byte      { return BoolToBytes(b.Value) }
func (b ObjBool) ToValue() interface{} { return b.Value }

func (b ObjBool) HashValue() HashKey {
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
func (a ObjArray) ShowValue() string {
	strVal := "|"
	for i := a.ElementCount - 1; i >= 0; i-- {
		strVal = strVal + a.Elements[i].ShowValue() + "|"
	}
	return strVal
}
func (a ObjArray) Type() ValueType { return VAL_ARRAY }
func (a ObjArray) ToBytes() []byte { return nil }

func (a ObjArray) Init(v ValueType, e int) {
	a.ElementCount = e
	a.ElementTypes = v
	a.Elements = make([]Obj, e)
}
func (a ObjArray) ToValue() interface{} { return a.Elements }

// Iterator interface
func (a ObjArray) Count() int {
	return a.ElementCount
}
func (a ObjArray) Next() Obj {
	a.ElementCount++
	return a.Elements[a.ElementCount-1]
}
func (a ObjArray) First() Obj {
	return a.Elements[0]
}
func (a ObjArray) Current() Obj {
	return a.Elements[a.ElementCount-1]
}
func (a ObjArray) Position() int {
	return a.ElementCount - 1
}

func (a ObjArray) GetElement(element int64) Obj {
	return a.Elements[element]
}

type ObjEnum struct {
}

// Utility functions
func MakeStringObj(str string) *ObjString {
	s := ObjString(str)
	return &s
}
