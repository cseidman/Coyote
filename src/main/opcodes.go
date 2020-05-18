package main

const (
	OP_HALT byte = iota
	OP_CONSTANT
	OP_PUSH
	OP_NIL
	OP_TRUE // 5
	OP_FALSE
	OP_EQUAL
	OP_GREATER
	OP_GREATER_EQUAL
	OP_LESS //10
	OP_LESS_EQUAL
	OP_INEGATE
	OP_FNEGATE
	OP_IADD
	OP_FADD
	OP_ISUBTRACT //15
	OP_FSUBTRACT
	OP_IMULTIPLY
	OP_FMULTIPLY
	OP_IDIVIDE
	OP_FDIVIDE // 20
	OP_INCREMENT
	OP_DECREMENT
	OP_PREINCREMENT
	OP_PREDECREMENT
	OP_SADD //25
	OP_NOT
	OP_RETURN
	OP_DECLARE
	OP_PRINT
	OP_POP //30
	OP_DEFINE_GLOBAL
	OP_DEFINE_LOCAL
	OP_GET_GLOBAL
	OP_SET_GLOBAL
	OP_GET_LOCAL //35
	OP_SET_LOCAL
	OP_JUMP_IF_FALSE
	OP_JUMP
	OP_CALL //40
	OP_ARRAY
	OP_INDEX
	OP_LIST
	OP_GET_ALOCAL
	OP_SET_ALOCAL //45
	OP_GET_AGLOBAL
	OP_SET_AGLOBAL
	OP_CLOSURE
	OP_GET_UPVALUE
	OP_SET_UPVALUE //50
	OP_GET_AUPVALUE
	OP_SET_AUPVALUE
	OP_CLOSE_UPVALUE
	OP_CLASS
	OP_METHOD //55
	OP_SET_PROPERTY
	OP_GET_PROPERTY
	OP_CREATE_TABLE
	OP_CREATE_COLUMN
	OP_INSERT //60
	OP_DISPLAY_TABLE
	OP_SQL_SELECT
	OP_ALL_COLUMNS
	OP_JOIN
	OP_LEFT_JOIN //65
	OP_RIGHT_JOIN
	OP_CROSS_JOIN
	OP_COLUMN
	OP_TABLE
	OP_EXPRESSION_COLUMN //70
	OP_JOIN_EXPRESSION
	OP_NOTHING
	OP_SET_REGISTER
	OP_GET_REGISTER
	OP_FOR_LOOP
	OP_ICONST
	OP_FCONST
	OP_SCONST
	OP_BCONST
	OP_BREAK
	OP_NOT_EQUAL
	OP_CONTINUE
	OP_FN_CONST
	OP_PUSH_0
	OP_PUSH_1
	OP_PUSH_2
	OP_PUSH_3
	OP_PUSH_4
	OP_PUSH_5
	OP_IEXP
	OP_FEXP
	OP_SCAN
	OP_AINDEX
	OP_ASIZE
	OP_SUBCLASS
	OP_GET_METHOD
	OP_SET_METHOD
	OP_BIND_METHOD
	OP_BIND_PROPERTY
	OP_CALL_METHOD
	OP_GET_HLOCAL
	OP_SET_HLOCAL
	OP_GET_HGLOBAL
	OP_SET_HGLOBAL
	OP_GET_LOCAL_0
	OP_GET_LOCAL_1
	OP_GET_LOCAL_2
	OP_GET_LOCAL_3
	OP_GET_LOCAL_4
	OP_GET_LOCAL_5
	OP_CALL_0
	OP_CALL_1
	OP_CALL_2
	OP_CALL_3
	OP_GET_GLOBAL_0
	OP_GET_GLOBAL_1
	OP_GET_GLOBAL_2
	OP_GET_GLOBAL_3
	OP_GET_GLOBAL_4
	OP_GET_GLOBAL_5
	OP_CALL_NATIVE
	OP_ENUM
	OP_ENUM_TAG
	OP_DROP_TABLE
	OP_IRANGE
	OP_MAKE_ARRAY
)

