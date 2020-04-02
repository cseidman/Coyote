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
)

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
)

type FunctionType byte

const (
	TYPE_FUNCTION FunctionType = iota
	TYPE_SCRIPT
	TYPE_METHOD
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
