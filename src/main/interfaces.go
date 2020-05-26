package main

type VarType byte

const (
	VAR_UNKNOWN VarType = iota
	VAR_SCALAR
	VAR_ARRAY
	VAR_FUNCTION
	VAR_HASH
	VAR_CLASS
	VAR_INSTANCE
	VAR_ENUM
	VAR_MATRIX
	VAR_TABLE
	VAR_RANGE
	VAR_OBJECT
)

var VarTypeLabel = map[VarType]string{
	VAR_UNKNOWN:  "Uknown",
	VAR_SCALAR:   "Scalar",
	VAR_ARRAY:    "Array",
	VAR_FUNCTION: "Function",
	VAR_HASH:     "Hash",
	VAR_CLASS:    "Class",
	VAR_INSTANCE: "Instance",
	VAR_ENUM:     "Enum",
	VAR_MATRIX:   "Matrix",
	VAR_TABLE: 	  "Table",
	VAR_RANGE:	  "Range",
	VAR_OBJECT:   "Object",

}

type ValueType byte

const (
	VAL_NIL ValueType = iota
	VAL_BOOL
	VAL_NUMBER
	VAL_BYTE
	VAL_INTEGER
	VAL_FLOAT
	VAL_STRING
	VAL_LIST
	VAL_ARRAY
	VAL_FUNCTION
	VAL_OBJ
	VAL_CLOSURE
	VAL_UPVALUE
	VAL_NATIVE
	VAL_CLASS
	VAL_INSTANCE
	VAL_SQLTABLE
	VAL_SQLCOLUMN
	VAL_COLUMN_DEF
	VAL_METHOD
	VAL_ENUM
	VAL_MATRIX
	VAL_TABLE
	VAL_RANGE
	VAL_OBJECT
)

var ValueTypeLabel = map[ValueType]string{
	VAL_NIL:        "nil",
	VAL_BOOL:       "bool",
	VAL_NUMBER:     "number",
	VAL_BYTE:       "byte",
	VAL_INTEGER:    "integer",
	VAL_FLOAT:      "float",
	VAL_STRING:     "string",
	VAL_LIST:       "list",
	VAL_ARRAY:      "array",
	VAL_FUNCTION:   "function",
	VAL_OBJ:        "obj",
	VAL_CLOSURE:    "closure",
	VAL_UPVALUE:    "upvalue",
	VAL_NATIVE:     "native",
	VAL_CLASS:      "class",
	VAL_INSTANCE:   "instance",
	VAL_SQLTABLE:   "SQLTable",
	VAL_SQLCOLUMN:  "SQLColumn",
	VAL_COLUMN_DEF: "ColumnDef",
	VAL_METHOD:     "Method",
	VAL_ENUM:       "Enum",
	VAL_MATRIX:     "Matrix",
	VAL_TABLE:      "Table" ,
	VAL_RANGE:      "Range" ,
	VAL_OBJECT:     "Object" ,
}

type FunctionType byte

const (
	TYPE_FUNCTION FunctionType = iota
	TYPE_SCRIPT
	TYPE_METHOD
	TYPE_NATIVE
)

type AccessorType byte

const (
	PUBLIC AccessorType = iota
	PRIVATE
	PROTECTED
)

type ClassComponentType byte

const (
	PROPERTY ClassComponentType = iota
	METHOD
)

type VariableScope byte

const (
	GLOBAL VariableScope = iota
	LOCAL
	UPVALUE
	REGISTER
	CLASS_PROPERTY
)

// This structure is going to get us to where we can
// create user defined data types
type VarInfo struct {
	Id int
	BaseType ValueType
	IsNative bool
}
var VarData = map[string]VarInfo {
	"int": {0,VAL_INTEGER,true},
	"float": {1,VAL_FLOAT, true},
	"bool": {2,VAL_BOOL, true },
	"string": {3, VAL_STRING, true},
	"byte":{4, VAL_BYTE, true },
	"nil": {5, VAL_NIL, true},
}