var OpLabel = map[byte]string{
	OP_HALT:     "OP_HALT",
	OP_CONSTANT: "OP_CONSTANT",
	OP_PUSH:     "OP_PUSH",
	OP_NIL:      "OP_NIL",
	OP_TRUE:     "OP_TRUE",

	OP_FALSE:         "OP_FALSE",
	OP_EQUAL:         "OP_IEQUAL",
	OP_GREATER:       "OP_GREATER",
	OP_GREATER_EQUAL: "OP_GREATER_EQUAL",
	OP_LESS:          "OP_LESS",

	OP_LESS_EQUAL: "OP_LESS_EQUAL",
	OP_INEGATE:    "OP_INEGATE",
	OP_NOT:        "OP_NOT",
	OP_IADD:       "OP_IADD",
	OP_FADD:       "OP_FADD",

	OP_ISUBTRACT: "OP_ISUBTRACT",
	OP_FSUBTRACT: "OP_FSUBTRACT",
	OP_IMULTIPLY: "OP_IMULTIPLY",
	OP_FMULTIPLY: "OP_FMULTIPLY",
	OP_IDIVIDE:   "OP_IDIVIDE",

	OP_FDIVIDE:      "OP_FDIVIDE",
	OP_INCREMENT:    "OP_INCREMENT",
	OP_DECREMENT:    "OP_DECREMENT",
	OP_PREINCREMENT: "OP_PREINCREMENT",
	OP_PREDECREMENT: "OP_PREDECREMENT",

	OP_SADD:    "OP_SADD",
	OP_RETURN:  "OP_RETURN",
	OP_DECLARE: "OP_DECLARE", //20
	OP_PRINT:   "OP_PRINT",
	OP_POP:     "OP_POP",

	OP_DEFINE_GLOBAL: "OP_DEFINE_GLOBAL",
	OP_DEFINE_LOCAL:  "OP_DEFINE_LOCAL",
	OP_GET_GLOBAL:    "OP_GET_GLOBAL",
	OP_SET_GLOBAL:    "OP_SET_GLOBAL",
	OP_GET_LOCAL:     "OP_GET_LOCAL",

	OP_SET_LOCAL:     "OP_SET_LOCAL",
	OP_JUMP_IF_FALSE: "OP_JUMP_IF_FALSE",
	OP_JUMP:          "OP_JUMP",
	OP_CALL:          "OP_CALL",

	OP_ARRAY:      "OP_ARRAY",
	OP_INDEX:      "OP_INDEX",
	OP_LIST:       "OP_LIST",
	OP_GET_ALOCAL: "OP_GET_ALOCAL",
	OP_SET_ALOCAL: "OP_SET_ALOCAL",

	OP_GET_AGLOBAL: "OP_GET_AGLOBAL",
	OP_SET_AGLOBAL: "OP_SET_AGLOBAL",
	OP_CLOSURE:     "OP_CLOSURE",
	OP_GET_UPVALUE: "OP_GET_UPVALUE",
	OP_SET_UPVALUE: "OP_SET_UPVALUE",

	OP_GET_AUPVALUE:  "OP_GET_AUPVALUE",
	OP_SET_AUPVALUE:  "OP_SET_AUPVALUE",
	OP_CLOSE_UPVALUE: "OP_CLOSE_UPVALUE",
	OP_CLASS:         "OP_CLASS",
	OP_METHOD:        "OP_METHOD",

	OP_SET_PROPERTY:  "OP_SET_PROPERTY",
	OP_GET_PROPERTY:  "OP_GET_PROPERTY",
	OP_CREATE_TABLE:  "OP_CREATE_TABLE",
	OP_CREATE_COLUMN: "OP_CREATE_COLUMN",
	OP_INSERT:        "OP_INSERT",

	OP_DISPLAY_TABLE: "OP_DISPLAY_TABLE",
	OP_SQL_SELECT:    "OP_SQL_SELECT",
	OP_ALL_COLUMNS:   "OP_ALL_COLUMNS",
	OP_JOIN:          "OP_JOIN",
	OP_LEFT_JOIN:     "OP_LEFT_JOIN",

	OP_RIGHT_JOIN:        "OP_RIGHT_JOIN",
	OP_CROSS_JOIN:        "OP_CROSS_JOIN",
	OP_COLUMN:            "OP_COLUMN",
	OP_TABLE:             "OP_TABLE",
	OP_EXPRESSION_COLUMN: "OP_EXPRESSION_COLUMN",

	OP_JOIN_EXPRESSION: "OP_JOIN_EXPRESSION",
	OP_NOTHING:         "OP_NOTHING",
	OP_SET_REGISTER:    "OP_SET_REGISTER",
	OP_GET_REGISTER:    "OP_GET_REGISTER",
	OP_FOR_LOOP:        "OP_FOR_LOOP",

	OP_ICONST: "OP_ICONST",
	OP_FCONST: "OP_FCONST",
	OP_SCONST: "OP_SCONST",
	OP_BCONST: "OP_BCONST",
	OP_BREAK:  "OP_BREAK",

	OP_NOT_EQUAL: "OP_NOT_EQUAL",
	OP_CONTINUE:  "OP_CONTINUE",
	OP_FN_CONST:  "OP_FN_CONST",
	OP_PUSH_0:    "OP_PUSH_0",
	OP_IEXP:      "OP_IEXP",

	OP_PUSH_1: "OP_PUSH_1",
	OP_PUSH_2: "OP_PUSH_2",
	OP_PUSH_3: "OP_PUSH_3",
	OP_PUSH_4: "OP_PUSH_4",
	OP_PUSH_5: "OP_PUSH_5",

	OP_FEXP:    "OP_FEXP",
	OP_FNEGATE: "OP_FNEGATE",
	OP_SCAN:    "OP_SCAN",
	OP_AINDEX:  "OP_AINDEX",
	OP_ASIZE:   "OP_ASIZE",

	OP_SUBCLASS:      "OP_SUBCLASS",
	OP_GET_METHOD:    "OP_GET_METHOD",
	OP_SET_METHOD:    "OP_SET_METHOD",
	OP_BIND_METHOD:   "OP_BIND_METHOD",
	OP_BIND_PROPERTY: "OP_BIND_PROPERTY",

	OP_CALL_METHOD: "OP_CALL_METHOD",
	OP_GET_HLOCAL:  "OP_GET_HLOCAL",
	OP_SET_HLOCAL:  "OP_SET_HLOCAL",
	OP_GET_HGLOBAL: "OP_GET_HGLOBAL",
	OP_SET_HGLOBAL: "OP_SET_HGLOBAL",

	OP_GET_LOCAL_0: "OP_GET_LOCAL_0",
	OP_GET_LOCAL_1: "OP_GET_LOCAL_1",
	OP_GET_LOCAL_2: "OP_GET_LOCAL_2",
	OP_GET_LOCAL_3: "OP_GET_LOCAL_3",
	OP_GET_LOCAL_4: "OP_GET_LOCAL_4",

	OP_GET_LOCAL_5: "OP_GET_LOCAL_5",
	OP_CALL_0:      "OP_CALL_0",
	OP_CALL_1:      "OP_CALL_1",
	OP_CALL_2:      "OP_CALL_2",
	OP_CALL_3:      "OP_CALL_3",

	OP_GET_GLOBAL_0: "OP_GET_GLOBAL_0",
	OP_GET_GLOBAL_1: "OP_GET_GLOBAL_1",
	OP_GET_GLOBAL_2: "OP_GET_GLOBAL_2",
	OP_GET_GLOBAL_3: "OP_GET_GLOBAL_3",
	OP_GET_GLOBAL_4: "OP_GET_GLOBAL_4",

	OP_GET_GLOBAL_5: "OP_GET_GLOBAL_5",
	OP_CALL_NATIVE:  "OP_CALL_NATIVE",
	OP_ENUM:         "OP_ENUM",
	OP_ENUM_TAG:     "OP_ENUM_TAG",
	OP_DROP_TABLE:   "OP_DROP_TABLE",

	OP_IRANGE:		 "OP_IRANGE",
	OP_MAKE_ARRAY:   "OP_MAKE_ARRAY",

}
